package candy_parser

import (
	"candy/candy_ast"
	"candy/candy_lexer"
	"candy/candy_report"
	"candy/candy_token"
	"fmt"
	"reflect"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	TUPLE
	ARROW
	OR
	AND
	TERNARY
	RANGE
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	POSTFIX
	NULL_COALESCE
)

var precedences = map[candy_token.TokenType]int{
	candy_token.ARROW:         ARROW,
	candy_token.ASSIGN:        ASSIGN,
	candy_token.PLUS_ASSIGN:   ASSIGN,
	candy_token.MINUS_ASSIGN:  ASSIGN,
	candy_token.STAR_ASSIGN:   ASSIGN,
	candy_token.SLASH_ASSIGN:  ASSIGN,
	candy_token.RANGE:         RANGE,
	candy_token.RANGE_EXCL:    RANGE,
	candy_token.EQ:            EQUALS,
	candy_token.NOT_EQ:        EQUALS,
	candy_token.IS:            EQUALS,
	candy_token.LT:            LESSGREATER,
	candy_token.GT:            LESSGREATER,
	candy_token.LTE:           LESSGREATER,
	candy_token.GTE:           LESSGREATER,
	candy_token.IN:            LESSGREATER,
	candy_token.PLUS:          SUM,
	candy_token.MINUS:         SUM,
	candy_token.BIT_OR:        SUM,
	candy_token.BIT_XOR:       SUM,
	candy_token.AMPERSAND:     PRODUCT,
	candy_token.SHL:           PRODUCT,
	candy_token.SHR:           PRODUCT,
	candy_token.SLASH:         PRODUCT,
	candy_token.ASTERISK:      PRODUCT,
	candy_token.PERCENT:       PRODUCT,
	candy_token.MOD:           PRODUCT,
	candy_token.LPAREN:        CALL,
	candy_token.LBRACK:        CALL,
	candy_token.DOT:           CALL,
	candy_token.SAFE_DOT:      CALL,
	candy_token.LBRACE:        CALL,
	candy_token.INC:           POSTFIX,
	candy_token.DEC:           POSTFIX,
	candy_token.NULL_COALESCE: NULL_COALESCE,
	candy_token.QUESTION:      TERNARY,
	candy_token.AND:           AND,
	candy_token.OR:            OR,
	candy_token.AND_AND:       AND,
	candy_token.OR_OR:         OR,
	candy_token.COMMA:         TUPLE,
}

type (
	prefixParseFn func() candy_ast.Expression
	infixParseFn  func(candy_ast.Expression) candy_ast.Expression
)

// Parser is a hand-written Pratt (TDOP) parser.
type Parser struct {
	l         *candy_lexer.Lexer
	errors    []candy_report.Diagnostic
	errorSeen map[string]struct{}
	recovery  []string
	curToken  candy_token.Token
	peekToken candy_token.Token

	prefixParseFns map[candy_token.TokenType]prefixParseFn
	infixParseFns  map[candy_token.TokenType]infixParseFn

	exprStopAtBrace bool
}

func New(l *candy_lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []candy_report.Diagnostic{},
		errorSeen:      make(map[string]struct{}),
		prefixParseFns: make(map[candy_token.TokenType]prefixParseFn),
		infixParseFns:  make(map[candy_token.TokenType]infixParseFn),
	}
	p.nextToken()
	p.nextToken()

	p.registerPrefix(candy_token.IDENT, p.parseIdentifier)
	p.registerPrefix(candy_token.INT, p.parseIntegerLiteral)
	p.registerPrefix(candy_token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(candy_token.STR, p.parseStringLiteral)
	p.registerPrefix(candy_token.TRUE, p.parseBoolean)
	p.registerPrefix(candy_token.FALSE, p.parseBoolean)
	p.registerPrefix(candy_token.NULL, p.parseNull)
	p.registerPrefix(candy_token.BANG, p.parsePrefixExpression)
	p.registerPrefix(candy_token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(candy_token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(candy_token.LBRACK, p.parseArrayLiteral)
	p.registerPrefix(candy_token.LBRACE, p.parseBraceMapLiteral)
	p.registerPrefix(candy_token.MAP, p.parseMapLiteral)
	p.registerPrefix(candy_token.WHEN, p.parseWhenExpression)
	p.registerPrefix(candy_token.MATCH, p.parseMatchExpression)
	p.registerPrefix(candy_token.IF, p.parseIfExpression)
	p.registerPrefix(candy_token.NEW, p.parseNewExpression)
	p.registerPrefix(candy_token.AMPERSAND, p.parsePrefixExpression)
	p.registerPrefix(candy_token.ASTERISK, p.parsePrefixExpression)
	p.registerPrefix(candy_token.NOT, p.parsePrefixExpression)
	p.registerPrefix(candy_token.BIT_NOT, p.parsePrefixExpression)
	p.registerPrefix(candy_token.AWAIT, p.parseAwaitExpression)
	p.registerPrefix(candy_token.TYPEOF, p.parseTypeofExpression)
	p.registerPrefix(candy_token.INC, p.parsePrefixExpression)
	p.registerPrefix(candy_token.DEC, p.parsePrefixExpression)

	p.registerInfix(candy_token.LBRACK, p.parseIndexExpression)
	p.registerInfix(candy_token.PLUS, p.parseInfixExpression)
	p.registerInfix(candy_token.MINUS, p.parseInfixExpression)
	p.registerInfix(candy_token.SLASH, p.parseInfixExpression)
	p.registerInfix(candy_token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(candy_token.PERCENT, p.parseInfixExpression)
	p.registerInfix(candy_token.EQ, p.parseInfixExpression)
	p.registerInfix(candy_token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(candy_token.LT, p.parseInfixExpression)
	p.registerInfix(candy_token.GT, p.parseInfixExpression)
	p.registerInfix(candy_token.LTE, p.parseInfixExpression)
	p.registerInfix(candy_token.GTE, p.parseInfixExpression)
	p.registerInfix(candy_token.IN, p.parseInfixExpression)
	p.registerInfix(candy_token.LPAREN, p.parseCallExpression)
	p.registerInfix(candy_token.DOT, p.parseDotExpression)
	p.registerInfix(candy_token.ARROW, p.parseLambdaInfix)
	p.registerInfix(candy_token.LBRACE, p.parseStructLiteralExpr)
	p.registerInfix(candy_token.IS, p.parseIsExpression)
	p.registerInfix(candy_token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(candy_token.PLUS_ASSIGN, p.parseAssignExpression)
	p.registerInfix(candy_token.MINUS_ASSIGN, p.parseAssignExpression)
	p.registerInfix(candy_token.STAR_ASSIGN, p.parseAssignExpression)
	p.registerInfix(candy_token.SLASH_ASSIGN, p.parseAssignExpression)
	p.registerInfix(candy_token.SAFE_DOT, p.parseDotExpression)
	p.registerInfix(candy_token.NULL_COALESCE, p.parseInfixExpression)
	p.registerInfix(candy_token.RANGE, p.parseRangeExpression)
	p.registerInfix(candy_token.RANGE_EXCL, p.parseRangeExpression)
	p.registerInfix(candy_token.QUESTION, p.parseTernaryExpression)
	p.registerInfix(candy_token.AND, p.parseInfixExpression)
	p.registerInfix(candy_token.OR, p.parseInfixExpression)
	p.registerInfix(candy_token.AND_AND, p.parseInfixExpression)
	p.registerInfix(candy_token.OR_OR, p.parseInfixExpression)
	p.registerInfix(candy_token.AMPERSAND, p.parseInfixExpression)
	p.registerInfix(candy_token.BIT_OR, p.parseInfixExpression)
	p.registerInfix(candy_token.BIT_XOR, p.parseInfixExpression)
	p.registerInfix(candy_token.SHL, p.parseInfixExpression)
	p.registerInfix(candy_token.SHR, p.parseInfixExpression)
	p.registerInfix(candy_token.COMMA, p.parseTupleInfix)
	p.registerInfix(candy_token.INC, p.parsePostfixIncDec)
	p.registerInfix(candy_token.DEC, p.parsePostfixIncDec)

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(t candy_token.TokenType, fn prefixParseFn) { p.prefixParseFns[t] = fn }
func (p *Parser) registerInfix(t candy_token.TokenType, fn infixParseFn)   { p.infixParseFns[t] = fn }

func (p *Parser) Errors() []candy_report.Diagnostic { return p.errors }
func (p *Parser) peekError(t candy_token.TokenType) {
	p.addErrorAt(p.peekToken, fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
}
func (p *Parser) addErrorf(f string, a ...any) {
	p.addErrorAt(p.curToken, fmt.Sprintf(f, a...))
}
func (p *Parser) addErrorAt(tok candy_token.Token, msg string) {
	full := fmt.Sprintf("%d:%d: %s", tok.Line, tok.Col, msg)
	if p.errorSeen == nil {
		p.errorSeen = make(map[string]struct{})
	}
	if _, ok := p.errorSeen[full]; ok {
		return
	}
	p.errorSeen[full] = struct{}{}
	p.errors = append(p.errors, candy_report.Diagnostic{
		Level:   candy_report.Error,
		Message: msg,
		Line:    tok.Line,
		Col:     tok.Col,
		Offset:  tok.Offset,
		Length:  len(tok.Literal),
	})
}
func (p *Parser) addRecoveryf(f string, a ...any) {
	p.recovery = append(p.recovery, fmt.Sprintf(f, a...))
}
func (p *Parser) curTokenIs(t candy_token.TokenType) bool  { return p.curToken.Type == t }
func (p *Parser) peekTokenIs(t candy_token.TokenType) bool { return p.peekToken.Type == t }

func (p *Parser) expectPeek(t candy_token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) expect(t candy_token.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	p.addErrorf("expected %s, got %s", t, p.curToken.Type)
	return false
}

// ParseProgram parses a full program. Statement parsers consume up to the closing semicolon; no extra next here.
func (p *Parser) ParseProgram() *candy_ast.Program {
	prog := &candy_ast.Program{Statements: []candy_ast.Statement{}}
	// If recovery fails, avoid unbounded work.
	lastOffset, lastType, stall := -1, candy_token.TokenType(""), 0
	for {
		if p.curTokenIs(candy_token.EOF) {
			break
		}
		if p.curToken.Offset == lastOffset && p.curToken.Type == lastType {
			stall++
		} else {
			stall = 0
			lastOffset, lastType = p.curToken.Offset, p.curToken.Type
		}
		if stall > 4 {
			p.addRecoveryf("parser recovery stalled at %s, advancing", p.curToken.Type)
			p.nextToken()
			stall = 0
			continue
		}
		if len(p.errors) > 500 {
			p.addErrorf("too many parse errors, stopping")
			break
		}
		if st := p.parseStatement(); !isNilStatement(st) {
			prog.Statements = append(prog.Statements, st)
		} else {
			// error recovery: synchronize to likely statement boundary/start.
			if p.curTokenIs(candy_token.EOF) {
				break
			}
			p.synchronize()
		}
		// After a successful statement, expect peek left cur on `;` — move past.
		if p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
		}
	}
	return prog
}

func isNilStatement(st candy_ast.Statement) bool {
	if st == nil {
		return true
	}
	v := reflect.ValueOf(st)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func (p *Parser) synchronize() {
	for !p.curTokenIs(candy_token.EOF) {
		if p.isStatementStart(p.curToken.Type) {
			return
		}
		if p.curTokenIs(candy_token.SEMICOLON) {
			p.nextToken()
			return
		}
		if p.curTokenIs(candy_token.RBRACE) {
			p.nextToken()
			return
		}
		if p.isStatementStart(p.peekToken.Type) {
			p.nextToken()
			return
		}
		p.nextToken()
	}
}

func (p *Parser) isStatementStart(tt candy_token.TokenType) bool {
	switch tt {
	case candy_token.VAL, candy_token.VAR, candy_token.CONST, candy_token.PRINT, candy_token.RETURN, candy_token.FUNCTION, candy_token.IF, candy_token.IMPORT, candy_token.STRUCT, candy_token.PACKAGE, candy_token.CLASS, candy_token.OBJECT, candy_token.INTERFACE, candy_token.TRAIT, candy_token.EXTERN, candy_token.SEALED, candy_token.SUSPEND,
		candy_token.DIM, candy_token.SUB, candy_token.SELECT, candy_token.FOR, candy_token.WHILE, candy_token.DO, candy_token.SWITCH, candy_token.GOTO, candy_token.IDENT, candy_token.VOID,
		candy_token.REF, candy_token.SHARED, candy_token.MAYBE, candy_token.EXPORT, candy_token.MODULE, candy_token.ENUM, candy_token.TRY, candy_token.CATCH, candy_token.RUN, candy_token.ASYNC,
		candy_token.LBRACK, candy_token.PRIVATE, candy_token.OPERATOR, candy_token.FINALLY,
		candy_token.BREAK, candy_token.CONTINUE, candy_token.FOREACH, candy_token.DELETE,
		candy_token.REPEAT, candy_token.LOOP:
		return true
	default:
		return false
	}
}

func (p *Parser) peekPrecedence() int {
	if n, ok := precedences[p.peekToken.Type]; ok {
		return n
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if n, ok := precedences[p.curToken.Type]; ok {
		return n
	}
	return LOWEST
}

func (p *Parser) parseIntegerLiteral() candy_ast.Expression {
	v, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addErrorf("int: %v", err)
		return nil
	}
	return &candy_ast.IntegerLiteral{Token: p.curToken, Value: v}
}

func (p *Parser) parseFloatLiteral() candy_ast.Expression {
	v, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addErrorf("float: %v", err)
		return nil
	}
	return &candy_ast.FloatLiteral{Token: p.curToken, Value: v}
}

func (p *Parser) parseBoolean() candy_ast.Expression {
	return &candy_ast.Boolean{Token: p.curToken, Value: p.curTokenIs(candy_token.TRUE)}
}

func (p *Parser) parseNull() candy_ast.Expression {
	// use identifier-like node or dedicated - use a NullLiteral
	return &candy_ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) parseIdentifier() candy_ast.Expression {
	ident := &candy_ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekTokenIs(candy_token.LT) {
		// `peek` is already `<` from the lexer; the lexer is positioned to read the token *after* `<`.
		// Disambiguate `n < 3` from `Box<int>`: generic has `ident` (or `>`/`comma` chain) and closing `>`.
		br := p.l.Branch()
		t2 := br.NextToken()
		t3 := br.NextToken()
		if t2.Type == candy_token.INT || t2.Type == candy_token.FLOAT {
			return ident
		}
		if t2.Type == candy_token.IDENT && (t3.Type == candy_token.GT || t3.Type == candy_token.COMMA) {
			return p.parseGenericTypeExpression(ident)
		}
		return ident
	}
	return ident
}

func (p *Parser) parseGenericTypeExpression(ident *candy_ast.Identifier) candy_ast.Expression {
	tok := p.curToken
	p.nextToken() // <
	p.nextToken() // first arg
	var args []candy_ast.Expression
	for !p.curTokenIs(candy_token.GT) && !p.curTokenIs(candy_token.EOF) {
		args = append(args, p.parseTypeIdentifier())
		if p.curTokenIs(candy_token.GT) {
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
	// We leave it at GT so infix operators like { can pick it up
	return &candy_ast.TypeExpression{Token: tok, Name: ident, Arguments: args}
}
