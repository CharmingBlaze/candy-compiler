package candy_lexer

import (
	"candy/candy_token"
	"fmt"
	"strconv"
	"strings"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	col          int
	prevToken    candy_token.TokenType
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0}
	l.readChar()
	return l
}

// Branch returns a lexer with the same scan state, for lookahead without mutating the main lexer.
func (l *Lexer) Branch() *Lexer {
	n := *l
	return &n
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.line++
		l.col = 0
	} else if l.ch != 0 {
		l.col++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) peek2Char() byte {
	if l.readPosition+1 >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+1]
}

// NextToken returns the next token.
func (l *Lexer) NextToken() candy_token.Token {
	l.skipSpaceAndComments()

	tLine, tCol, tOff := l.line, l.col, l.position

	if l.ch == 0 {
		if l.canInsertSemicolon() {
			l.prevToken = candy_token.SEMICOLON
			return candy_token.Token{Type: candy_token.SEMICOLON, Literal: ";", Line: tLine, Col: tCol, Offset: tOff}
		}
		return candy_token.Token{Type: candy_token.EOF, Literal: "", Line: tLine, Col: tCol, Offset: tOff}
	}

	var tok candy_token.Token
	switch l.ch {
	case '"':
		l.readChar() // opening quote
		s, err := l.readStringContent('"')
		if err != nil {
			tok = candy_token.Token{Type: candy_token.ILLEGAL, Literal: err.Error(), Line: tLine, Col: tCol, Offset: tOff}
		} else if l.ch != '"' {
			tok = candy_token.Token{Type: candy_token.ILLEGAL, Literal: "unclosed string", Line: tLine, Col: tCol, Offset: tOff}
		} else {
			l.readChar() // closing "
			tok = candy_token.Token{Type: candy_token.STR, Literal: s, Line: tLine, Col: tCol, Offset: tOff}
		}
	case '\'':
		l.readChar() // opening quote
		s, err := l.readStringContent('\'')
		if err != nil {
			tok = candy_token.Token{Type: candy_token.ILLEGAL, Literal: err.Error(), Line: tLine, Col: tCol, Offset: tOff}
		} else if l.ch != '\'' {
			tok = candy_token.Token{Type: candy_token.ILLEGAL, Literal: "unclosed string", Line: tLine, Col: tCol, Offset: tOff}
		} else {
			l.readChar() // closing '
			tok = candy_token.Token{Type: candy_token.STR, Literal: s, Line: tLine, Col: tCol, Offset: tOff}
		}
	case '\n':
		// This happens if skipSpaceAndComments() stopped because canInsertSemicolon() was true
		tok = candy_token.Token{Type: candy_token.SEMICOLON, Literal: ";", Line: tLine, Col: tCol, Offset: tOff}
		l.readChar()
	case ';':
		tok = l.tok1(candy_token.SEMICOLON, tLine, tCol, tOff)
	case '=':
		if l.peekChar() == '=' {
			tok = l.tok2(candy_token.EQ, tLine, tCol, tOff)
		} else if l.peekChar() == '>' {
			tok = l.tok2(candy_token.ARROW, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.ASSIGN, tLine, tCol, tOff)
		}
	case '+':
		if l.peekChar() == '+' {
			tok = l.tok2(candy_token.INC, tLine, tCol, tOff)
		} else if l.peekChar() == '=' {
			tok = l.tok2(candy_token.PLUS_ASSIGN, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.PLUS, tLine, tCol, tOff)
		}
	case '-':
		if l.peekChar() == '-' {
			tok = l.tok2(candy_token.DEC, tLine, tCol, tOff)
		} else if l.peekChar() == '=' {
			tok = l.tok2(candy_token.MINUS_ASSIGN, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.MINUS, tLine, tCol, tOff)
		}
	case '!':
		if l.peekChar() == '=' {
			tok = l.tok2(candy_token.NOT_EQ, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.BANG, tLine, tCol, tOff)
		}
	case '/':
		if l.peekChar() == '=' {
			tok = l.tok2(candy_token.SLASH_ASSIGN, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.SLASH, tLine, tCol, tOff)
		}
	case '*':
		if l.peekChar() == '=' {
			tok = l.tok2(candy_token.STAR_ASSIGN, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.ASTERISK, tLine, tCol, tOff)
		}
	case '%':
		tok = l.tok1(candy_token.PERCENT, tLine, tCol, tOff)
	case '<':
		if l.peekChar() == '<' {
			tok = l.tok2(candy_token.SHL, tLine, tCol, tOff)
		} else if l.peekChar() == '=' {
			tok = l.tok2(candy_token.LTE, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.LT, tLine, tCol, tOff)
		}
	case '>':
		if l.peekChar() == '>' {
			tok = l.tok2(candy_token.SHR, tLine, tCol, tOff)
		} else if l.peekChar() == '=' {
			tok = l.tok2(candy_token.GTE, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.GT, tLine, tCol, tOff)
		}
	case ':':
		tok = l.tok1(candy_token.COLON, tLine, tCol, tOff)
	case '?':
		if l.peekChar() == '.' {
			tok = l.tok2(candy_token.SAFE_DOT, tLine, tCol, tOff)
		} else if l.peekChar() == '?' {
			if l.peek2Char() == '=' {
				tok = l.tok3(candy_token.NULL_COALESCE_ASSIGN, tLine, tCol, tOff)
			} else {
				tok = l.tok2(candy_token.NULL_COALESCE, tLine, tCol, tOff)
			}
		} else {
			tok = l.tok1(candy_token.QUESTION, tLine, tCol, tOff)
		}
	case '&':
		if l.peekChar() == '&' {
			tok = l.tok2(candy_token.AND_AND, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.AMPERSAND, tLine, tCol, tOff)
		}
	case '|':
		if l.peekChar() == '|' {
			if l.peek2Char() == '=' {
				tok = l.tok3(candy_token.OR_ASSIGN, tLine, tCol, tOff)
			} else {
				tok = l.tok2(candy_token.OR_OR, tLine, tCol, tOff)
			}
		} else if l.peekChar() == '>' {
			tok = l.tok2(candy_token.PIPELINE, tLine, tCol, tOff)
		} else {
			tok = l.tok1(candy_token.BIT_OR, tLine, tCol, tOff)
		}
	case '`':
		if l.peekChar() == '`' {
			if l.readPosition+1 < len(l.input) && l.input[l.readPosition+1] == '`' {
				tok = l.readTripleBacktickString(tLine, tCol, tOff)
			} else {
				tok = l.tok1(candy_token.ILLEGAL, tLine, tCol, tOff)
			}
		} else {
			tok = l.tok1(candy_token.ILLEGAL, tLine, tCol, tOff)
		}
	case '^':
		tok = l.tok1(candy_token.BIT_XOR, tLine, tCol, tOff)
	case '~':
		tok = l.tok1(candy_token.BIT_NOT, tLine, tCol, tOff)
	case '(':
		tok = l.tok1(candy_token.LPAREN, tLine, tCol, tOff)
	case ')':
		tok = l.tok1(candy_token.RPAREN, tLine, tCol, tOff)
	case ',':
		tok = l.tok1(candy_token.COMMA, tLine, tCol, tOff)
	case '{':
		tok = l.tok1(candy_token.LBRACE, tLine, tCol, tOff)
	case '}':
		tok = l.tok1(candy_token.RBRACE, tLine, tCol, tOff)
	case '[':
		tok = l.tok1(candy_token.LBRACK, tLine, tCol, tOff)
	case ']':
		tok = l.tok1(candy_token.RBRACK, tLine, tCol, tOff)
	case '.':
		if l.peekChar() == '.' {
			// Support both inclusive `..` and exclusive `..<` range operators.
			if l.readPosition+1 < len(l.input) && l.input[l.readPosition+1] == '<' {
				tok = l.tok3(candy_token.RANGE_EXCL, tLine, tCol, tOff)
			} else {
				tok = l.tok2(candy_token.RANGE, tLine, tCol, tOff)
			}
		} else {
			tok = l.tok1(candy_token.DOT, tLine, tCol, tOff)
		}
	default:
		if isLetter(l.ch) {
			raw := l.readIdentifier()
			id := strings.ToLower(raw)
			ty := candy_token.LookupIdent(id)
			tok = candy_token.Token{Type: ty, Literal: id, Line: tLine, Col: tCol, Offset: tOff}
		} else if isDigit(l.ch) {
			tok = l.readNumber(tLine, tCol, tOff)
		} else {
			b := l.ch
			l.readChar()
			tok = candy_token.Token{Type: candy_token.ILLEGAL, Literal: string(b), Line: tLine, Col: tCol, Offset: tOff}
		}
	}

	l.prevToken = tok.Type
	return tok
}

func (l *Lexer) canInsertSemicolon() bool {
	switch l.prevToken {
	case candy_token.IDENT, candy_token.INT, candy_token.FLOAT, candy_token.STR,
		candy_token.TRUE, candy_token.FALSE, candy_token.NULL, candy_token.RETURN,
		candy_token.RPAREN, candy_token.RBRACE, candy_token.RBRACK,
		candy_token.END, candy_token.WEND, candy_token.NEXT,
		candy_token.FUNCTION, candy_token.SUB, candy_token.IF,
		candy_token.CLASS, candy_token.STRUCT, candy_token.OBJECT,
		candy_token.INTERFACE, candy_token.TRAIT,
		candy_token.MODULE, candy_token.ENUM,
		candy_token.BREAK, candy_token.CONTINUE:
		return true
	}
	return false
}

func (l *Lexer) tok1(ty candy_token.TokenType, line, col, off int) candy_token.Token {
	t := candy_token.Token{Type: ty, Literal: string(l.ch), Line: line, Col: col, Offset: off}
	l.readChar()
	return t
}

func (l *Lexer) tok2(ty candy_token.TokenType, line, col, off int) candy_token.Token {
	ch1 := l.ch
	l.readChar()
	ch2 := l.ch
	l.readChar()
	return candy_token.Token{Type: ty, Literal: string(ch1) + string(ch2), Line: line, Col: col, Offset: off}
}

func (l *Lexer) tok3(ty candy_token.TokenType, line, col, off int) candy_token.Token {
	ch1 := l.ch
	l.readChar()
	ch2 := l.ch
	l.readChar()
	ch3 := l.ch
	l.readChar()
	return candy_token.Token{Type: ty, Literal: string(ch1) + string(ch2) + string(ch3), Line: line, Col: col, Offset: off}
}

func (l *Lexer) skipSpaceAndComments() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
			l.readChar()
			continue
		}
		if l.ch == '\n' {
			if l.canInsertSemicolon() {
				// Don't consume it; NextToken will return SEMICOLON
				// Wait, if I don't consume it, NextToken will call skipSpaceAndComments again.
				// I should transform this \n into a SEMICOLON token.
				// To do that, I'll change l.ch to ';' (virtually) or just return from here.
				return
			}
			l.readChar()
			continue
		}

		if l.ch == '/' && l.peekChar() == '/' {
			for l.ch != 0 && l.ch != '\n' {
				l.readChar()
			}
			continue
		}
		if l.ch == '/' && l.peekChar() == '*' {
			l.readChar()
			l.readChar()
			for l.ch != 0 {
				if l.ch == '*' && l.peekChar() == '/' {
					l.readChar()
					l.readChar()
					break
				}
				l.readChar()
			}
			continue
		}
		break
	}
}

func (l *Lexer) readStringContent(quote byte) (string, error) {
	var b strings.Builder
	for l.ch != quote && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case '"':
				b.WriteByte('"')
			case '\'':
				b.WriteByte('\'')
			case '{':
				// Keep `\` + `{` in the literal so the string interpolator can treat it as a literal brace.
				b.WriteByte('\\')
				b.WriteByte('{')
			case '\\':
				b.WriteByte('\\')
			case '0':
				b.WriteByte(0)
			case 0:
				return b.String(), fmt.Errorf("unclosed string")
			default:
				return b.String(), fmt.Errorf("bad escape")
			}
		} else {
			b.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == 0 {
		return b.String(), fmt.Errorf("unclosed string")
	}
	return b.String(), nil
}

func (l *Lexer) readNumber(line, col, off int) candy_token.Token {
	start := l.position
	isFloat := false
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // .
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	if l.ch == 'e' || l.ch == 'E' {
		isFloat = true
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		if !isDigit(l.ch) {
			lit := l.input[start:l.position]
			return candy_token.Token{Type: candy_token.ILLEGAL, Literal: lit, Line: line, Col: col, Offset: off}
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	lit := l.input[start:l.position]
	if isFloat {
		if _, err := strconv.ParseFloat(lit, 64); err != nil {
			return candy_token.Token{Type: candy_token.ILLEGAL, Literal: lit, Line: line, Col: col, Offset: off}
		}
		return candy_token.Token{Type: candy_token.FLOAT, Literal: lit, Line: line, Col: col, Offset: off}
	}
	return candy_token.Token{Type: candy_token.INT, Literal: lit, Line: line, Col: col, Offset: off}
}

func (l *Lexer) readIdentifier() string {
	s := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	// QBASIC style type suffixes (excluding & which is now address-of)
	if l.ch == '$' || l.ch == '%' || l.ch == '#' || l.ch == '!' {
		l.readChar()
	}
	return l.input[s:l.position]
}

func (l *Lexer) readTripleBacktickString(line, col, off int) candy_token.Token {
	l.readChar() // `
	l.readChar() // `
	l.readChar() // `
	start := l.position
	for {
		if l.ch == 0 {
			return candy_token.Token{Type: candy_token.ILLEGAL, Literal: "unclosed triple-backtick string", Line: line, Col: col, Offset: off}
		}
		if l.ch == '`' && l.peekChar() == '`' {
			if l.readPosition+1 < len(l.input) && l.input[l.readPosition+1] == '`' {
				lit := l.input[start:l.position]
				l.readChar() // `
				l.readChar() // `
				l.readChar() // `
				return candy_token.Token{Type: candy_token.STR, Literal: lit, Line: line, Col: col, Offset: off}
			}
		}
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool { return '0' <= ch && ch <= '9' }
