package candy_ast

import "candy/candy_token"

// ImportStatement: import "path";
type ImportStatement struct {
	Token   candy_token.Token
	Path    string   // unquoted path string
	Alias   string   // optional alias: import math as m
	From    string   // optional module path for from-import syntax
	Symbols []string // optional symbols for from-import: from math import sin, cos
}

func (i *ImportStatement) statementNode()       {}
func (i *ImportStatement) TokenLiteral() string { return i.Token.Literal }
func (i *ImportStatement) String() string {
	if i.From != "" {
		return "from " + i.From + " import ...;"
	}
	s := "import \"" + i.Path + "\""
	if i.Alias != "" {
		s += " as " + i.Alias
	}
	return s + ";"
}
