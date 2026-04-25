package candy_ast

import "candy/candy_token"

type IntegerLiteral struct {
	Token candy_token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return IntString(il.Value) }
