package candy_ast

import (
	"candy/candy_token"
	"strings"
)

// Field is a struct field: name: Type
type Field struct {
	Token      candy_token.Token
	Attributes []*Attribute
	IsPrivate  bool
	Name       *Identifier
	TypeName   Expression
	Init       Expression
}

// StructStatement: struct Name : Base { ... }
type StructStatement struct {
	Token          candy_token.Token
	Attributes     []*Attribute
	Name           *Identifier
	TypeParameters []*Identifier
	Bases          []*Identifier // Inheritance and interfaces
	Fields         []Field
	Properties     []*PropertyStatement
	Methods        []*FunctionStatement
	Operators      []*OperatorOverloadStatement
}

func (s *StructStatement) statementNode()       {}
func (s *StructStatement) TokenLiteral() string { return s.Token.Literal }
func (s *StructStatement) String() string {
	var out strings.Builder
	out.WriteString("struct ")
	out.WriteString(s.Name.String())
	if len(s.TypeParameters) > 0 {
		out.WriteString("<")
		var params []string
		for _, p := range s.TypeParameters {
			params = append(params, p.String())
		}
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(">")
	}
	out.WriteString(" { ")
	for _, f := range s.Fields {
		out.WriteString(f.Name.String() + ": " + StringExpr(f.TypeName) + "; ")
	}
	out.WriteString("}")
	return out.String()
}
