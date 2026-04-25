package candy_ast

import (
	"candy/candy_token"
	"strings"
)

type StructLiteral struct {
	Token  candy_token.Token // The { token
	Name   Expression
	Fields map[string]Expression
}

func (sl *StructLiteral) expressionNode()      {}
func (sl *StructLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StructLiteral) String() string {
	var out strings.Builder
	out.WriteString(StringExpr(sl.Name))
	out.WriteString("{ ")
	for k, v := range sl.Fields {
		out.WriteString(k + ": " + StringExpr(v) + ", ")
	}
	out.WriteString("}")
	return out.String()
}
