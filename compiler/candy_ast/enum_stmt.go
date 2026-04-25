package candy_ast

import (
	"candy/candy_token"
	"strings"
)

type EnumVariant struct {
	Name  *Identifier
	Value Expression // Optional, for explicitly valued enums
}

type EnumStatement struct {
	Token    candy_token.Token // the 'enum' token
	Name     *Identifier
	Variants []*EnumVariant
}

func (es *EnumStatement) statementNode()       {}
func (es *EnumStatement) TokenLiteral() string { return es.Token.Literal }
func (es *EnumStatement) String() string {
	var sb strings.Builder
	sb.WriteString("enum ")
	sb.WriteString(es.Name.String())
	sb.WriteString(" { ")
	for i, v := range es.Variants {
		sb.WriteString(v.Name.String())
		if v.Value != nil {
			sb.WriteString(" = ")
			sb.WriteString(StringExpr(v.Value))
		}
		if i < len(es.Variants)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" }")
	return sb.String()
}
