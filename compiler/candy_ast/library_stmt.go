package candy_ast

import "candy/candy_token"

// LibraryStatement supports checklist-style source form:
// library "name" { ... }
type LibraryStatement struct {
	Token candy_token.Token
	Name  string
	Body  *BlockStatement
}

func (s *LibraryStatement) statementNode()       {}
func (s *LibraryStatement) TokenLiteral() string { return s.Token.Literal }
func (s *LibraryStatement) String() string {
	if s.Body == nil {
		return "library \"" + s.Name + "\" {}"
	}
	return "library \"" + s.Name + "\" " + s.Body.String()
}
