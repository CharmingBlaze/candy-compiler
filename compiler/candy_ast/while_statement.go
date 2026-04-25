package candy_ast

import "candy/candy_token"

type WhileStatement struct {
	Token     candy_token.Token
	Condition Expression
	Body      *BlockStatement
}

func (s *WhileStatement) statementNode()       {}
func (s *WhileStatement) TokenLiteral() string { return s.Token.Literal }
func (s *WhileStatement) String() string {
	return "while " + StringExpr(s.Condition) + " { " + s.Body.String() + " }"
}
