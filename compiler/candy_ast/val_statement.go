package candy_ast

import (
	"candy/candy_token"
	"strings"
)

type ValStatement struct {
	Token      candy_token.Token
	Attributes []*Attribute
	IsMaybe    bool
	IsShared   bool
	IsRef      bool
	Name       *Identifier
	TypeName   Expression
	Value      Expression
}

func (vs *ValStatement) statementNode()       {}
func (vs *ValStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *ValStatement) String() string {
	var out strings.Builder
	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	if vs.TypeName != nil {
		out.WriteString(": " + StringExpr(vs.TypeName))
	}
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(StringExpr(vs.Value))
	}
	out.WriteString(";")
	return out.String()
}
