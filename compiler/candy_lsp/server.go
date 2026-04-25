package candy_lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Server implements a minimal diagnostics-only LSP loop.
type Server struct {
	docs           map[string]documentState
	lastDiag       map[string]string
	workspaceIndex map[string][]symbolDef
	shutdown       bool
	exited         bool
	writeError     error
}

type documentState struct {
	URI      string
	Path     string
	Text     string
	Version  int
	Hash     string
	Imports  []string
	Symbols  []symbolDef
	ImportTo []string
}

// Run serves LSP messages over stdio-like streams.
func Run(in io.Reader, out io.Writer) error {
	s := &Server{
		docs:           map[string]documentState{},
		lastDiag:       map[string]string{},
		workspaceIndex: map[string][]symbolDef{},
	}
	br := bufio.NewReader(in)
	for {
		reqBytes, err := readLSPPacket(br)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		var req rpcRequest
		if err := json.Unmarshal(reqBytes, &req); err != nil {
			return err
		}
		if err := s.handle(req, out); err != nil {
			return err
		}
		if s.exited {
			return nil
		}
	}
}

func (s *Server) handle(req rpcRequest, out io.Writer) error {
	if req.JSONRPC != "" && req.JSONRPC != "2.0" {
		return s.writeRPCError(out, req.ID, -32600, "invalid jsonrpc version")
	}
	switch req.Method {
	case methodInitialize:
		return writeResponse(out, rpcResponse{
			JSONRPC: "2.0",
			ID:      decodeID(req.ID),
			Result: initializeResult{
				Capabilities: serverCapabilities{
					TextDocumentSync:   1, // full sync
					HoverProvider:      true,
					DefinitionProvider: true,
				},
			},
		})
	case methodInitialized:
		return nil
	case methodDidOpen:
		var p didOpenParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return s.writeRPCError(out, req.ID, -32602, "invalid didOpen params")
		}
		ds := documentState{
			URI:     p.TextDocument.URI,
			Path:    uriToPath(p.TextDocument.URI),
			Text:    p.TextDocument.Text,
			Version: 1,
		}
		s.docs[p.TextDocument.URI] = s.reindexDocument(documentState{}, ds)
		return s.publishDiagnostics(out, p.TextDocument.URI, p.TextDocument.Text)
	case methodDidChange:
		var p didChangeParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return s.writeRPCError(out, req.ID, -32602, "invalid didChange params")
		}
		old := s.docs[p.TextDocument.URI]
		ds := old
		if len(p.ContentChanges) > 0 {
			ds.Text = p.ContentChanges[len(p.ContentChanges)-1].Text
		}
		if p.TextDocument.Version != 0 {
			ds.Version = p.TextDocument.Version
		} else {
			ds.Version++
		}
		ds.URI = p.TextDocument.URI
		ds.Path = uriToPath(p.TextDocument.URI)
		s.docs[p.TextDocument.URI] = s.reindexDocument(old, ds)
		return s.publishDiagnostics(out, p.TextDocument.URI, ds.Text)
	case methodDidClose:
		var p didCloseParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return s.writeRPCError(out, req.ID, -32602, "invalid didClose params")
		}
		if ds, ok := s.docs[p.TextDocument.URI]; ok {
			s.removeSymbols(ds.Symbols)
		}
		delete(s.docs, p.TextDocument.URI)
		delete(s.lastDiag, p.TextDocument.URI)
		return writeNotification(out, rpcNotification{
			JSONRPC: "2.0",
			Method:  methodPublishDiagnostics,
			Params: publishDiagnosticsParams{
				URI:         p.TextDocument.URI,
				Diagnostics: []lspDiagnostic{},
			},
		})
	case methodHover:
		var p textDocumentPositionParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return s.writeRPCError(out, req.ID, -32602, "invalid hover params")
		}
		ds, ok := s.docs[p.TextDocument.URI]
		if !ok {
			return writeResponse(out, rpcResponse{JSONRPC: "2.0", ID: decodeID(req.ID), Result: nil})
		}
		h := buildHover(ds.Text, p.Position)
		return writeResponse(out, rpcResponse{
			JSONRPC: "2.0",
			ID:      decodeID(req.ID),
			Result:  h,
		})
	case methodDefinition:
		var p textDocumentPositionParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return s.writeRPCError(out, req.ID, -32602, "invalid definition params")
		}
		ds, ok := s.docs[p.TextDocument.URI]
		if !ok {
			return writeResponse(out, rpcResponse{JSONRPC: "2.0", ID: decodeID(req.ID), Result: []lspLocation{}})
		}
		def := s.findDefinition(ds, p.Position)
		return writeResponse(out, rpcResponse{
			JSONRPC: "2.0",
			ID:      decodeID(req.ID),
			Result:  def,
		})
	case methodShutdown:
		s.shutdown = true
		return writeResponse(out, rpcResponse{
			JSONRPC: "2.0",
			ID:      decodeID(req.ID),
			Result:  nil,
		})
	case methodExit:
		s.exited = true
		return nil
	default:
		if len(req.ID) > 0 {
			return s.writeRPCError(out, req.ID, -32601, "method not found: "+req.Method)
		}
		return nil // notification: ignore
	}
}

func (s *Server) publishDiagnostics(out io.Writer, uri, source string) error {
	raw := AnalyzeSource(source)
	diags := make([]lspDiagnostic, 0, len(raw))
	for _, d := range raw {
		diags = append(diags, toLSPDiagnostic(d))
	}
	b, _ := json.Marshal(diags)
	hash := string(b)
	if prev, ok := s.lastDiag[uri]; ok && prev == hash {
		return nil
	}
	s.lastDiag[uri] = hash
	return writeNotification(out, rpcNotification{
		JSONRPC: "2.0",
		Method:  methodPublishDiagnostics,
		Params: publishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diags,
		},
	})
}

func (s *Server) writeRPCError(out io.Writer, id json.RawMessage, code int, msg string) error {
	return writeResponse(out, rpcResponse{
		JSONRPC: "2.0",
		ID:      decodeID(id),
		Error: rpcError{
			Code:    code,
			Message: msg,
		},
	})
}

var identRE = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]*`)

func buildHover(src string, pos lspPosition) interface{} {
	w, r := wordAtPosition(src, pos)
	if w == "" {
		return nil
	}
	info := "identifier"
	lw := strings.ToLower(w)
	switch lw {
	case "len", "ok", "err":
		info = "builtin"
	}
	if strings.Contains(r, "(") {
		info = "function"
	}
	return hoverResult{
		Contents: fmt.Sprintf("`%s` (%s)", w, info),
		Range:    rangeForWord(pos.Line, w, r),
	}
}

func (s *Server) findDefinition(ds documentState, pos lspPosition) []lspLocation {
	w, _ := wordAtPosition(ds.Text, pos)
	if w == "" {
		return []lspLocation{}
	}
	// 1) Local definitions first.
	if m := bestLocalMatch(ds.Symbols, w, pos); m != nil {
		return []lspLocation{{URI: m.URI, Range: m.Range}}
	}
	// 2) Definitions in imported documents.
	for _, impURI := range ds.ImportTo {
		imp, ok := s.docs[impURI]
		if !ok {
			continue
		}
		if m := bestGlobalMatch(imp.Symbols, w); m != nil {
			return []lspLocation{{URI: m.URI, Range: m.Range}}
		}
	}
	// 3) Workspace fallback.
	if m := bestGlobalMatch(s.workspaceIndex[w], w); m != nil {
		return []lspLocation{{URI: m.URI, Range: m.Range}}
	}
	return []lspLocation{}
}

func (s *Server) reindexDocument(old, ds documentState) documentState {
	hash := ds.Text
	if old.Hash == hash && len(old.Symbols) > 0 {
		old.URI = ds.URI
		old.Path = ds.Path
		old.Version = ds.Version
		old.Text = ds.Text
		return old
	}
	s.removeSymbols(old.Symbols)
	symbols, imports := buildDocumentIndex(ds.URI, ds.Text)
	importTo := make([]string, 0, len(imports))
	for _, imp := range imports {
		if uri := resolveImportURI(ds.URI, ds.Path, imp); uri != "" {
			importTo = append(importTo, uri)
		}
	}
	sort.Strings(importTo)
	ds.Hash = hash
	ds.Symbols = symbols
	ds.Imports = imports
	ds.ImportTo = importTo
	s.addSymbols(symbols)
	return ds
}

func (s *Server) addSymbols(symbols []symbolDef) {
	for _, sym := range symbols {
		name := strings.ToLower(sym.Name)
		s.workspaceIndex[name] = append(s.workspaceIndex[name], sym)
	}
	for name := range s.workspaceIndex {
		sort.Slice(s.workspaceIndex[name], func(i, j int) bool {
			a, b := s.workspaceIndex[name][i], s.workspaceIndex[name][j]
			if a.URI != b.URI {
				return a.URI < b.URI
			}
			if a.Range.Start.Line != b.Range.Start.Line {
				return a.Range.Start.Line < b.Range.Start.Line
			}
			return a.Range.Start.Character < b.Range.Start.Character
		})
	}
}

func (s *Server) removeSymbols(symbols []symbolDef) {
	for _, sym := range symbols {
		name := strings.ToLower(sym.Name)
		all := s.workspaceIndex[name]
		if len(all) == 0 {
			continue
		}
		next := all[:0]
		for _, ex := range all {
			if ex.URI == sym.URI && ex.Range == sym.Range {
				continue
			}
			next = append(next, ex)
		}
		if len(next) == 0 {
			delete(s.workspaceIndex, name)
			continue
		}
		s.workspaceIndex[name] = next
	}
}

func bestLocalMatch(symbols []symbolDef, word string, pos lspPosition) *symbolDef {
	word = strings.ToLower(word)
	var best *symbolDef
	bestDist := 1<<31 - 1
	for i := range symbols {
		s := &symbols[i]
		if strings.ToLower(s.Name) != word {
			continue
		}
		if s.Range.Start.Line > pos.Line || (s.Range.Start.Line == pos.Line && s.Range.Start.Character > pos.Character) {
			continue
		}
		dist := (pos.Line-s.Range.Start.Line)*10000 + (pos.Character - s.Range.Start.Character)
		if best == nil || dist < bestDist {
			best = s
			bestDist = dist
		}
	}
	if best != nil {
		return best
	}
	return bestGlobalMatch(symbols, word)
}

func bestGlobalMatch(symbols []symbolDef, word string) *symbolDef {
	word = strings.ToLower(word)
	var cands []symbolDef
	for _, s := range symbols {
		if strings.ToLower(s.Name) == word {
			cands = append(cands, s)
		}
	}
	if len(cands) == 0 {
		return nil
	}
	sort.Slice(cands, func(i, j int) bool {
		if cands[i].URI != cands[j].URI {
			return cands[i].URI < cands[j].URI
		}
		if cands[i].Range.Start.Line != cands[j].Range.Start.Line {
			return cands[i].Range.Start.Line < cands[j].Range.Start.Line
		}
		return cands[i].Range.Start.Character < cands[j].Range.Start.Character
	})
	return &cands[0]
}

func wordAtPosition(src string, pos lspPosition) (string, string) {
	lines := strings.Split(src, "\n")
	if pos.Line < 0 || pos.Line >= len(lines) {
		return "", ""
	}
	line := lines[pos.Line]
	if pos.Character < 0 {
		return "", line
	}
	col := pos.Character
	if col >= len(line) {
		col = len(line) - 1
	}
	if col < 0 {
		return "", line
	}
	matches := identRE.FindAllStringIndex(line, -1)
	for _, m := range matches {
		if col >= m[0] && col <= m[1] {
			return strings.ToLower(line[m[0]:m[1]]), line
		}
	}
	return "", line
}

func rangeForWord(line int, word string, sourceLine string) lspRange {
	idx := strings.Index(strings.ToLower(sourceLine), strings.ToLower(word))
	if idx < 0 {
		idx = 0
	}
	return lspRange{
		Start: lspPosition{Line: line, Character: idx},
		End:   lspPosition{Line: line, Character: idx + len(word)},
	}
}

func uriToPath(uri string) string {
	const p = "file:///"
	if strings.HasPrefix(strings.ToLower(uri), strings.ToLower(p)) {
		return filepath.FromSlash(uri[len(p):])
	}
	return ""
}

func toLSPDiagnostic(d Diagnostic) lspDiagnostic {
	line, col, msg := parseDiagnosticMessage(d.Message)
	sev := 1
	if d.Severity == SeverityWarn {
		sev = 2
	}
	startLine := max(line-1, 0)
	startCol := max(col-1, 0)
	return lspDiagnostic{
		Severity: sev,
		Message:  msg,
		Range: lspRange{
			Start: lspPosition{Line: startLine, Character: startCol},
			End:   lspPosition{Line: startLine, Character: startCol + 1},
		},
	}
}

func parseDiagnosticMessage(s string) (line int, col int, msg string) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) == 3 {
		if l, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
			if c, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
				return l, c, strings.TrimSpace(parts[2])
			}
		}
	}
	return 1, 1, strings.TrimSpace(s)
}

func readLSPPacket(r *bufio.Reader) ([]byte, error) {
	contentLength := -1
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			nv := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			n, err := strconv.Atoi(nv)
			if err != nil {
				return nil, fmt.Errorf("bad content-length: %w", err)
			}
			contentLength = n
		}
	}
	if contentLength < 0 {
		return nil, fmt.Errorf("missing content-length")
	}
	buf := make([]byte, contentLength)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func writeResponse(w io.Writer, resp rpcResponse) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return writePacket(w, b)
}

func writeNotification(w io.Writer, n rpcNotification) error {
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return writePacket(w, b)
}

func writePacket(w io.Writer, payload []byte) error {
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(payload))
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	_, err := io.Copy(w, bytes.NewReader(payload))
	return err
}

func decodeID(raw json.RawMessage) interface{} {
	if len(raw) == 0 {
		return nil
	}
	var n int
	if err := json.Unmarshal(raw, &n); err == nil {
		return n
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
