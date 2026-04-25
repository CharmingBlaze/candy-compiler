package candy_ast

import "candy/candy_token"

type DoWhileStatement struct {
	Token     candy_token.Token
	Body      *BlockStatement
	Condition Expression
}

func (s *DoWhileStatement) statementNode()       {}
func (s *DoWhileStatement) TokenLiteral() string { return s.Token.Literal }
