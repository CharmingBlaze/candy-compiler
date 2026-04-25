package candy_parser

import (
	"candy/candy_ast"
	"candy/candy_token"
	"strings"
)

func (p *Parser) parseStatement() candy_ast.Statement {
	var attrs []*candy_ast.Attribute
	if p.curTokenIs(candy_token.LBRACK) {
		attrs = p.parseAttributes()
	}
	if p.curTokenIs(candy_token.IDENT) && p.curToken.Literal == "library" && p.peekTokenIs(candy_token.STR) {
		return p.parseLibraryStatement()
	}

	switch p.curToken.Type {
	case candy_token.EXPORT:
		p.nextToken()
		st := p.parseCStyleFunctionDeclaration(false, true, attrs)
		if st == nil {
			return nil
		}
		return st
	case candy_token.MODULE:
		return p.parseModuleStatement()
	case candy_token.ENUM:
		return p.parseEnumStatement()
	case candy_token.TRY:
		return p.parseTryStatement()
	case candy_token.RUN:
		return p.parseRunStatement()
	case candy_token.MAYBE:
		st := p.parseTypeLedStatement()
		if st != nil {
			if vs, ok := st.(*candy_ast.VarStatement); ok {
				vs.Attributes = attrs
			}
		}
		return st
	case candy_token.VAL, candy_token.DIM:
		st := p.parseValStatement()
		if st != nil {
			st.Attributes = attrs
		}
		return st
	case candy_token.VAR:
		st := p.parseVarStatement()
		if st != nil {
			st.Attributes = attrs
		}
		return st
	case candy_token.CONST:
		return p.parseConstStatement()
	case candy_token.PRINT:
		return p.parsePrintStatement()
	case candy_token.RETURN:
		return p.parseReturnStatement()
	case candy_token.FUNCTION:
		st := p.parseFunctionStatement()
		if st != nil {
			st.Attributes = attrs
		}
		return st
	case candy_token.SUB:
		st := p.parseSubStatement()
		if st != nil {
			st.Attributes = attrs
		}
		return st
	case candy_token.FOR:
		return p.parseForStatement()
	case candy_token.WHILE:
		return p.parseWhileStatement()
	case candy_token.DO:
		return p.parseDoWhileStatement()
	case candy_token.SWITCH:
		return p.parseSwitchStatement()
	case candy_token.IF:
		return p.parseIfWrapper()
	case candy_token.IMPORT:
		return p.parseImportStatement()
	case candy_token.STRUCT:
		st := p.parseStructStatement()
		if st != nil {
			st.Attributes = attrs
		}
		return st
	case candy_token.PACKAGE:
		return p.parsePackageStatement()
	case candy_token.CLASS:
		return p.parseClassStatement(false)
	case candy_token.SEALED:
		if p.peekTokenIs(candy_token.CLASS) {
			p.nextToken()
			return p.parseClassStatement(true)
		}
		return nil
	case candy_token.OBJECT:
		return p.parseObjectStatement()
	case candy_token.INTERFACE:
		return p.parseInterfaceStatement()
	case candy_token.TRAIT:
		return p.parseTraitStatement()
	case candy_token.DEFER:
		return p.parseDeferStatement()
	case candy_token.EXTERN:
		return p.parseExternFunctionStatement()
	case candy_token.SUSPEND:
		return p.parseSuspendFunctionStatement()
	case candy_token.FOREACH:
		return p.parseForEachStatement()
	case candy_token.REPEAT:
		return p.parseRepeatStatement()
	case candy_token.LOOP:
		return p.parseLoopStatement()
	case candy_token.BREAK:
		return p.parseBreakStatement()
	case candy_token.CONTINUE:
		return p.parseContinueStatement()
	case candy_token.DELETE:
		return p.parseDeleteStatement()
	case candy_token.WITH:
		return p.parseWithStatement()
	case candy_token.ASYNC:
		p.nextToken()
		if p.curTokenIs(candy_token.FUNCTION) {
			st := p.parseFunctionStatement()
			if st != nil {
				st.IsAsync = true
				st.Attributes = attrs
			}
			return st
		}
		st2 := p.parseCStyleFunctionDeclaration(true, false, attrs)
		if st2 == nil {
			return nil
		}
		return st2
	case candy_token.OVERRIDE:
		p.nextToken()
		if p.curTokenIs(candy_token.FUNCTION) {
			st := p.parseFunctionStatement()
			if st != nil {
				st.IsOverride = true
				st.Attributes = attrs
			}
			return st
		}
		return nil
	case candy_token.IDENT:
		if p.peekTokenIs(candy_token.AS) {
			return p.parseAsTypeDeclaration()
		}
		if p.peekTokenIs(candy_token.COLON) {
			return p.parseColonTypeDeclaration()
		}
		if p.peekTokenIs(candy_token.IDENT) || p.peekTokenIs(candy_token.ASTERISK) {
			st := p.parseTypeLedStatement()
			if st != nil {
				if vs, ok := st.(*candy_ast.VarStatement); ok {
					vs.Attributes = attrs
				}
			}
			return st
		}
		return p.parseExpressionStatement()
	case candy_token.SEMICOLON:
		return nil
	case candy_token.RBRACE, candy_token.EOF, candy_token.END, candy_token.WEND, candy_token.NEXT:
		return nil
	default:
		if p.curTokenIs(candy_token.VOID) || p.looksLikeTypeLedStatement() {
			return p.parseTypeLedStatement()
		}
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBlockStatement() *candy_ast.BlockStatement {
	b := &candy_ast.BlockStatement{Token: p.curToken}
	p.nextToken()
	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if st := p.parseStatement(); st != nil {
			b.Statements = append(b.Statements, st)
		} else {
			if p.curTokenIs(candy_token.RBRACE) || p.curTokenIs(candy_token.EOF) {
				break
			}
			p.nextToken()
		}
	}
	if p.curTokenIs(candy_token.EOF) {
		return b
	}
	_ = p.expect(candy_token.RBRACE)
	return b
}

func (p *Parser) parseBasicBlockStatement(blockType candy_token.TokenType) *candy_ast.BlockStatement {
	b := &candy_ast.BlockStatement{Token: p.curToken}

	// If we are at a newline (inserted semicolon), skip it to start the body
	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}

	for !p.curTokenIs(candy_token.END) && !p.curTokenIs(candy_token.EOF) && !p.curTokenIs(candy_token.WEND) && !p.curTokenIs(candy_token.NEXT) {
		if st := p.parseStatement(); st != nil {
			b.Statements = append(b.Statements, st)
		} else {
			if p.curTokenIs(candy_token.END) || p.curTokenIs(candy_token.EOF) || p.curTokenIs(candy_token.WEND) || p.curTokenIs(candy_token.NEXT) {
				break
			}
			p.nextToken()
		}
	}

	if p.curTokenIs(candy_token.END) {
		p.nextToken() // skip END
		// Optional block type after END: IF, FUNCTION, SUB, etc.
		if p.curToken.Type == blockType || p.curTokenIs(candy_token.SUB) || p.curTokenIs(candy_token.IF) || p.curTokenIs(candy_token.FUNCTION) || p.curTokenIs(candy_token.FOR) {
			p.nextToken()
		}
	} else if p.curTokenIs(candy_token.WEND) || p.curTokenIs(candy_token.NEXT) {
		p.nextToken()
		if p.curTokenIs(candy_token.IDENT) {
			p.nextToken()
		}
	}

	return b
}

func (p *Parser) expectSemicolon() bool {
	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
		return true
	}
	// Satisfied by newline (ASI) or end of block/file.
	// But we must move past the last token of the previous expression/statement
	// to the next one (which might be the ASI semicolon or the next identifier).
	if !p.curTokenIs(candy_token.EOF) && !p.curTokenIs(candy_token.RBRACE) {
		p.nextToken()
	}
	return true
}

func (p *Parser) parseExpressionStatement() *candy_ast.ExpressionStatement {
	s := &candy_ast.ExpressionStatement{Token: p.curToken}
	s.Expression = p.parseExpression(LOWEST)
	if !p.expectSemicolon() {
		return nil
	}
	return s
}

func (p *Parser) parseValStatement() *candy_ast.ValStatement {
	s := &candy_ast.ValStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekTokenIs(candy_token.COLON) || p.peekTokenIs(candy_token.AS) {
		p.nextToken()
		p.nextToken()
		s.TypeName = p.parseTypeIdentifier()
	}
	if !p.expectPeek(candy_token.ASSIGN) {
		return nil
	}
	p.nextToken()
	s.Value = p.parseExpression(LOWEST)
	if !p.expectSemicolon() {
		return nil
	}
	return s
}

func (p *Parser) parseVarStatement() *candy_ast.VarStatement {
	s := &candy_ast.VarStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekTokenIs(candy_token.COLON) {
		p.nextToken()
		p.nextToken()
		s.TypeName = p.parseTypeIdentifier()
	}
	if !p.expectPeek(candy_token.ASSIGN) {
		return nil
	}
	p.nextToken()
	s.Value = p.parseExpression(LOWEST)
	if !p.expectSemicolon() {
		return nil
	}
	return s
}

func (p *Parser) parseConstStatement() *candy_ast.ValStatement {
	s := &candy_ast.ValStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(candy_token.ASSIGN) {
		return nil
	}
	p.nextToken()
	s.Value = p.parseExpression(LOWEST)
	if !p.expectSemicolon() {
		return nil
	}
	return s
}

func (p *Parser) parsePrintStatement() *candy_ast.ExpressionStatement {
	tok := p.curToken
	p.nextToken() // past PRINT
	arg := p.parseExpression(LOWEST)
	if arg == nil {
		return nil
	}
	if !p.expectSemicolon() {
		return nil
	}
	// Lower to call so interpreter and LLVM see a normal call.
	return &candy_ast.ExpressionStatement{
		Token: tok,
		Expression: &candy_ast.CallExpression{
			Token:     tok,
			Function:  &candy_ast.Identifier{Token: tok, Value: "println"},
			Arguments: []candy_ast.Expression{arg},
		},
	}
}

func (p *Parser) parseReturnStatement() *candy_ast.ReturnStatement {
	s := &candy_ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	if p.curTokenIs(candy_token.SEMICOLON) {
		_ = p.expect(candy_token.SEMICOLON)
		return s
	}
	s.ReturnValue = p.parseExpression(LOWEST)
	if !p.expectSemicolon() {
		return nil
	}
	return s
}

func (p *Parser) parseIfWrapper() *candy_ast.IfExpression {
	e := p.parseIfExpression()
	if e == nil {
		return nil
	}
	ie, ok := e.(*candy_ast.IfExpression)
	if !ok {
		return nil
	}
	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}
	return ie
}

func (p *Parser) looksLikeTypeLedStatement() bool {
	if p.curTokenIs(candy_token.VOID) {
		return p.peekTokenIs(candy_token.IDENT)
	}
	if p.curTokenIs(candy_token.REF) || p.curTokenIs(candy_token.SHARED) {
		return p.peekTokenIs(candy_token.IDENT)
	}
	if !p.curTokenIs(candy_token.IDENT) {
		return false
	}
	return p.peekTokenIs(candy_token.IDENT)
}

func (p *Parser) parseTypeLedStatement() candy_ast.Statement {
	isMaybe := p.curTokenIs(candy_token.MAYBE)
	if isMaybe {
		p.nextToken()
	}
	isRef := p.curTokenIs(candy_token.REF)
	if isRef {
		p.nextToken()
	}
	isShared := p.curTokenIs(candy_token.SHARED)
	if isShared {
		p.nextToken()
	}

	typeName := p.parseTypeIdentifier()
	if typeName == nil {
		return nil
	}
	p.nextToken()

	if !p.curTokenIs(candy_token.IDENT) {
		p.addErrorf("expected identifier after type, got %s", p.curToken.Type)
		return nil
	}
	name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(candy_token.LPAREN) {
		return p.parseCStyleFunctionStatement(typeName, name)
	}

	// Support C-style arrays: `int scores[10]`
	if p.peekTokenIs(candy_token.LBRACK) {
		p.nextToken() // into [
		// For now, we'll just parse the size expression
		p.nextToken()
		// _ = p.parseExpression(LOWEST)
		if !p.expect(candy_token.RBRACK) {
			return nil
		}
		// We'll mark it as an array type in typeName or similar.
		if id, ok := typeName.(*candy_ast.Identifier); ok {
			id.Value += "[]"
		}
	}

	s := &candy_ast.VarStatement{
		Token:    candy_ast.ExprToken(typeName),
		Name:     name,
		TypeName: typeName,
		IsMaybe:  isMaybe,
		IsRef:    isRef,
		IsShared: isShared,
	}

	if p.peekTokenIs(candy_token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		s.Value = p.parseExpression(LOWEST)
	}

	if !p.expectSemicolon() {
		return nil
	}
	return s
}

func (p *Parser) parseTypedVarStatement(typeName, name *candy_ast.Identifier) *candy_ast.VarStatement {
	s := &candy_ast.VarStatement{Token: typeName.Token, Name: name, TypeName: typeName}
	if !p.expectPeek(candy_token.ASSIGN) {
		return nil
	}
	p.nextToken()
	s.Value = p.parseExpression(LOWEST)
	if !p.expectSemicolon() {
		return nil
	}
	return s
}

func (p *Parser) parseAsTypeDeclaration() candy_ast.Statement {
	name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // past name
	p.nextToken() // past as
	typeName := p.parseTypeIdentifier()
	if typeName == nil {
		return nil
	}
	p.nextToken()
	s := &candy_ast.VarStatement{
		Token:    name.Token,
		Name:     name,
		TypeName: typeName,
	}
	if p.curTokenIs(candy_token.ASSIGN) {
		p.nextToken()
		s.Value = p.parseExpression(LOWEST)
	}
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseColonTypeDeclaration() candy_ast.Statement {
	name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // past name
	p.nextToken() // past :
	typeName := p.parseTypeIdentifier()
	if typeName == nil {
		return nil
	}
	p.nextToken()
	s := &candy_ast.VarStatement{
		Token:    name.Token,
		Name:     name,
		TypeName: typeName,
	}
	if p.curTokenIs(candy_token.ASSIGN) {
		p.nextToken()
		s.Value = p.parseExpression(LOWEST)
	}
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseCStyleFunctionStatement(returnType candy_ast.Expression, name *candy_ast.Identifier) *candy_ast.FunctionStatement {
	s := &candy_ast.FunctionStatement{
		Token:      candy_ast.ExprToken(returnType),
		Name:       name,
		ReturnType: returnType,
	}
	if !p.expectPeek(candy_token.LPAREN) {
		return nil
	}
	p.nextToken()
	s.Parameters = p.parseCStyleFunctionParameters()
	if !p.curTokenIs(candy_token.LBRACE) {
		if !p.expect(candy_token.LBRACE) {
			return nil
		}
	}
	s.Body = p.parseBlockStatement()
	return s
}

// parseCStyleFunctionDeclaration parses `T name ( params ) { }` (caller has already consumed
// `export` / the leading `async` if applicable). isExported tags `export` declarations.
func (p *Parser) parseCStyleFunctionDeclaration(isAsync, isExported bool, attrs []*candy_ast.Attribute) *candy_ast.FunctionStatement {
	typeName := p.parseTypeIdentifier()
	if typeName == nil {
		return nil
	}
	p.nextToken()
	if !p.curTokenIs(candy_token.IDENT) {
		p.addErrorf("expected function name, got %s", p.curToken.Type)
		return nil
	}
	name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	s := p.parseCStyleFunctionStatement(typeName, name)
	if s != nil {
		s.IsAsync = isAsync
		s.Exported = isExported
		s.Attributes = attrs
	}
	return s
}

func (p *Parser) parseTypedDeclStatement() candy_ast.Statement {
	s := &candy_ast.VarStatement{Token: p.curToken}
	s.TypeName = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // move to name
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(candy_token.ASSIGN) {
		p.nextToken() // move to =
		p.nextToken() // move past =
		s.Value = p.parseExpression(LOWEST)
	}

	if !p.expectSemicolon() {
		return nil
	}
	return s
}

func (p *Parser) parseFunctionStatement() *candy_ast.FunctionStatement {
	s := &candy_ast.FunctionStatement{Token: p.curToken}

	// If the caller already set IsAsync or IsOverride, they'll be preserved if we manage them correctly.
	// But usually they are consumed by parseStatement and passed here somehow?
	// Actually, the current parseStatement logic:
	/*
		case candy_token.ASYNC:
			p.nextToken()
			if p.curTokenIs(candy_token.FUNCTION) {
				st := p.parseFunctionStatement()
				if st != nil {
					st.IsAsync = true
				}
				return st
			}
	*/
	// This is fine for now.

	// `function` — only this spelling uses the "BASIC" function body rules in a few
	// branches; `fun` and `func` are the same FUNCTION token and work for all normal
	// signatures, including `fun f() { }` and `fun f(): int { }`.
	isBasicFunction := strings.ToLower(p.curToken.Literal) == "function"
	p.nextToken() // past fun/function/sub

	var recv *candy_ast.Parameter
	if p.curTokenIs(candy_token.LPAREN) {
		recv = p.parseReceiverParameter()
		if recv == nil {
			return nil
		}
	}

	// `fun` name can be a keyword like `maybe` (lexed as MAYBE). After a receiver, `cur` is already
	// the name; otherwise the name is the next token.
	if p.curTokenIs(candy_token.IDENT) || p.curTokenIs(candy_token.MAYBE) {
		// keep cur
	} else if p.peekTokenIs(candy_token.IDENT) {
		p.nextToken()
	} else if p.peekTokenIs(candy_token.MAYBE) {
		p.nextToken()
	} else {
		p.addErrorf("expected function name, got %s", p.peekToken.Type)
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	for _, tp := range p.parseTypeParamsIfAny() {
		s.TypeParameters = append(s.TypeParameters, tp)
	}

	if !p.expectPeek(candy_token.LPAREN) {
		return nil
	}
	p.nextToken()
	s.Parameters = p.parseFunctionParameters()
	s.Receiver = recv

	if isBasicFunction {
		// BASIC style: return type optional or after block?
		// New spec: `int add(int a, int b)` or `function add(a, b)`
		// If it started with `function`, it might not have a return type explicitly.
	}

	// Optional return type: `int add(...)` vs `void greet(...)`
	// Wait, the current logic for `int add(...)` is handled by `parseTypeLedStatement`.
	// This `parseFunctionStatement` is for when it starts with `function` or `sub`.

	if p.curTokenIs(candy_token.COLON) || p.curTokenIs(candy_token.AS) {
		p.nextToken()
		s.ReturnType = p.parseTypeIdentifier()
		p.nextToken()
	} else if !isBasicFunction && p.curTokenIs(candy_token.IDENT) && p.peekTokenIs(candy_token.LBRACE) {
		// `fun` uses `: ReturnType` before `{`. Reject `fun f() Int {` (missing `:`) so we don't skip `Int`.
		p.addErrorf("expected `:` or `as` before return type in function signature")
		return nil
	}

	if p.curTokenIs(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else if p.expectPeek(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else {
		return nil
	}
	return s
}

func (p *Parser) parseSubStatement() *candy_ast.FunctionStatement {
	s := &candy_ast.FunctionStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(candy_token.LPAREN) {
		return nil
	}
	p.nextToken() // past (
	s.Parameters = p.parseFunctionParameters()

	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}
	if p.curTokenIs(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else if p.expectPeek(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else {
		return nil
	}
	return s
}

func (p *Parser) parseForStatement() candy_ast.Statement {
	if p.peekTokenIs(candy_token.LPAREN) {
		return p.parseCStyleForStatement()
	}
	s := &candy_ast.ForStatement{Token: p.curToken}
	if p.peekTokenIs(candy_token.EACH) {
		p.nextToken() // past EACH
	}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Var = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekTokenIs(candy_token.COMMA) {
		p.nextToken() // past Var
		p.nextToken() // past ,
		if !p.curTokenIs(candy_token.IDENT) {
			p.addErrorf("expected identifier after comma in for-each, got %s", p.curToken.Type)
			return nil
		}
		s.ValueVar = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}
	if p.peekTokenIs(candy_token.IN) {
		p.nextToken() // IN
		p.nextToken() // iterable expr start
		// Do not let `a {` after `in` become a struct literal; the `{` is the for body.
		p.exprStopAtBrace = true
		s.Iterable = p.parseExpression(LOWEST)
		p.exprStopAtBrace = false
		if !p.curTokenIs(candy_token.LBRACE) {
			if !p.expectPeek(candy_token.LBRACE) {
				return nil
			}
		}
		s.Body = p.parseBlockStatement()
		return s
	}
	if !p.expectPeek(candy_token.ASSIGN) {
		return nil
	}
	p.nextToken() // past =
	s.Start = p.parseExpression(LOWEST)
	if !p.expectPeek(candy_token.TO) {
		return nil
	}
	p.nextToken() // past TO
	p.exprStopAtBrace = true
	s.End = p.parseExpression(LOWEST)
	p.exprStopAtBrace = false
	if p.peekTokenIs(candy_token.STEP) {
		p.nextToken() // to STEP
		p.nextToken() // past STEP
		s.Step = p.parseExpression(LOWEST)
	}

	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}

	if p.curTokenIs(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else if p.expectPeek(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else {
		return nil
	}

	return s
}

func (p *Parser) parseCStyleForStatement() *candy_ast.CForStatement {
	s := &candy_ast.CForStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.LPAREN) {
		return nil
	}
	p.nextToken() // past (
	s.Init = p.parseStatement()
	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}

	s.Cond = p.parseExpression(LOWEST)
	if !p.expectSemicolon() {
		return nil
	}
	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}

	s.Post = p.parseExpression(LOWEST)

	if !p.expectPeek(candy_token.RPAREN) {
		return nil
	}

	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}

	if p.curTokenIs(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else if p.expectPeek(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else {
		return nil
	}

	return s
}

func (p *Parser) parseWhileStatement() *candy_ast.WhileStatement {
	s := &candy_ast.WhileStatement{Token: p.curToken}
	p.nextToken() // past while

	p.exprStopAtBrace = true
	s.Condition = p.parseExpression(LOWEST)
	p.exprStopAtBrace = false

	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}

	if p.curTokenIs(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else if p.expectPeek(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else {
		return nil
	}
	return s
}

func (p *Parser) parseDoWhileStatement() *candy_ast.DoWhileStatement {
	s := &candy_ast.DoWhileStatement{Token: p.curToken}
	p.nextToken() // past do

	if p.curTokenIs(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else if p.expectPeek(candy_token.LBRACE) {
		s.Body = p.parseBlockStatement()
	} else {
		return nil
	}

	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}

	if !p.expect(candy_token.WHILE) {
		return nil
	}
	if !p.expect(candy_token.LPAREN) {
		return nil
	}
	s.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(candy_token.RPAREN) {
		return nil
	}
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseSwitchStatement() *candy_ast.SwitchStatement {
	sw := &candy_ast.SwitchStatement{Token: p.curToken}
	p.nextToken() // past switch

	p.exprStopAtBrace = true
	sw.Subject = p.parseExpression(LOWEST)
	p.exprStopAtBrace = false

	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}
	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	p.nextToken() // into {

	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.RBRACE) {
			break
		}

		var c candy_ast.SwitchCase
		c.Token = p.curToken

		if p.curTokenIs(candy_token.DEFAULT) {
			c.IsDefault = true
			p.nextToken()
			if p.curTokenIs(candy_token.COLON) || p.curTokenIs(candy_token.ARROW) {
				p.nextToken()
			}
		} else if p.curTokenIs(candy_token.CASE) {
			p.nextToken()
			for {
				pat := p.parseExpression(LOWEST)
				if pat != nil {
					c.Patterns = append(c.Patterns, pat)
				}
				if p.peekTokenIs(candy_token.COMMA) {
					p.nextToken()
					p.nextToken()
					continue
				}
				break
			}
			if p.curTokenIs(candy_token.COLON) || p.curTokenIs(candy_token.ARROW) {
				p.nextToken()
			} else if p.expectPeek(candy_token.COLON) || p.expectPeek(candy_token.ARROW) {
				p.nextToken()
			}
		} else {
			// tolerate stray tokens
			p.nextToken()
			continue
		}

		if p.curTokenIs(candy_token.LBRACE) {
			c.Body = p.parseBlockStatement()
		} else {
			expr := p.parseExpression(LOWEST)
			esTok := c.Token
			if expr != nil {
				switch e := expr.(type) {
				case *candy_ast.Identifier:
					esTok = e.Token
				case *candy_ast.StringLiteral:
					esTok = e.Token
				case *candy_ast.IntegerLiteral:
					esTok = e.Token
				case *candy_ast.FloatLiteral:
					esTok = e.Token
				case *candy_ast.Boolean:
					esTok = e.Token
				case *candy_ast.NullLiteral:
					esTok = e.Token
				case *candy_ast.PrefixExpression:
					esTok = e.Token
				case *candy_ast.InfixExpression:
					esTok = e.Token
				case *candy_ast.CallExpression:
					esTok = e.Token
				case *candy_ast.DotExpression:
					esTok = e.Token
				case *candy_ast.IndexExpression:
					esTok = e.Token
				case *candy_ast.GroupedExpression:
					esTok = e.Token
				case *candy_ast.IfExpression:
					esTok = e.Token
				case *candy_ast.LambdaExpression:
					esTok = e.Token
				default:
					esTok = c.Token
				}
			}
			c.Body = &candy_ast.ExpressionStatement{Token: esTok, Expression: expr}
			_ = p.expectSemicolon()
		}

		sw.Cases = append(sw.Cases, c)
	}

	_ = p.expect(candy_token.RBRACE)
	_ = p.expectSemicolon()
	return sw
}

func (p *Parser) parseDeferStatement() *candy_ast.DeferStatement {
	s := &candy_ast.DeferStatement{Token: p.curToken}
	p.nextToken() // past defer

	expr := p.parseExpression(LOWEST)
	call, ok := expr.(*candy_ast.CallExpression)
	if !ok {
		p.addErrorf("defer must be followed by a function call")
		return nil
	}
	s.Call = call

	if !p.expectSemicolon() {
		return nil
	}
	return s
}

func (p *Parser) parseReceiverParameter() *candy_ast.Parameter {
	if !p.expect(candy_token.LPAREN) {
		return nil
	}
	if p.curTokenIs(candy_token.RPAREN) {
		return nil
	}
	if !p.curTokenIs(candy_token.IDENT) {
		p.addErrorf("expected receiver name, got %s", p.curToken.Type)
		return nil
	}
	param := candy_ast.Parameter{
		Token: p.curToken,
		Name:  &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	// Optional typed receiver: (self: Box) / (self as Box)
	if p.peekTokenIs(candy_token.COLON) || p.peekTokenIs(candy_token.AS) {
		p.nextToken() // : / as
		p.nextToken() // first type token
		param.TypeName = p.parseTypeIdentifier()
	}
	// Advance to ')' if we're still on name/type token.
	if !p.curTokenIs(candy_token.RPAREN) {
		p.nextToken()
	}
	if !p.curTokenIs(candy_token.RPAREN) {
		p.addErrorf("expected ), got %s", p.curToken.Type)
		return nil
	}
	p.nextToken() // move past ')'
	return &param
}

func (p *Parser) parseFunctionParameters() []candy_ast.Parameter {
	var r []candy_ast.Parameter
	// Caller advanced past `(`; cur is first param or `)`.
	if p.curTokenIs(candy_token.RPAREN) {
		p.nextToken()
		return r
	}
	for {
		r = append(r, p.parseOneParam())
		if p.curTokenIs(candy_token.RPAREN) {
			p.nextToken()
			return r
		}
		if p.curTokenIs(candy_token.COMMA) {
			p.nextToken()
			continue
		}
		p.addErrorf("expected , or ) in param list, got %s", p.curToken.Type)
		return r
	}
}

func (p *Parser) parseCStyleFunctionParameters() []candy_ast.Parameter {
	var r []candy_ast.Parameter
	if p.curTokenIs(candy_token.RPAREN) {
		p.nextToken()
		return r
	}
	for {
		var p1 candy_ast.Parameter
		first := p.curToken
		if p.peekTokenIs(candy_token.IDENT) {
			// TYPE IDENT
			typeId := &candy_ast.Identifier{Token: first, Value: first.Literal}
			p.nextToken()
			nameId := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			p1 = candy_ast.Parameter{Token: first, Name: nameId, TypeName: typeId}
			p.nextToken()
		} else {
			// IDENT (Dynamic)
			nameId := &candy_ast.Identifier{Token: first, Value: first.Literal}
			p1 = candy_ast.Parameter{Token: first, Name: nameId, TypeName: nameId}
			p.nextToken()
		}
		r = append(r, p1)

		if p.curTokenIs(candy_token.RPAREN) {
			p.nextToken()
			return r
		}
		if p.curTokenIs(candy_token.COMMA) {
			p.nextToken()
			continue
		}
		break
	}
	return r
}

func (p *Parser) parseOneParam() candy_ast.Parameter {
	n := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // move past name
	var tn candy_ast.Expression
	if p.curTokenIs(candy_token.COLON) || p.curTokenIs(candy_token.AS) {
		p.nextToken() // move past : or as
		tn = p.parseTypeIdentifier()
		if tn != nil {
			p.nextToken()
		}
	}
	var def candy_ast.Expression
	if p.curTokenIs(candy_token.ASSIGN) {
		p.nextToken() // move past =
		def = p.parseExpression(LOWEST)
	}
	return candy_ast.Parameter{Token: n.Token, Name: n, TypeName: tn, Default: def}
}

func (p *Parser) parseTypeIdentifier() candy_ast.Expression {
	if p.curTokenIs(candy_token.LPAREN) {
		tok := p.curToken
		p.nextToken() // past (
		var types []candy_ast.Expression
		for !p.curTokenIs(candy_token.RPAREN) && !p.curTokenIs(candy_token.EOF) {
			ty := p.parseTypeIdentifier()
			if ty != nil {
				types = append(types, ty)
			}
			p.nextToken()
			if p.curTokenIs(candy_token.COMMA) {
				p.nextToken()
			}
		}
		if !p.curTokenIs(candy_token.RPAREN) {
			p.addErrorf("expected ), got %s", p.curToken.Type)
			return nil
		}
		// cur is on )
		return &candy_ast.TupleTypeExpression{Token: tok, Types: types}
	}

	isMaybe := false
	if p.curTokenIs(candy_token.MAYBE) {
		isMaybe = true
		p.nextToken()
	}

	if !p.curTokenIs(candy_token.IDENT) && !p.curTokenIs(candy_token.VOID) {
		p.addErrorf("expected type identifier, got %s", p.curToken.Type)
		return nil
	}
	tok := p.curToken
	name := strings.ToLower(p.curToken.Literal)
	if isMaybe {
		name = "maybe " + name
	}

	isPointer := false
	if p.peekTokenIs(candy_token.ASTERISK) {
		p.nextToken()
		isPointer = true
		name += "*"
	}

	if p.peekTokenIs(candy_token.QUESTION) {
		p.nextToken()
		name += "?"
	}
	ident := &candy_ast.Identifier{Token: tok, Value: name, IsPointer: isPointer}

	if p.peekTokenIs(candy_token.LT) {
		p.nextToken() // <
		p.nextToken() // first arg
		var args []candy_ast.Expression
		for !p.curTokenIs(candy_token.GT) && !p.curTokenIs(candy_token.EOF) {
			args = append(args, p.parseTypeIdentifier())
			if p.peekTokenIs(candy_token.GT) {
				p.nextToken()
				break
			}
			if p.curTokenIs(candy_token.COMMA) {
				p.nextToken()
				continue
			}
			if !p.curTokenIs(candy_token.EOF) {
				p.nextToken()
			}
		}
		if !p.curTokenIs(candy_token.GT) {
			p.addErrorf("expected >, got %s", p.curToken.Type)
			return nil
		}
		return &candy_ast.TypeExpression{Token: tok, Name: ident, Arguments: args}
	}
	return ident
}

func (p *Parser) parseGenericParameters() []*candy_ast.Identifier {
	var params []*candy_ast.Identifier
	if !p.peekTokenIs(candy_token.LT) {
		return params
	}
	p.nextToken() // <
	p.nextToken() // first ident
	for !p.curTokenIs(candy_token.GT) && !p.curTokenIs(candy_token.EOF) {
		params = append(params, &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		p.nextToken()
		if p.curTokenIs(candy_token.COMMA) {
			p.nextToken()
		}
	}
	if !p.curTokenIs(candy_token.GT) {
		p.addErrorf("expected >, got %s", p.curToken.Type)
		return nil
	}
	return params
}

func (p *Parser) parseStructStatement() *candy_ast.StructStatement {
	s := &candy_ast.StructStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	s.TypeParameters = p.parseGenericParameters()

	if p.peekTokenIs(candy_token.COLON) {
		p.nextToken()
		for {
			if !p.expectPeek(candy_token.IDENT) {
				break
			}
			s.Bases = append(s.Bases, &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			if p.peekTokenIs(candy_token.COMMA) {
				p.nextToken()
				continue
			}
			break
		}
	}

	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	p.nextToken()

	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
			continue
		}

		memberAttrs := p.parseAttributesIfAny()

		if p.curTokenIs(candy_token.IDENT) && strings.EqualFold(p.curToken.Literal, "using") {
			p.nextToken()
			base := p.parseTypeIdentifier()
			if base != nil {
				if ident, ok := base.(*candy_ast.Identifier); ok {
					s.Bases = append(s.Bases, ident)
				}
			}
			_ = p.expectSemicolon()
			continue
		}

		// `name: Type? = Value?` struct fields (require `:` so `V operator+` is not parsed as a field name `V`)
		if p.curTokenIs(candy_token.IDENT) && p.peekTokenIs(candy_token.COLON) {
			fname := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			p.nextToken() // name
			p.nextToken() // :
			typ := p.parseTypeIdentifier()

			var init candy_ast.Expression
			if p.peekTokenIs(candy_token.ASSIGN) {
				p.nextToken() // =
				p.nextToken() // first expr token
				init = p.parseExpression(LOWEST)
			}

			_ = p.expectSemicolon()
			s.Fields = append(s.Fields, candy_ast.Field{
				Token:      fname.Token,
				Name:       fname,
				TypeName:   typ,
				Init:       init,
				Attributes: memberAttrs,
			})
			continue
		}

		isPrivate := false
		if p.curTokenIs(candy_token.PRIVATE) {
			isPrivate = true
			p.nextToken()
		}

		if p.curTokenIs(candy_token.OPERATOR) {
			if op := p.parseOperatorOverloadWithReturnType(nil); op != nil {
				s.Operators = append(s.Operators, op)
				continue
			}
			p.addErrorf("expected `operator+(` after leading `operator`")
		}

		// C#-style: `T operator+(...)` — return type, then the `operator` keyword, then the symbol.
		if p.curTokenIs(candy_token.IDENT) && p.peekTokenIs(candy_token.OPERATOR) {
			retT := p.parseTypeIdentifier()
			if retT == nil {
				p.nextToken()
			} else {
				p.nextToken() // to OPERATOR
				if op := p.parseOperatorOverloadWithReturnType(retT); op != nil {
					s.Operators = append(s.Operators, op)
					continue
				}
				p.addErrorf("expected `operator+(` after type+operator in struct")
			}
		}

		// `type name` (field, method, or property) or `type operator+` after a two-token lookahead.
		// `void` / `Maybe` (MAYBE) are valid type starts (e.g. `void init(...) { }`, `Maybe X?` via parseTypeIdentifier).
		if p.curTokenIs(candy_token.IDENT) || p.curTokenIs(candy_token.VOID) || p.curTokenIs(candy_token.MAYBE) {
			var entryType candy_ast.Expression
			if p.peekTokenIs(candy_token.IDENT) || p.peekTokenIs(candy_token.ASTERISK) {
				entryType = p.parseTypeIdentifier()
				p.nextToken()
			}
			if p.curTokenIs(candy_token.OPERATOR) {
				if op := p.parseOperatorOverloadWithReturnType(entryType); op != nil {
					s.Operators = append(s.Operators, op)
					continue
				}
				p.addErrorf("expected operator method body")
			}
			if p.curTokenIs(candy_token.IDENT) {
				name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
				if p.peekTokenIs(candy_token.LPAREN) {
					if m := p.parseStructMethodAfterType(entryType); m != nil {
						s.Methods = append(s.Methods, m)
					}
				} else if p.peekTokenIs(candy_token.LBRACE) {
					if prop := p.parsePropertyStatementAfterName(entryType, name); prop != nil {
						s.Properties = append(s.Properties, prop)
					}
				} else {
					s.Fields = append(s.Fields, p.parseStructFields(entryType, name, isPrivate, memberAttrs)...)
				}
			}
		}

		if p.curTokenIs(candy_token.COMMA) || p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
		} else {
			p.nextToken() // Always advance to avoid infinite loop
		}
	}

	_ = p.expect(candy_token.RBRACE)
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseStructFields(typeExpr candy_ast.Expression, firstIdent *candy_ast.Identifier, isPrivate bool, fieldAttrs []*candy_ast.Attribute) []candy_ast.Field {
	var fields []candy_ast.Field
	addField := func(name *candy_ast.Identifier) {
		fields = append(fields, candy_ast.Field{
			Token:      name.Token,
			Attributes: fieldAttrs,
			IsPrivate:  isPrivate,
			Name:       name,
			TypeName:   typeExpr,
		})
	}
	addField(firstIdent)
	if p.peekTokenIs(candy_token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		if n := len(fields); n > 0 {
			fields[n-1].Init = p.parseExpression(LOWEST)
		}
	}
	for p.peekTokenIs(candy_token.COMMA) {
		p.nextToken() // ,
		if !p.expectPeek(candy_token.IDENT) {
			break
		}
		name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		addField(name)
		if p.peekTokenIs(candy_token.ASSIGN) {
			p.nextToken()
			p.nextToken()
			if n := len(fields); n > 0 {
				fields[n-1].Init = p.parseExpression(LOWEST)
			}
		}
	}
	p.nextToken()
	return fields
}

func (p *Parser) parsePropertyStatementAfterName(typeName candy_ast.Expression, name *candy_ast.Identifier) *candy_ast.PropertyStatement {
	prop := &candy_ast.PropertyStatement{
		Token: name.Token,
		Name:  name,
		Type:  typeName,
	}
	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	p.nextToken()

	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		switch p.curToken.Literal {
		case "get":
			p.nextToken()
			if p.curTokenIs(candy_token.LBRACE) {
				prop.Getter = p.parseBlockStatement()
			} else if p.curTokenIs(candy_token.SEMICOLON) {
				prop.IsAuto = true
				p.nextToken()
			}
		case "set":
			p.nextToken()
			if p.curTokenIs(candy_token.LBRACE) {
				prop.Setter = p.parseBlockStatement()
			} else if p.curTokenIs(candy_token.SEMICOLON) {
				prop.IsAuto = true
				p.nextToken()
			}
		default:
			p.nextToken()
		}
	}
	_ = p.expect(candy_token.RBRACE)
	// `} = "default"` — after `}` the lexer leaves `cur` on `=`.
	if p.curTokenIs(candy_token.ASSIGN) {
		p.nextToken()
		prop.DefaultValue = p.parseExpression(LOWEST)
		// parseExpression can leave cur on the last rvalue token; advance for the next struct member.
		if prop.DefaultValue != nil {
			p.nextToken()
		}
	}
	return prop
}

func (p *Parser) parseAttributesIfAny() []*candy_ast.Attribute {
	if p.curTokenIs(candy_token.LBRACK) {
		return p.parseAttributes()
	}
	return nil
}

func (p *Parser) parseOperatorOverloadStatement() *candy_ast.OperatorOverloadStatement {
	return p.parseOperatorOverloadWithReturnType(nil)
}

// parseOperatorOverloadWithReturnType reads `operator` ... `(` C-style formals `)`.
// For a bare `operator+(...)` in a struct, ret may be nil.
func (p *Parser) parseOperatorOverloadWithReturnType(ret candy_ast.Expression) *candy_ast.OperatorOverloadStatement {
	op := &candy_ast.OperatorOverloadStatement{Token: p.curToken, ReturnType: ret}
	p.nextToken() // past operator

	// operator symbol: +, -, *, /, or multi-char tokens; skip optional whitespace tokenization detail.
	op.Operator = p.curToken.Literal
	p.nextToken()

	if !p.expect(candy_token.LPAREN) {
		return nil
	}
	op.Parameters = p.parseCStyleFunctionParameters()
	if !p.curTokenIs(candy_token.LBRACE) {
		if !p.expect(candy_token.LBRACE) {
			return nil
		}
	}
	op.Body = p.parseBlockStatement()
	return op
}

func (p *Parser) parseStructField() candy_ast.Field {
	var f candy_ast.Field
	if !p.curTokenIs(candy_token.IDENT) {
		return f
	}
	f.Token = p.curToken
	f.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(candy_token.COLON) {
		return f
	}
	p.nextToken()
	tn := p.parseTypeIdentifier()
	if tn == nil {
		return f
	}
	f.TypeName = tn
	p.nextToken()
	return f
}

func (p *Parser) parseStructTypeLedFields(typeName *candy_ast.Identifier) []candy_ast.Field {
	var fields []candy_ast.Field
	for p.curTokenIs(candy_token.IDENT) {
		name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		fields = append(fields, candy_ast.Field{
			Token:    name.Token,
			Name:     name,
			TypeName: typeName,
		})
		p.nextToken()
		if p.curTokenIs(candy_token.COMMA) {
			p.nextToken()
			continue
		}
		break
	}
	return fields
}

func (p *Parser) parseStructMethodAfterType(returnType candy_ast.Expression) *candy_ast.FunctionStatement {
	if returnType == nil {
		return nil
	}
	name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	fn := &candy_ast.FunctionStatement{
		Token:      candy_ast.ExprToken(returnType),
		Name:       name,
		ReturnType: returnType,
	}
	if !p.expectPeek(candy_token.LPAREN) {
		return nil
	}
	p.nextToken()
	fn.Parameters = p.parseCStyleFunctionParameters()
	if !p.curTokenIs(candy_token.LBRACE) {
		if !p.expect(candy_token.LBRACE) {
			return nil
		}
	}
	fn.Body = p.parseBlockStatement()
	return fn
}

func (p *Parser) parseStructCtorLike(structName *candy_ast.Identifier) *candy_ast.FunctionStatement {
	if !p.curTokenIs(candy_token.IDENT) || p.curToken.Literal != structName.Value {
		return nil
	}
	fn := &candy_ast.FunctionStatement{
		Token:      p.curToken,
		Name:       &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		ReturnType: structName,
	}
	if !p.expectPeek(candy_token.LPAREN) {
		return nil
	}
	p.nextToken()
	fn.Parameters = p.parseCStyleFunctionParameters()
	if p.curTokenIs(candy_token.COLON) {
		p.nextToken()
		for !p.curTokenIs(candy_token.LBRACE) && !p.curTokenIs(candy_token.EOF) {
			p.nextToken()
		}
	}
	if !p.curTokenIs(candy_token.LBRACE) {
		if !p.expect(candy_token.LBRACE) {
			return nil
		}
	}
	fn.Body = p.parseBlockStatement()
	return fn
}

func (p *Parser) parseBreakStatement() candy_ast.Statement {
	s := &candy_ast.BreakStatement{Token: p.curToken}
	p.nextToken()
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseContinueStatement() candy_ast.Statement {
	s := &candy_ast.ContinueStatement{Token: p.curToken}
	p.nextToken()
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseDeleteStatement() candy_ast.Statement {
	s := &candy_ast.DeleteStatement{Token: p.curToken}
	p.nextToken() // past delete
	if p.curTokenIs(candy_token.LPAREN) {
		p.nextToken()
		s.Value = p.parseExpression(LOWEST)
		if !p.expectPeek(candy_token.RPAREN) {
			return nil
		}
		p.nextToken() // move past )
	} else {
		s.Value = p.parseExpression(LOWEST)
		p.nextToken() // move past last token of expr
	}
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseWithStatement() candy_ast.Statement {
	s := &candy_ast.WithStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(candy_token.ASSIGN) {
		return nil
	}
	p.nextToken()
	oldStop := p.exprStopAtBrace
	p.exprStopAtBrace = true
	s.Value = p.parseExpression(LOWEST)
	p.exprStopAtBrace = oldStop
	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	s.Body = p.parseBlockStatement()
	return s
}

func (p *Parser) parseLibraryStatement() candy_ast.Statement {
	s := &candy_ast.LibraryStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.STR) {
		return nil
	}
	s.Name = p.curToken.Literal
	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	body := &candy_ast.BlockStatement{Token: p.curToken}
	p.nextToken()
	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.IDENT) && p.curToken.Literal == "type" && p.peekTokenIs(candy_token.IDENT) {
			if st := p.parseTypeAsStructStatement(); st != nil {
				body.Statements = append(body.Statements, st)
			}
		} else if st := p.parseStatement(); st != nil {
			body.Statements = append(body.Statements, st)
		} else if !p.curTokenIs(candy_token.RBRACE) {
			p.nextToken()
		}
	}
	_ = p.expect(candy_token.RBRACE)
	s.Body = body
	return s
}

func (p *Parser) parseTypeAsStructStatement() candy_ast.Statement {
	// current token is IDENT("type"), next is struct name.
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	st := &candy_ast.StructStatement{Token: p.curToken, Name: name}
	p.nextToken()
	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.IDENT) && p.peekTokenIs(candy_token.COLON) {
			fname := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			p.nextToken() // :
			p.nextToken() // type
			ft := p.parseTypeIdentifier()
			st.Fields = append(st.Fields, candy_ast.Field{Name: fname, TypeName: ft})
			if p.curTokenIs(candy_token.COMMA) || p.curTokenIs(candy_token.SEMICOLON) {
				p.nextToken()
			}
			continue
		}
		p.nextToken()
	}
	_ = p.expect(candy_token.RBRACE)
	return st
}

func (p *Parser) parseForEachStatement() candy_ast.Statement {
	s := &candy_ast.ForEachStatement{Token: p.curToken}
	p.nextToken() // past foreach
	
	hasParen := false
	if p.curTokenIs(candy_token.LPAREN) {
		p.nextToken()
		hasParen = true
	}

	if !p.curTokenIs(candy_token.IDENT) {
		p.addErrorf("expected identifier in foreach, got %s", p.curToken.Type)
		return nil
	}
	s.Var = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()

	if !p.curTokenIs(candy_token.IN) {
		p.addErrorf("expected 'in' in foreach, got %s", p.curToken.Type)
		return nil
	}
	p.nextToken()

	p.exprStopAtBrace = true
	s.Iterable = p.parseExpression(LOWEST)
	p.exprStopAtBrace = false

	if hasParen {
		if !p.expectPeek(candy_token.RPAREN) {
			return nil
		}
		p.nextToken() // move past )
	} else {
		p.nextToken() // move past last token of iterable
	}

	if !p.curTokenIs(candy_token.LBRACE) {
		if !p.expectPeek(candy_token.LBRACE) {
			return nil
		}
		// p.nextToken() // BlockStatement parser starts with current token as {
	}
	s.Body = p.parseBlockStatement()
	return s
}

func (p *Parser) parseRepeatStatement() *candy_ast.RepeatStatement {
	s := &candy_ast.RepeatStatement{Token: p.curToken}
	p.nextToken() // past repeat
	p.exprStopAtBrace = true
	s.Count = p.parseExpression(LOWEST)
	p.exprStopAtBrace = false
	if !p.curTokenIs(candy_token.LBRACE) {
		if !p.expectPeek(candy_token.LBRACE) {
			return nil
		}
	}
	s.Body = p.parseBlockStatement()
	return s
}

func (p *Parser) parseLoopStatement() *candy_ast.LoopStatement {
	s := &candy_ast.LoopStatement{Token: p.curToken}
	p.nextToken() // past loop
	if !p.curTokenIs(candy_token.LBRACE) {
		if !p.expectPeek(candy_token.LBRACE) {
			return nil
		}
	}
	s.Body = p.parseBlockStatement()
	return s
}
