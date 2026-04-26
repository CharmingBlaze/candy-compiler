package candy_lexer

import (
	"testing"

	"candy/candy_token"
)

type tokCase struct {
	expectedType    candy_token.TokenType
	expectedLiteral string
}

func expectTokens(t *testing.T, input string, cases []tokCase) {
	t.Helper()
	l := New(input)
	for i, tc := range cases {
		tok := l.NextToken()
		if tok.Type != tc.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tc.expectedType, tok.Type)
		}
		if tok.Literal != tc.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tc.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_FullProgram(t *testing.T) {
	t.Run("val_var_and_function", func(t *testing.T) {
		input := `val five = 5;
var ten = 10;

fun add(x: Int, y: Int): Int {
    return x + y;
};

val result = add(five, ten);
`
		cases := []tokCase{
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "five"},
			{candy_token.ASSIGN, "="},
			{candy_token.INT, "5"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.VAR, "var"},
			{candy_token.IDENT, "ten"},
			{candy_token.ASSIGN, "="},
			{candy_token.INT, "10"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.FUNCTION, "fun"},
			{candy_token.IDENT, "add"},
			{candy_token.LPAREN, "("},
			{candy_token.IDENT, "x"},
			{candy_token.COLON, ":"},
			{candy_token.IDENT, "int"},
			{candy_token.COMMA, ","},
			{candy_token.IDENT, "y"},
			{candy_token.COLON, ":"},
			{candy_token.IDENT, "int"},
			{candy_token.RPAREN, ")"},
			{candy_token.COLON, ":"},
			{candy_token.IDENT, "int"},
			{candy_token.LBRACE, "{"},
			{candy_token.RETURN, "return"},
			{candy_token.IDENT, "x"},
			{candy_token.PLUS, "+"},
			{candy_token.IDENT, "y"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.RBRACE, "}"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "result"},
			{candy_token.ASSIGN, "="},
			{candy_token.IDENT, "add"},
			{candy_token.LPAREN, "("},
			{candy_token.IDENT, "five"},
			{candy_token.COMMA, ","},
			{candy_token.IDENT, "ten"},
			{candy_token.RPAREN, ")"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.EOF, ""},
		}
		expectTokens(t, input, cases)
	})
	t.Run("ops_and_comparisons", func(t *testing.T) {
		input := `!- / *5;
5 < 10 > 5;
10 == 10;
10 != 9;`
		cases := []tokCase{
			{candy_token.BANG, "!"},
			{candy_token.MINUS, "-"},
			{candy_token.SLASH, "/"},
			{candy_token.ASTERISK, "*"},
			{candy_token.INT, "5"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.INT, "5"},
			{candy_token.LT, "<"},
			{candy_token.INT, "10"},
			{candy_token.GT, ">"},
			{candy_token.INT, "5"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.INT, "10"},
			{candy_token.EQ, "=="},
			{candy_token.INT, "10"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.INT, "10"},
			{candy_token.NOT_EQ, "!="},
			{candy_token.INT, "9"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.EOF, ""},
		}
		expectTokens(t, input, cases)
	})
	t.Run("if_else", func(t *testing.T) {
		input := `if (5 < 10) {
    return true;
} else {
    return false;
};
`
		cases := []tokCase{
			{candy_token.IF, "if"},
			{candy_token.LPAREN, "("},
			{candy_token.INT, "5"},
			{candy_token.LT, "<"},
			{candy_token.INT, "10"},
			{candy_token.RPAREN, ")"},
			{candy_token.LBRACE, "{"},
			{candy_token.RETURN, "return"},
			{candy_token.TRUE, "true"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.RBRACE, "}"},
			{candy_token.ELSE, "else"},
			{candy_token.LBRACE, "{"},
			{candy_token.RETURN, "return"},
			{candy_token.FALSE, "false"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.RBRACE, "}"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.EOF, ""},
		}
		expectTokens(t, input, cases)
	})
	t.Run("line_comment_skipped", func(t *testing.T) {
		input := `val a = 1; // cmt
return a;`
		cases := []tokCase{
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "a"},
			{candy_token.ASSIGN, "="},
			{candy_token.INT, "1"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.RETURN, "return"},
			{candy_token.IDENT, "a"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.EOF, ""},
		}
		expectTokens(t, input, cases)
	})
	t.Run("strings_and_floats", func(t *testing.T) {
		input := `val s = "hi";
val f = 1.25;
val e = 1.5e-2;
val e2 = 2E3;`
		cases := []tokCase{
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "s"},
			{candy_token.ASSIGN, "="},
			{candy_token.STR, "hi"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "f"},
			{candy_token.ASSIGN, "="},
			{candy_token.FLOAT, "1.25"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "e"},
			{candy_token.ASSIGN, "="},
			{candy_token.FLOAT, "1.5e-2"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "e2"},
			{candy_token.ASSIGN, "="},
			{candy_token.FLOAT, "2E3"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.EOF, ""},
		}
		expectTokens(t, input, cases)
	})
	t.Run("single_quoted_strings_and_empty", func(t *testing.T) {
		input := `val a = '';
val b = 'hi';
val c = "ok";`
		cases := []tokCase{
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "a"},
			{candy_token.ASSIGN, "="},
			{candy_token.STR, ""},
			{candy_token.SEMICOLON, ";"},
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "b"},
			{candy_token.ASSIGN, "="},
			{candy_token.STR, "hi"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "c"},
			{candy_token.ASSIGN, "="},
			{candy_token.STR, "ok"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.EOF, ""},
		}
		expectTokens(t, input, cases)
	})
	t.Run("nullable_type_suffix_token", func(t *testing.T) {
		input := `fun f(x: Int?): Int? { return null; };`
		cases := []tokCase{
			{candy_token.FUNCTION, "fun"},
			{candy_token.IDENT, "f"},
			{candy_token.LPAREN, "("},
			{candy_token.IDENT, "x"},
			{candy_token.COLON, ":"},
			{candy_token.IDENT, "int"},
			{candy_token.QUESTION, "?"},
			{candy_token.RPAREN, ")"},
			{candy_token.COLON, ":"},
			{candy_token.IDENT, "int"},
			{candy_token.QUESTION, "?"},
			{candy_token.LBRACE, "{"},
			{candy_token.RETURN, "return"},
			{candy_token.NULL, "null"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.RBRACE, "}"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.EOF, ""},
		}
		expectTokens(t, input, cases)
	})
	t.Run("block_comment", func(t *testing.T) {
		input := `val z = 2; /* x
y */ return z;`
		cases := []tokCase{
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "z"},
			{candy_token.ASSIGN, "="},
			{candy_token.INT, "2"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.RETURN, "return"},
			{candy_token.IDENT, "z"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.EOF, ""},
		}
		expectTokens(t, input, cases)
	})
	t.Run("case_insensitive_keywords_and_identifiers", func(t *testing.T) {
		input := `VaL MyVar = 3; ReTuRn mYvAr;`
		cases := []tokCase{
			{candy_token.VAL, "val"},
			{candy_token.IDENT, "myvar"},
			{candy_token.ASSIGN, "="},
			{candy_token.INT, "3"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.RETURN, "return"},
			{candy_token.IDENT, "myvar"},
			{candy_token.SEMICOLON, ";"},
			{candy_token.EOF, ""},
		}
		expectTokens(t, input, cases)
	})
}

func TestExclusiveRangeToken(t *testing.T) {
	cases := []tokCase{
		{candy_token.INT, "1"},
		{candy_token.RANGE_EXCL, "..<"},
		{candy_token.INT, "10"},
		{candy_token.SEMICOLON, ";"},
		{candy_token.EOF, ""},
	}
	expectTokens(t, "1..<10", cases)
}
