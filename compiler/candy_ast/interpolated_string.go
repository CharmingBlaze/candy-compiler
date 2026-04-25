package candy_ast

import (
	"candy/candy_token"
	"strings"
)

type InterpolatedStringLiteral struct {
	Token candy_token.Token
	Parts []Expression // Can be StringLiteral or any other Expression
}

func (isl *InterpolatedStringLiteral) expressionNode()      {}
func (isl *InterpolatedStringLiteral) TokenLiteral() string { return isl.Token.Literal }
func (isl *InterpolatedStringLiteral) String() string {
	var out strings.Builder
	out.WriteString("\"")
	for _, p := range isl.Parts {
		if s, ok := p.(*StringLiteral); ok {
			out.WriteString(s.Value)
		} else {
			out.WriteString("{")
			out.WriteString(StringExpr(p))
			out.WriteString("}")
		}
	}
	out.WriteString("\"")
	return out.String()
}
