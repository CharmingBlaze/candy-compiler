package candy_ast

import "candy/candy_token"

type StringLiteral struct {
	Token candy_token.Token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return s.Token.Literal }
