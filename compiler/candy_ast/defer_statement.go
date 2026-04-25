package candy_ast

import "candy/candy_token"

type DeferStatement struct {
	Token candy_token.Token // The DEFER token
	Call  *CallExpression
}

func (s *DeferStatement) statementNode()       {}
func (s *DeferStatement) TokenLiteral() string { return s.Token.Literal }
func (s *DeferStatement) String() string       { return "defer " + s.Call.String() }
