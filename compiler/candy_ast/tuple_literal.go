package candy_ast

import (
	"candy/candy_token"
	"strings"
)

type TupleLiteral struct {
	Token candy_token.Token
	Elems []Expression
}

func (t *TupleLiteral) expressionNode()      {}
func (t *TupleLiteral) TokenLiteral() string { return t.Token.Literal }
func (t *TupleLiteral) String() string {
	var out []string
	for _, e := range t.Elems {
		out = append(out, StringExpr(e))
	}
	return "(" + strings.Join(out, ", ") + ")"
}
