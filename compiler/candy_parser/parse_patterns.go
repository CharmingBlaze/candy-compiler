package candy_parser

import (
	"candy/candy_ast"
	"candy/candy_token"
)

// parseWhenExpression parses:
//   - when { cond: expr; ... else: expr; }
//   - when (subject) { arm => expr; else => expr; }  // sugar via subject matching
func (p *Parser) parseWhenExpression() candy_ast.Expression {
	tok := p.curToken
	w := &candy_ast.WhenExpression{Token: tok}
	if p.peekTokenIs(candy_token.LPAREN) {
		p.nextToken() // to (
		p.nextToken() // start subject
		w.Subject = p.parseExpression(LOWEST)
		if !p.expectPeek(candy_token.RPAREN) {
			return w
		}
	}
	if !p.expectPeek(candy_token.LBRACE) {
		return w
	}
	p.nextToken()
	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.ELSE) {
			if !p.expectPeek(candy_token.COLON) {
				return w
			}
			p.nextToken()
			w.ElseV = p.parseExpression(LOWEST)
			if p.peekTokenIs(candy_token.SEMICOLON) {
				p.nextToken()
			}
			p.nextToken()
			continue
		}
		cond := p.parseExpression(LOWEST)
		if !p.expectPeek(candy_token.COLON) {
			return w
		}
		p.nextToken()
		body := p.parseExpression(LOWEST)
		w.Arms = append(w.Arms, candy_ast.WhenArm{Cond: cond, Body: body})
		if p.peekTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
		}
		p.nextToken()
	}
	return w
}

// parseMatchExpression: `match ( sub ) { pat: body; }`
func (p *Parser) parseMatchExpression() candy_ast.Expression {
	m := &candy_ast.MatchExpression{Token: p.curToken}
	if !p.expectPeek(candy_token.LPAREN) {
		return m
	}
	p.nextToken()
	m.Subject = p.parseExpression(LOWEST)
	if !p.expectPeek(candy_token.RPAREN) {
		return m
	}
	if !p.expectPeek(candy_token.LBRACE) {
		return m
	}
	p.nextToken()
	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if p.curTokenIs(candy_token.ELSE) {
			p.nextToken()
			if p.peekTokenIs(candy_token.COLON) {
				p.nextToken()
			}
			m.Default = p.parseExpression(LOWEST)
			if p.peekTokenIs(candy_token.SEMICOLON) {
				p.nextToken()
			}
			continue
		}
		pat := p.parseExpression(LOWEST)

		var guard candy_ast.Expression
		if p.curTokenIs(candy_token.IF) {
			p.nextToken()
			guard = p.parseExpression(LOWEST)
		} else if p.peekTokenIs(candy_token.IF) {
			p.nextToken() // IF
			p.nextToken() // start of expr
			guard = p.parseExpression(LOWEST)
		}

		if p.peekTokenIs(candy_token.ARROW) || p.peekTokenIs(candy_token.COLON) {
			p.nextToken()
			p.nextToken()
		}

		bod := p.parseExpression(LOWEST)
		m.Branches = append(m.Branches, candy_ast.MatchBranch{Pat: pat, Guard: guard, Body: bod})

		if p.peekTokenIs(candy_token.SEMICOLON) || p.peekTokenIs(candy_token.COMMA) {
			p.nextToken()
			p.nextToken()
		} else if p.curTokenIs(candy_token.SEMICOLON) || p.curTokenIs(candy_token.COMMA) {
			p.nextToken()
		}
	}
	_ = p.expect(candy_token.RBRACE)
	return m
}

func (p *Parser) parseIsExpression(left candy_ast.Expression) candy_ast.Expression {
	expr := &candy_ast.IsExpression{Token: p.curToken, Left: left}
	p.nextToken()
	tn := p.parseTypeIdentifier()
	if tn == nil {
		return nil
	}
	expr.TypeName = tn
	return expr
}
