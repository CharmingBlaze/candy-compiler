package candy_ast

import (
	"candy/candy_token"
)

type CForStatement struct {
	Token candy_token.Token // The FOR token
	Init  Statement
	Cond  Expression
	Post  Expression
	Body  *BlockStatement
}

func (s *CForStatement) statementNode()       {}
func (s *CForStatement) TokenLiteral() string { return s.Token.Literal }
func (s *CForStatement) String() string {
	return "for ( ... ) { ... }"
}
