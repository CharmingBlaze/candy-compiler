package candy_parser

import (
	"candy/candy_ast"
	"candy/candy_token"
)

func (p *Parser) parsePackageStatement() *candy_ast.PackageStatement {
	s := &candy_ast.PackageStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = p.curToken.Literal
	if !p.expectPeek(candy_token.SEMICOLON) {
		return nil
	}
	return s
}

func (p *Parser) parseClassStatement(sealed bool) *candy_ast.ClassStatement {
	s := &candy_ast.ClassStatement{Token: p.curToken, Sealed: sealed}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	s.TypeParameters = p.parseTypeParamsIfAny()

	if p.peekTokenIs(candy_token.LPAREN) {
		p.nextToken() // to (
		p.nextToken() // past (
		s.Parameters = p.parseFunctionParameters()
	}

	if p.peekTokenIs(candy_token.EXTENDS) {
		p.nextToken()
		if !p.expectPeek(candy_token.IDENT) {
			return nil
		}
		s.Base = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if p.curTokenIs(candy_token.LBRACE) || p.peekTokenIs(candy_token.LBRACE) {
		if !p.curTokenIs(candy_token.LBRACE) {
			p.nextToken() // to {
		}
		s.Members = p.parseClassMembers()
	}

	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseClassMembers() []candy_ast.Statement {
	var members []candy_ast.Statement
	p.nextToken() // past {

	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
			continue
		}

		m := p.parseStatement()
		if m != nil {
			members = append(members, m)
		} else {
			p.nextToken()
		}
	}

	if !p.expect(candy_token.RBRACE) {
		return nil
	}
	return members
}

func (p *Parser) parseObjectStatement() *candy_ast.ObjectStatement {
	s := &candy_ast.ObjectStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(candy_token.LBRACE) {
		p.nextToken()
		s.Members = p.parseClassMembers()
	}

	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseInterfaceStatement() *candy_ast.InterfaceStatement {
	s := &candy_ast.InterfaceStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	s.TypeParameters = p.parseTypeParamsIfAny()

	if !p.expectPeek(candy_token.LBRACE) {
		return nil
	}
	p.nextToken()

	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
			continue
		}

		if p.curTokenIs(candy_token.IDENT) {
			method := &candy_ast.InterfaceMethod{}
			method.ReturnType = p.parseTypeIdentifier()
			p.nextToken()
			if !p.curTokenIs(candy_token.IDENT) {
				p.addErrorf("expected method name, got %s", p.curToken.Type)
				break
			}
			method.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			if !p.expectPeek(candy_token.LPAREN) {
				break
			}
			p.nextToken()
			method.Parameters = p.parseFunctionParameters()
			s.Methods = append(s.Methods, method)
		} else {
			p.nextToken()
		}

		if p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
		}
	}

	_ = p.expect(candy_token.RBRACE)
	_ = p.expectSemicolon()
	return s
}

func (p *Parser) parseTraitStatement() *candy_ast.TraitStatement {
	s := &candy_ast.TraitStatement{Token: p.curToken}
	if !p.expectPeek(candy_token.IDENT) {
		return nil
	}
	s.Name = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekTokenIs(candy_token.LBRACE) {
		p.nextToken()
		p.nextToken()
		for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
			p.nextToken()
		}
	}
	if !p.curTokenIs(candy_token.RBRACE) {
		if !p.expectPeek(candy_token.RBRACE) {
			return nil
		}
	}
	if !p.expectPeek(candy_token.SEMICOLON) {
		return nil
	}
	return s
}

func (p *Parser) parseExternFunctionStatement() *candy_ast.ExternFunctionStatement {
	s := &candy_ast.ExternFunctionStatement{Token: p.curToken}
	// Support both:
	//   extern fun name(params): Ret { }
	//   extern name(params): Ret
	if p.peekTokenIs(candy_token.FUNCTION) {
		p.nextToken() // to fun/function token
	}
	p.nextToken() // function name
	if !p.curTokenIs(candy_token.IDENT) && !p.curTokenIs(candy_token.MAYBE) {
		p.addErrorf("expected extern function name, got %s", p.curToken.Type)
		return nil
	}
	fn := &candy_ast.FunctionStatement{
		Token: p.curToken,
		Name:  &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	if !p.expectPeek(candy_token.LPAREN) {
		return nil
	}
	p.nextToken() // first param token or ')'
	fn.Parameters = p.parseExternParameters(fn)
	if p.curTokenIs(candy_token.COLON) || p.curTokenIs(candy_token.AS) {
		p.nextToken()
		fn.ReturnType = p.parseTypeIdentifier()
		if fn.ReturnType != nil {
			p.nextToken()
		}
	}
	if p.curTokenIs(candy_token.LBRACE) {
		fn.Body = p.parseBlockStatement()
	} else if p.peekTokenIs(candy_token.LBRACE) {
		p.nextToken()
		fn.Body = p.parseBlockStatement()
	}
	_ = p.expectSemicolon()
	s.Function = fn
	return s
}

func (p *Parser) parseExternParameters(fn *candy_ast.FunctionStatement) []candy_ast.Parameter {
	var params []candy_ast.Parameter
	if p.curTokenIs(candy_token.RPAREN) {
		p.nextToken()
		return params
	}
	for !p.curTokenIs(candy_token.EOF) {
		// Variadic marker: `...` (tokenized as `..` + `.`).
		if p.curTokenIs(candy_token.RANGE) && p.peekTokenIs(candy_token.DOT) {
			fn.Variadic = true
			p.nextToken() // dot
			if p.peekTokenIs(candy_token.IDENT) {
				p.nextToken()
				name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
				params = append(params, candy_ast.Parameter{
					Token:    name.Token,
					Name:     name,
					TypeName: &candy_ast.Identifier{Token: name.Token, Value: "variadic"},
				})
			}
			if p.curTokenIs(candy_token.RPAREN) {
				p.nextToken()
			} else if p.peekTokenIs(candy_token.RPAREN) {
				p.nextToken() // to )
				p.nextToken() // past )
			}
			return params
		}

		param := p.parseOneParam()
		params = append(params, param)
		if p.curTokenIs(candy_token.RPAREN) {
			p.nextToken()
			return params
		}
		if p.curTokenIs(candy_token.COMMA) {
			p.nextToken()
			continue
		}
		p.addErrorf("expected , or ) in extern param list, got %s", p.curToken.Type)
		return params
	}
	return params
}

func (p *Parser) parseSuspendFunctionStatement() *candy_ast.FunctionStatement {
	if !p.expectPeek(candy_token.FUNCTION) {
		return nil
	}
	fn := p.parseFunctionStatement()
	if fn != nil {
		fn.Suspend = true
	}
	return fn
}

func (p *Parser) parseTypeParamsIfAny() []*candy_ast.Identifier {
	var out []*candy_ast.Identifier
	if !p.peekTokenIs(candy_token.LT) {
		return out
	}
	p.nextToken() // <
	p.nextToken() // first ident
	for !p.curTokenIs(candy_token.GT) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.IDENT) {
			out = append(out, &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		}
		if p.peekTokenIs(candy_token.COMMA) {
			p.nextToken()
			p.nextToken()
			continue
		}
		if p.peekTokenIs(candy_token.GT) {
			p.nextToken()
			break
		}
		break
	}
	return out
}
