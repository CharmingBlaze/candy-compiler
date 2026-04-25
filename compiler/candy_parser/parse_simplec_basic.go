package candy_parser

import (
	"candy/candy_ast"
	"candy/candy_token"
)

func (p *Parser) parseImportStatement() *candy_ast.ImportStatement {
	s := &candy_ast.ImportStatement{Token: p.curToken}
	p.nextToken() // import
	if p.curTokenIs(candy_token.STR) {
		s.Path = p.curToken.Literal
		p.nextToken()
	} else if p.curTokenIs(candy_token.IDENT) {
		s.Path = p.curToken.Literal
		p.nextToken()
		for p.curTokenIs(candy_token.DOT) {
			p.nextToken()
			if !p.expect(candy_token.IDENT) {
				break
			}
			s.Path += "." + p.curToken.Literal
			p.nextToken()
		}
	} else {
		p.addErrorf("expected import path string or identifier, got %s", p.curToken.Type)
		return nil
	}
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseModuleStatement() *candy_ast.ModuleStatement {
	s := &candy_ast.ModuleStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	s.Body = p.parseBlockStatement()
	return s
}

func (p *Parser) parseEnumStatement() *candy_ast.EnumStatement {
	s := &candy_ast.EnumStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	p.nextToken() // past {
	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
			continue
		}
		if !p.curTokenIs(candy_token.IDENT) {
			p.addErrorf("expected variant name, got %s", p.curToken.Type)
			return nil
		}
		variant := &candy_ast.EnumVariant{
			Name: &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekTokenIs(candy_token.ASSIGN) {
			p.nextToken()
			p.nextToken()
			// Keep enum assignments scalar: stop before `,` so
			// `Playing = 10, Paused` doesn't become tuple `(10, paused)`.
			variant.Value = p.parseExpression(TUPLE)
		}
		s.Variants = append(s.Variants, variant)
		if p.peekTokenIs(candy_token.COMMA) {
			p.nextToken()
			p.nextToken()
		} else {
			// `name` or `= value` is fully consumed: step to the next `}` or variant
			if !p.curTokenIs(candy_token.RBRACE) {
				p.nextToken()
			}
		}
	}
	_ = p.expect(candy_token.RBRACE)
	return s
}

func (p *Parser) parseTryStatement() *candy_ast.TryStatement {
	ts := &candy_ast.TryStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	ts.TryBody = p.parseBlockStatement()

	for p.curTokenIs(candy_token.CATCH) {
		cc := &candy_ast.CatchClause{Token: p.curToken}
		p.nextToken() // past 'catch'
		cc.Type = p.parseTypeIdentifier()
		if p.peekTokenIs(candy_token.IDENT) {
			p.nextToken()
			cc.Identifier = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}
		if !p.expectPeek(candy_token.LBRACE) {
			return nil
		}
		cc.Body = p.parseBlockStatement()
		ts.CatchClauses = append(ts.CatchClauses, cc)
	}

	if p.curTokenIs(candy_token.FINALLY) {
		if !p.expectPeek(candy_token.LBRACE) {
			return nil
		}
		ts.FinallyBody = p.parseBlockStatement()
	}

	return ts
}

func (p *Parser) parseRunStatement() *candy_ast.RunStatement {
	s := &candy_ast.RunStatement{Token: p.curToken}
	p.nextToken()
	s.Value = p.parseExpression(LOWEST)
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseAwaitExpression() candy_ast.Expression {
	ae := &candy_ast.AwaitExpression{Token: p.curToken}
	p.nextToken()
	ae.Value = p.parseExpression(PREFIX)
	return ae
}

func (p *Parser) parseTypeofExpression() candy_ast.Expression {
	tok := p.curToken
	if !p.expectPeek(candy_token.LPAREN) {
		return nil
	}
	p.nextToken() // past (
	// In SimpleC-BASIC, typeof(Type) returns type info.
	// For now, let's just parse an identifier or type.
	val := p.parseExpression(LOWEST)
	if !p.expect(candy_token.RPAREN) {
		return nil
	}
	// We'll wrap it in a CallExpression or dedicated node if needed.
	// For simplicity, let's use a CallExpression with "typeof" name for now,
	// or create a dedicated node if we had one. (We didn't create typeof_expr.go yet).
	return &candy_ast.CallExpression{
		Token:     tok,
		Function:  &candy_ast.Identifier{Token: tok, Value: "typeof"},
		Arguments: []candy_ast.Expression{val},
	}
}

func (p *Parser) parseAttributes() []*candy_ast.Attribute {
	var attrs []*candy_ast.Attribute
	for p.curTokenIs(candy_token.LBRACK) {
		attr := &candy_ast.Attribute{Token: p.curToken}
		if !p.expectPeek(candy_token.IDENT) {
			return nil
		}
		attr.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		attr.Arguments = make(map[string]candy_ast.Expression)
		if p.peekTokenIs(candy_token.LPAREN) {
			p.nextToken()
			p.nextToken()
			for !p.curTokenIs(candy_token.RPAREN) && !p.curTokenIs(candy_token.EOF) {
				if !p.curTokenIs(candy_token.IDENT) {
					break
				}
				key := p.curToken.Literal
				if !p.expectPeek(candy_token.COLON) {
					break
				}
				p.nextToken()
				val := p.parseExpression(LOWEST)
				attr.Arguments[key] = val
				if p.peekTokenIs(candy_token.COMMA) {
					p.nextToken()
					p.nextToken()
				} else {
					// rvalue: parseExpression can leave cur on the STR/INT/… token while peek is `)`.
					if p.peekTokenIs(candy_token.RPAREN) && !p.curTokenIs(candy_token.RPAREN) {
						p.nextToken()
					}
				}
			}
			_ = p.expect(candy_token.RPAREN)
		}
		// `);` then a newline can insert `;` before `]`.
		for p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
		}
		if p.curTokenIs(candy_token.RBRACK) {
			p.nextToken()
		} else {
			if !p.expectPeek(candy_token.RBRACK) {
				return nil
			}
			p.nextToken()
		}
		attrs = append(attrs, attr)
	}
	return attrs
}
