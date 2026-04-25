package candy_ast

import (
	"candy/candy_token"
	"strings"
)

type Attribute struct {
	Token     candy_token.Token // the '[' token
	Name      *Identifier
	Arguments map[string]Expression
}

func (a *Attribute) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(a.Name.String())
	if len(a.Arguments) > 0 {
		sb.WriteString("(")
		i := 0
		for k, v := range a.Arguments {
			sb.WriteString(k)
			sb.WriteString(": ")
			sb.WriteString(StringExpr(v))
			if i < len(a.Arguments)-1 {
				sb.WriteString(", ")
			}
			i++
		}
		sb.WriteString(")")
	}
	sb.WriteString("]")
	return sb.String()
}
