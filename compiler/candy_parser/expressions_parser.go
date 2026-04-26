package candy_parser

import (
	"candy/candy_ast"
	"candy/candy_lexer"
	"candy/candy_token"
	"strings"
)

func (p *Parser) parseExpression(precedence int) candy_ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.addErrorf("no prefix for %q", p.curToken.Type)
		return nil
	}
	left := prefix()
	if left == nil {
		return nil
	}
	for n := 0; n < 10000; n++ {
		if p.curTokenIs(candy_token.SEMICOLON) || p.curTokenIs(candy_token.RBRACE) || p.curTokenIs(candy_token.EOF) ||
			p.curTokenIs(candy_token.THEN) || p.curTokenIs(candy_token.END) || p.curTokenIs(candy_token.WEND) || p.curTokenIs(candy_token.NEXT) {
			break
		}
		// Support infix "not in" sugar by desugaring:
		//   a not in b  =>  not (a in b)
		if p.peekTokenIs(candy_token.NOT) {
			p.nextToken() // move to NOT
			if p.peekTokenIs(candy_token.IN) {
				notTok := p.curToken
				p.nextToken() // move to IN
				inExpr := p.parseInfixExpression(left)
				if inExpr == nil {
					return nil
				}
				left = &candy_ast.PrefixExpression{
					Token:    notTok,
					Operator: "not",
					Right:    inExpr,
				}
				continue
			}
			p.addErrorf("expected IN after NOT in infix expression")
			return nil
		}
		if p.peekTokenIs(candy_token.SEMICOLON) || p.peekTokenIs(candy_token.EOF) || p.peekTokenIs(candy_token.RBRACE) ||
			p.peekTokenIs(candy_token.RPAREN) || p.peekTokenIs(candy_token.RBRACK) ||
			p.peekTokenIs(candy_token.THEN) || p.peekTokenIs(candy_token.END) || p.peekTokenIs(candy_token.WEND) || p.peekTokenIs(candy_token.NEXT) ||
			p.peekTokenIs(candy_token.TO) || p.peekTokenIs(candy_token.STEP) {
			break
		}
		if !(precedence < p.peekPrecedence()) {
			break
		}
		if p.exprStopAtBrace && p.peekTokenIs(candy_token.LBRACE) {
			break
		}
		inf := p.infixParseFns[p.peekToken.Type]
		if inf == nil {
			break
		}
		p.nextToken()
		left = inf(left)
		if left == nil {
			return nil
		}
	}
	return left
}

func (p *Parser) parsePrefixExpression() candy_ast.Expression {
	expr := &candy_ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parsePostfixIncDec(left candy_ast.Expression) candy_ast.Expression {
	expr := &candy_ast.PostfixExpression{Token: p.curToken, Left: left, Operator: p.curToken.Literal}
	p.nextToken() // past ++ or --
	return expr
}

func (p *Parser) parseInfixExpression(left candy_ast.Expression) candy_ast.Expression {
	expr := &candy_ast.InfixExpression{Token: p.curToken, Left: left, Operator: p.curToken.Literal}
	pr := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(pr)
	return expr
}

func (p *Parser) parseRangeExpression(left candy_ast.Expression) candy_ast.Expression {
	expr := &candy_ast.RangeExpression{Token: p.curToken, Left: left}
	pr := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(pr)
	return expr
}

func (p *Parser) parseAssignExpression(left candy_ast.Expression) candy_ast.Expression {
	expr := &candy_ast.AssignExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}
	p.nextToken()
	expr.Value = p.parseExpression(LOWEST)
	return expr
}

func (p *Parser) parseGroupedExpression() candy_ast.Expression {
	tok := p.curToken
	p.nextToken() // past (

	exprs := []candy_ast.Expression{}
	exprs = append(exprs, p.parseExpression(TUPLE))

	for p.peekTokenIs(candy_token.COMMA) {
		p.nextToken() // to ,
		p.nextToken() // past ,
		exprs = append(exprs, p.parseExpression(TUPLE))
	}

	if !p.expectPeek(candy_token.RPAREN) {
		return nil
	}

	if len(exprs) == 1 {
		return &candy_ast.GroupedExpression{Token: tok, Expr: exprs[0]}
	}

	return &candy_ast.TupleLiteral{Token: tok, Elems: exprs}
}

func (p *Parser) parseIfExpression() candy_ast.Expression {
	expr := &candy_ast.IfExpression{Token: p.curToken}
	if p.peekTokenIs(candy_token.LPAREN) {
		p.nextToken()
		p.nextToken()
		expr.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(candy_token.RPAREN) {
			return nil
		}
	} else {
		p.nextToken()
		p.exprStopAtBrace = true
		expr.Condition = p.parseExpression(LOWEST)
		p.exprStopAtBrace = false
	}

	if p.curTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}
	if p.curTokenIs(candy_token.LBRACE) {
		expr.Consequence = p.parseBlockStatement()
	} else if p.peekTokenIs(candy_token.THEN) || p.curTokenIs(candy_token.THEN) {
		if p.curTokenIs(candy_token.THEN) {
			p.nextToken()
		} else {
			p.nextToken()
			p.nextToken()
		}
		expr.Consequence = &candy_ast.ExpressionStatement{Token: p.curToken, Expression: p.parseExpression(LOWEST)}
		if p.peekTokenIs(candy_token.ELSE) {
			p.nextToken()
			p.nextToken()
			expr.Alternative = &candy_ast.ExpressionStatement{Token: p.curToken, Expression: p.parseExpression(LOWEST)}
		}
	} else if p.expectPeek(candy_token.LBRACE) {
		expr.Consequence = p.parseBlockStatement()
	} else {
		return nil
	}
	if p.curTokenIs(candy_token.ELSEIF) {
		expr.Alternative = p.parseIfExpression().(candy_ast.Statement)
	} else if p.curTokenIs(candy_token.ELSE) {
		p.nextToken() // past ELSE
		if p.curTokenIs(candy_token.IF) {
			expr.Alternative = p.parseIfExpression().(candy_ast.Statement)
		} else if p.curTokenIs(candy_token.LBRACE) {
			expr.Alternative = p.parseBlockStatement()
		} else if p.expectPeek(candy_token.LBRACE) {
			expr.Alternative = p.parseBlockStatement()
		}
	}
	return expr
}

func (p *Parser) parseStructLiteralExpr(name candy_ast.Expression) candy_ast.Expression {
	sl := &candy_ast.StructLiteral{Token: p.curToken, Name: name, Fields: make(map[string]candy_ast.Expression)}
	p.nextToken() // skip {

	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		if !p.curTokenIs(candy_token.IDENT) && !p.isKeyword(p.curToken.Type) {
			p.addErrorf("expected field name, got %s", p.curToken.Type)
			return nil
		}
		fname := p.curToken.Literal
		p.nextToken()
		if !p.expect(candy_token.COLON) {
			return nil
		}
		// Use TUPLE precedence so `,` separates struct fields instead of being
		// absorbed into a tuple expression for the field value.
		val := p.parseExpression(TUPLE)
		sl.Fields[fname] = val
		p.nextToken() // move past last token of expression

		// consume separator
		if p.curTokenIs(candy_token.COMMA) || p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
		}
	}

	if !p.curTokenIs(candy_token.RBRACE) {
		p.addErrorf("expected }, got %s", p.curToken.Type)
		return nil
	}
	return sl
}

func (p *Parser) parseCallExpression(fn candy_ast.Expression) candy_ast.Expression {
	tok := p.curToken
	var typeArgs []candy_ast.Expression
	if p.curTokenIs(candy_token.LT) {
		p.nextToken() // <
		for !p.curTokenIs(candy_token.GT) && !p.curTokenIs(candy_token.EOF) {
			typeArgs = append(typeArgs, p.parseExpression(LOWEST))
			if p.curTokenIs(candy_token.COMMA) {
				p.nextToken()
			}
		}
		if !p.expectPeek(candy_token.GT) {
			return nil
		}
		if !p.peekTokenIs(candy_token.LPAREN) {
			p.addErrorf("expected ( after generic type params")
			return nil
		}
		p.nextToken() // move to (
	}

	return &candy_ast.CallExpression{Token: tok, Function: fn, TypeArguments: typeArgs, Arguments: p.parseCallArgs()}
}

func (p *Parser) parseDotExpression(left candy_ast.Expression) candy_ast.Expression {
	tok := p.curToken
	isSafe := tok.Type == candy_token.SAFE_DOT
	p.nextToken() // skip . or ?.
	id := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	return &candy_ast.DotExpression{Token: tok, Left: left, Right: id, IsSafe: isSafe}
}

func (p *Parser) parseSafeAccessExpression(left candy_ast.Expression) candy_ast.Expression {
	tok := p.curToken // ?.
	if p.peekTokenIs(candy_token.LPAREN) {
		p.nextToken() // (
		return &candy_ast.CallExpression{
			Token:     tok,
			Function:  left,
			Arguments: p.parseCallArgs(),
			IsSafe:    true,
		}
	}
	if p.peekTokenIs(candy_token.LBRACK) {
		p.nextToken() // [
		p.nextToken() // first token inside []
		idx := p.parseExpression(LOWEST)
		if !p.expectPeek(candy_token.RBRACK) {
			return nil
		}
		return &candy_ast.IndexExpression{
			Token:  tok,
			Base:   left,
			Index:  idx,
			IsSafe: true,
		}
	}
	return p.parseDotExpression(left)
}

func (p *Parser) parseCallArgs() []candy_ast.Expression {
	args := []candy_ast.Expression{}
	if p.peekTokenIs(candy_token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseSingleCallArg())
	for p.peekTokenIs(candy_token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseSingleCallArg())
	}
	// Allow semicolon insertion before closing ')' in multiline argument lists.
	for p.peekTokenIs(candy_token.SEMICOLON) {
		p.nextToken()
	}
	if p.peekTokenIs(candy_token.RPAREN) {
		p.nextToken()
		return args
	}
	// Recovery for legacy `Type { ... }` call-arg style where semicolon insertion
	// can consume/skip the closing ')' token in some newline forms.
	if p.peekTokenIs(candy_token.EOF) || p.peekTokenIs(candy_token.RBRACE) || p.isStatementStart(p.peekToken.Type) {
		return args
	}
	if !p.expectPeek(candy_token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseSingleCallArg() candy_ast.Expression {
	if p.curTokenIs(candy_token.IDENT) && p.peekTokenIs(candy_token.COLON) {
		namedTok := p.curToken
		name := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken() // :
		p.nextToken() // value
		// Parse with TUPLE precedence so lambda arrows inside args are accepted
		// while commas still terminate the argument.
		val := p.parseExpression(TUPLE)
		return &candy_ast.NamedArgumentExpression{
			Token: namedTok,
			Name:  name,
			Value: val,
		}
	}
	// Parse with TUPLE precedence so `x => ...` is parsed as a lambda argument.
	return p.parseExpression(TUPLE)
}

func (p *Parser) parseIndexExpression(left candy_ast.Expression) candy_ast.Expression {
	tok := p.curToken
	p.nextToken()
	idx := p.parseExpression(LOWEST)
	if !p.expectPeek(candy_token.RBRACK) {
		return left
	}
	return &candy_ast.IndexExpression{Token: tok, Base: left, Index: idx}
}

func (p *Parser) parseLambdaInfix(left candy_ast.Expression) candy_ast.Expression {
	l := &candy_ast.LambdaExpression{Token: p.curToken}
	l.Parameters = p.convertToParams(left)
	p.nextToken()
	l.Body = p.parseExpression(LOWEST)
	return l
}

func (p *Parser) convertToParams(expr candy_ast.Expression) []candy_ast.Parameter {
	switch e := expr.(type) {
	case *candy_ast.Identifier:
		return []candy_ast.Parameter{{Token: e.Token, Name: e, TypeName: e}}
	case *candy_ast.GroupedExpression:
		return p.convertToParams(e.Expr)
	case *candy_ast.TupleLiteral:
		var res []candy_ast.Parameter
		for _, el := range e.Elems {
			res = append(res, p.convertToParams(el)...)
		}
		return res
	}
	return nil
}

func (p *Parser) parseStringLiteral() candy_ast.Expression {
	lit := p.curToken.Literal
	// Support both ${} and {} for interpolation
	if !strings.Contains(lit, "{") {
		return &candy_ast.StringLiteral{Token: p.curToken, Value: lit}
	}

	tok := p.curToken
	parts := []candy_ast.Expression{}
	current := lit

	for {
		// First `{` that starts an interpolation, not a lexer-preserved `\{` literal.
		start := indexInterpolationBrace(current)
		if start == -1 {
			if len(current) > 0 {
				parts = append(parts, &candy_ast.StringLiteral{Token: tok, Value: current})
			}
			break
		}

		// Optional: handle ${} by checking if previous char is $
		realStart := start
		if start > 0 && current[start-1] == '$' {
			realStart = start - 1
		}

		if realStart > 0 {
			parts = append(parts, &candy_ast.StringLiteral{Token: tok, Value: current[:realStart]})
		}

		depth := 1
		end := -1
		for i := start + 1; i < len(current); i++ {
			if current[i] == '{' {
				depth++
			} else if current[i] == '}' {
				depth--
				if depth == 0 {
					end = i
					break
				}
			}
		}

		if end != -1 {
			exprStr := current[start+1 : end]
			if len(strings.TrimSpace(exprStr)) == 0 {
				// Handle empty {} as literal for use with format()
				parts = append(parts, &candy_ast.StringLiteral{Token: tok, Value: "{}"})
			} else {
				// Parse the expression inside braces
				innerParser := New(candy_lexer.New(exprStr))
				expr := innerParser.parseExpression(LOWEST)
				if expr != nil {
					parts = append(parts, expr)
				}
			}
			current = current[end+1:]
		} else {
			// Unclosed brace, treat as literal
			parts = append(parts, &candy_ast.StringLiteral{Token: tok, Value: current[:start+1]})
			current = current[start+1:]
		}
	}

	if len(parts) == 1 {
		if s, ok := parts[0].(*candy_ast.StringLiteral); ok {
			return s
		}
	}

	return &candy_ast.InterpolatedStringLiteral{Token: tok, Parts: parts}
}

// indexInterpolationBrace is the first `{` that starts an interpolation; `{` is skipped
// when it follows an **odd** number of backslashes (from `\{` in a string literal).
func indexInterpolationBrace(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] != '{' {
			continue
		}
		bs := 0
		for j := i - 1; j >= 0 && s[j] == '\\'; j-- {
			bs++
		}
		if bs%2 == 1 {
			continue
		}
		return i
	}
	return -1
}

func (p *Parser) parseArrayLiteral() candy_ast.Expression {
	tok := p.curToken
	p.nextToken() // past `[`
	var el []candy_ast.Expression
	if p.curTokenIs(candy_token.RBRACK) {
		p.nextToken()
		return &candy_ast.ArrayLiteral{Token: tok, Elem: el}
	}
	for {
		el = append(el, p.parseExpression(TUPLE))
		if p.peekTokenIs(candy_token.RBRACK) {
			p.nextToken() // move to ]
			return &candy_ast.ArrayLiteral{Token: tok, Elem: el}
		}
		if p.peekTokenIs(candy_token.COMMA) {
			p.nextToken() // to ,
			p.nextToken() // past ,
			continue
		}
		if !p.expectPeek(candy_token.RBRACK) {
			return nil
		}
		return &candy_ast.ArrayLiteral{Token: tok, Elem: el}
	}
}

func (p *Parser) parseTupleInfix(left candy_ast.Expression) candy_ast.Expression {
	tok := p.curToken
	var elems []candy_ast.Expression
	if tl, ok := left.(*candy_ast.TupleLiteral); ok {
		elems = append(elems, tl.Elems...)
	} else {
		elems = append(elems, left)
	}

	p.nextToken() // past ,
	right := p.parseExpression(TUPLE)
	if tl, ok := right.(*candy_ast.TupleLiteral); ok {
		elems = append(elems, tl.Elems...)
	} else {
		elems = append(elems, right)
	}

	return &candy_ast.TupleLiteral{Token: tok, Elems: elems}
}

func (p *Parser) parseMapLiteral() candy_ast.Expression {
	tok := p.curToken
	endTok := candy_token.TokenType(candy_token.RBRACE)
	if p.peekTokenIs(candy_token.LBRACE) {
		p.nextToken()
	} else if p.peekTokenIs(candy_token.LPAREN) {
		p.nextToken()
		endTok = candy_token.RPAREN
	} else {
		p.addErrorf("expected { or ( after map, got %s", p.peekToken.Type)
		return nil
	}
	p.nextToken() // first token inside map literal
	var pairs []candy_ast.MapPair
	for !p.curTokenIs(endTok) && !p.curTokenIs(candy_token.EOF) {
		// Allow optional separators between entries.
		if p.curTokenIs(candy_token.SEMICOLON) || p.curTokenIs(candy_token.COMMA) {
			p.nextToken()
			continue
		}
		ke := p.parseExpression(LOWEST)
		if p.peekTokenIs(endTok) {
			p.nextToken()
			break
		}
		if !p.expectPeek(candy_token.COLON) {
			break
		}
		p.nextToken()
		ve := p.parseExpression(LOWEST)
		pairs = append(pairs, candy_ast.MapPair{Key: ke, Value: ve})
		if p.peekTokenIs(candy_token.COMMA) || p.peekTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
			continue
		}
		if p.peekTokenIs(endTok) {
			p.nextToken()
			break
		}
	}
	if !p.curTokenIs(endTok) {
		p.addErrorf("expected %s, got %s", endTok, p.curToken.Type)
		return nil
	}
	return &candy_ast.MapLiteral{Token: tok, Pairs: pairs}
}

// parseBraceMapLiteral parses bare object literals like `{x: 1, y: 2}`.
// This is used in expression position (e.g. `point = {x: 10, y: 20}`).
func (p *Parser) parseBraceMapLiteral() candy_ast.Expression {
	tok := p.curToken
	p.nextToken() // past '{'
	var pairs []candy_ast.MapPair
	for !p.curTokenIs(candy_token.RBRACE) && !p.curTokenIs(candy_token.EOF) {
		// Optional separators in multiline object literals.
		if p.curTokenIs(candy_token.SEMICOLON) || p.curTokenIs(candy_token.COMMA) {
			p.nextToken()
			continue
		}
		var ke candy_ast.Expression
		var shorthandIdent *candy_ast.Identifier
		if p.curTokenIs(candy_token.IDENT) || p.isKeyword(p.curToken.Type) {
			ke = &candy_ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
			shorthandIdent = &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		} else {
			ke = p.parseExpression(LOWEST)
		}
		if p.peekTokenIs(candy_token.RBRACE) {
			if shorthandIdent != nil {
				pairs = append(pairs, candy_ast.MapPair{Key: ke, Value: shorthandIdent})
			}
			p.nextToken()
			break
		}
		if shorthandIdent != nil && p.peekTokenIs(candy_token.COMMA) {
			pairs = append(pairs, candy_ast.MapPair{Key: ke, Value: shorthandIdent})
			p.nextToken()
			p.nextToken()
			continue
		}
		if !p.expectPeek(candy_token.COLON) {
			break
		}
		p.nextToken()
		ve := p.parseExpression(LOWEST)
		pairs = append(pairs, candy_ast.MapPair{Key: ke, Value: ve})
		if p.peekTokenIs(candy_token.COMMA) || p.peekTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
			p.nextToken()
			continue
		}
		if p.peekTokenIs(candy_token.RBRACE) {
			p.nextToken()
			break
		}
	}
	if !p.curTokenIs(candy_token.RBRACE) {
		p.addErrorf("expected }, got %s", p.curToken.Type)
		return nil
	}
	return &candy_ast.MapLiteral{Token: tok, Pairs: pairs}
}

func (p *Parser) parseNewExpression() candy_ast.Expression {
	tok := p.curToken
	p.nextToken()
	// Usually followed by a call: new Player(...)
	// or just a type: new Player
	expr := p.parseExpression(PREFIX)
	return &candy_ast.CallExpression{
		Token:     tok,
		Function:  &candy_ast.Identifier{Token: tok, Value: "new"},
		Arguments: []candy_ast.Expression{expr},
	}
}

func (p *Parser) parseTernaryExpression(condition candy_ast.Expression) candy_ast.Expression {
	expr := &candy_ast.TernaryExpression{
		Token:     p.curToken,
		Condition: condition,
	}

	p.nextToken() // past ?
	expr.Consequence = p.parseExpression(TERNARY)

	if !p.expectPeek(candy_token.COLON) {
		return nil
	}

	p.nextToken() // past :
	expr.Alternative = p.parseExpression(TERNARY)

	return expr
}

func (p *Parser) parsePipelineExpression(left candy_ast.Expression) candy_ast.Expression {
	tok := p.curToken
	p.nextToken() // past |>
	right := p.parseExpression(PIPELINE)

	if right == nil {
		return left
	}

	// Desugar: left |> right(...) => right(left, ...)
	// If right is already a CallExpression, prepend left to its arguments.
	if call, ok := right.(*candy_ast.CallExpression); ok {
		call.Arguments = append([]candy_ast.Expression{left}, call.Arguments...)
		return call
	}

	// Otherwise, wrap it in a call: left |> right => right(left)
	return &candy_ast.CallExpression{
		Token:     tok,
		Function:  right,
		Arguments: []candy_ast.Expression{left},
	}
}
