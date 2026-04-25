package candy_ast

import "candy/candy_token"

// ArrayLiteral is [ expr, ... ].
type ArrayLiteral struct {
	Token candy_token.Token
	Elem  []Expression
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	return "[" + "…" + "]"
}
