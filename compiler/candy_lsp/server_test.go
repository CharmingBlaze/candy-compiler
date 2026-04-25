package candy_lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseDiagnosticMessage(t *testing.T) {
	l, c, m := parseDiagnosticMessage("12:34: bad thing")
	if l != 12 || c != 34 || m != "bad thing" {
		t.Fatalf("unexpected parse: %d %d %q", l, c, m)
	}
}

func TestToLSPDiagnosticDefaults(t *testing.T) {
	d := toLSPDiagnostic(Diagnostic{Severity: SeverityWarn, Message: "oops"})
	if d.Range.Start.Line != 0 || d.Range.Start.Character != 0 {
		t.Fatalf("unexpected default range: %#v", d.Range)
	}
	if d.Severity != 2 {
		t.Fatalf("severity mismatch: %d", d.Severity)
	}
}

func TestServerLifecycleAndPublishDiagnostics(t *testing.T) {
	in := &bytes.Buffer{}
	writeReq := func(payload string) {
		io.WriteString(in, "Content-Length: "+itoa(len(payload))+"\r\n\r\n"+payload)
	}
	writeReq(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"initialized","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///a.candy","text":"return true + 1;"}}}`)
	writeReq(`{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"exit","params":{}}`)

	out := &bytes.Buffer{}
	if err := Run(bytes.NewReader(in.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	packets := readAllPackets(t, out.Bytes())
	if len(packets) < 3 {
		t.Fatalf("expected >=3 packets, got %d", len(packets))
	}
	joined := strings.Join(packets, "\n")
	if !strings.Contains(joined, `"method":"textDocument/publishDiagnostics"`) {
		t.Fatalf("missing diagnostics notification: %s", joined)
	}
	if !strings.Contains(joined, `"id":1`) || !strings.Contains(joined, `"id":2`) {
		t.Fatalf("missing initialize/shutdown responses: %s", joined)
	}
}

func TestServerDidClosePublishesEmptyDiagnostics(t *testing.T) {
	in := &bytes.Buffer{}
	writeReq := func(payload string) {
		io.WriteString(in, "Content-Length: "+itoa(len(payload))+"\r\n\r\n"+payload)
	}
	writeReq(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///a.candy","text":"return 1;"}}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didClose","params":{"textDocument":{"uri":"file:///a.candy"}}}`)
	writeReq(`{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"exit","params":{}}`)
	out := &bytes.Buffer{}
	if err := Run(bytes.NewReader(in.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	packets := strings.Join(readAllPackets(t, out.Bytes()), "\n")
	if !strings.Contains(packets, `"diagnostics":[]`) {
		t.Fatalf("expected empty diagnostics on close: %s", packets)
	}
}

func TestServerHoverAndDefinition(t *testing.T) {
	in := &bytes.Buffer{}
	writeReq := func(payload string) {
		io.WriteString(in, "Content-Length: "+itoa(len(payload))+"\r\n\r\n"+payload)
	}
	writeReq(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///a.candy","text":"val user = 1;\nreturn user;"}}}`)
	writeReq(`{"jsonrpc":"2.0","id":3,"method":"textDocument/hover","params":{"textDocument":{"uri":"file:///a.candy"},"position":{"line":1,"character":8}}}`)
	writeReq(`{"jsonrpc":"2.0","id":4,"method":"textDocument/definition","params":{"textDocument":{"uri":"file:///a.candy"},"position":{"line":1,"character":8}}}`)
	writeReq(`{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"exit","params":{}}`)
	out := &bytes.Buffer{}
	if err := Run(bytes.NewReader(in.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	packets := strings.Join(readAllPackets(t, out.Bytes()), "\n")
	if !strings.Contains(packets, `"id":3`) || !strings.Contains(strings.ToLower(packets), "identifier") {
		t.Fatalf("missing hover response: %s", packets)
	}
	if !strings.Contains(packets, `"id":4`) || !strings.Contains(packets, `"uri":"file:///a.candy"`) {
		t.Fatalf("missing definition response: %s", packets)
	}
}

func TestServerMethodNotFound(t *testing.T) {
	in := &bytes.Buffer{}
	writeReq := func(payload string) {
		io.WriteString(in, "Content-Length: "+itoa(len(payload))+"\r\n\r\n"+payload)
	}
	writeReq(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","id":99,"method":"wat/nope","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"exit","params":{}}`)
	out := &bytes.Buffer{}
	if err := Run(bytes.NewReader(in.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	packets := strings.Join(readAllPackets(t, out.Bytes()), "\n")
	if !strings.Contains(packets, `"id":99`) || !strings.Contains(packets, `"code":-32601`) {
		t.Fatalf("expected method-not-found error response: %s", packets)
	}
}

func TestServerDedupesUnchangedDiagnostics(t *testing.T) {
	in := &bytes.Buffer{}
	writeReq := func(payload string) {
		io.WriteString(in, "Content-Length: "+itoa(len(payload))+"\r\n\r\n"+payload)
	}
	writeReq(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///a.candy","text":"return true + 1;"}}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didChange","params":{"textDocument":{"uri":"file:///a.candy","version":2},"contentChanges":[{"text":"return true + 1;"}]}}`)
	writeReq(`{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"exit","params":{}}`)
	out := &bytes.Buffer{}
	if err := Run(bytes.NewReader(in.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	packets := readAllPackets(t, out.Bytes())
	diagCount := 0
	for _, p := range packets {
		if strings.Contains(p, `"method":"textDocument/publishDiagnostics"`) {
			diagCount++
		}
	}
	if diagCount != 1 {
		t.Fatalf("expected 1 publishDiagnostics due to dedupe, got %d", diagCount)
	}
}

func TestServerDefinition_ImportedFile(t *testing.T) {
	root := filepath.ToSlash(t.TempDir())
	aURI := "file:///" + root + "/a.candy"
	bURI := "file:///" + root + "/b.candy"
	in := &bytes.Buffer{}
	writeReq := func(payload string) {
		io.WriteString(in, "Content-Length: "+itoa(len(payload))+"\r\n\r\n"+payload)
	}
	writeReq(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"` + bURI + `","text":"val fromb = 7;"}}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"` + aURI + `","text":"import \"b.candy\";\nreturn fromb;"}}}`)
	writeReq(`{"jsonrpc":"2.0","id":11,"method":"textDocument/definition","params":{"textDocument":{"uri":"` + aURI + `"},"position":{"line":1,"character":8}}}`)
	writeReq(`{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"exit","params":{}}`)
	out := &bytes.Buffer{}
	if err := Run(bytes.NewReader(in.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	defURI, ok := firstDefinitionURI(readAllPackets(t, out.Bytes()), 11)
	if !ok || defURI != bURI {
		t.Fatalf("expected definition in imported file %q, got %q (ok=%v)", bURI, defURI, ok)
	}
}

func TestServerDefinition_DeterministicWorkspaceFallback(t *testing.T) {
	root := filepath.ToSlash(t.TempDir())
	aURI := "file:///" + root + "/a.candy"
	bURI := "file:///" + root + "/b.candy"
	cURI := "file:///" + root + "/c.candy"
	in := &bytes.Buffer{}
	writeReq := func(payload string) {
		io.WriteString(in, "Content-Length: "+itoa(len(payload))+"\r\n\r\n"+payload)
	}
	writeReq(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"` + cURI + `","text":"return xref;"}}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"` + bURI + `","text":"val xref = 2;"}}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"` + aURI + `","text":"val xref = 1;"}}}`)
	writeReq(`{"jsonrpc":"2.0","id":12,"method":"textDocument/definition","params":{"textDocument":{"uri":"` + cURI + `"},"position":{"line":0,"character":7}}}`)
	writeReq(`{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"exit","params":{}}`)
	out := &bytes.Buffer{}
	if err := Run(bytes.NewReader(in.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	defURI, ok := firstDefinitionURI(readAllPackets(t, out.Bytes()), 12)
	if !ok || defURI != aURI {
		t.Fatalf("expected deterministic first URI %q, got %q (ok=%v)", aURI, defURI, ok)
	}
}

func TestServerDefinition_DidChangeReindexesTarget(t *testing.T) {
	root := filepath.ToSlash(t.TempDir())
	aURI := "file:///" + root + "/a.candy"
	bURI := "file:///" + root + "/b.candy"
	in := &bytes.Buffer{}
	writeReq := func(payload string) {
		io.WriteString(in, "Content-Length: "+itoa(len(payload))+"\r\n\r\n"+payload)
	}
	writeReq(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"` + aURI + `","text":"val fromx = 1;"}}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"` + bURI + `","text":"return fromx;"}}}`)
	writeReq(`{"jsonrpc":"2.0","id":13,"method":"textDocument/definition","params":{"textDocument":{"uri":"` + bURI + `"},"position":{"line":0,"character":8}}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didChange","params":{"textDocument":{"uri":"` + aURI + `","version":2},"contentChanges":[{"text":"val fromx = 99;"}]}}`)
	writeReq(`{"jsonrpc":"2.0","id":14,"method":"textDocument/definition","params":{"textDocument":{"uri":"` + bURI + `"},"position":{"line":0,"character":8}}}`)
	writeReq(`{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"exit","params":{}}`)
	out := &bytes.Buffer{}
	if err := Run(bytes.NewReader(in.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	packets := readAllPackets(t, out.Bytes())
	uri1, ok1 := firstDefinitionURI(packets, 13)
	uri2, ok2 := firstDefinitionURI(packets, 14)
	if !ok1 || !ok2 || uri1 != aURI || uri2 != aURI {
		t.Fatalf("expected both definitions to resolve to %q, got first=%q second=%q", aURI, uri1, uri2)
	}
}

func TestServerDefinition_DidCloseRemovesIndexSymbols(t *testing.T) {
	root := filepath.ToSlash(t.TempDir())
	aURI := "file:///" + root + "/a.candy"
	bURI := "file:///" + root + "/b.candy"
	in := &bytes.Buffer{}
	writeReq := func(payload string) {
		io.WriteString(in, "Content-Length: "+itoa(len(payload))+"\r\n\r\n"+payload)
	}
	writeReq(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"` + aURI + `","text":"val gone = 1;"}}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"` + bURI + `","text":"return gone;"}}}`)
	writeReq(`{"jsonrpc":"2.0","id":15,"method":"textDocument/definition","params":{"textDocument":{"uri":"` + bURI + `"},"position":{"line":0,"character":8}}}`)
	writeReq(`{"jsonrpc":"2.0","method":"textDocument/didClose","params":{"textDocument":{"uri":"` + aURI + `"}}}`)
	writeReq(`{"jsonrpc":"2.0","id":16,"method":"textDocument/definition","params":{"textDocument":{"uri":"` + bURI + `"},"position":{"line":0,"character":8}}}`)
	writeReq(`{"jsonrpc":"2.0","id":2,"method":"shutdown","params":{}}`)
	writeReq(`{"jsonrpc":"2.0","method":"exit","params":{}}`)
	out := &bytes.Buffer{}
	if err := Run(bytes.NewReader(in.Bytes()), out); err != nil {
		t.Fatal(err)
	}
	packets := readAllPackets(t, out.Bytes())
	if _, ok := firstDefinitionURI(packets, 15); !ok {
		t.Fatalf("expected definition before didClose")
	}
	if _, ok := firstDefinitionURI(packets, 16); ok {
		t.Fatalf("expected no definition after didClose")
	}
}

func readAllPackets(t *testing.T, b []byte) []string {
	t.Helper()
	r := bufio.NewReader(bytes.NewReader(b))
	var out []string
	for {
		p, err := readLSPPacket(r)
		if err == io.EOF {
			return out
		}
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, string(p))
	}
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}

func firstDefinitionURI(packets []string, id int) (string, bool) {
	for _, pkt := range packets {
		var raw map[string]json.RawMessage
		if err := json.Unmarshal([]byte(pkt), &raw); err != nil {
			continue
		}
		if !rawIDMatches(raw["id"], id) {
			continue
		}
		var locs []lspLocation
		if err := json.Unmarshal(raw["result"], &locs); err != nil {
			return "", false
		}
		if len(locs) == 0 {
			return "", false
		}
		return locs[0].URI, true
	}
	return "", false
}

func rawIDMatches(raw json.RawMessage, id int) bool {
	if len(raw) == 0 {
		return false
	}
	var n int
	if err := json.Unmarshal(raw, &n); err != nil {
		return false
	}
	return n == id
}
