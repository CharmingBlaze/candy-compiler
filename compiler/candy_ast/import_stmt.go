package candy_ast

import "candy/candy_token"

// ImportStatement: import "path";
type ImportStatement struct {
	Token candy_token.Token
	Path  string // unquoted path string
}

func (i *ImportStatement) statementNode()       {}
func (i *ImportStatement) TokenLiteral() string { return i.Token.Literal }
func (i *ImportStatement) String() string       { return "import \"" + i.Path + "\";" }
