package candy_ast

import "candy/candy_token"

type Identifier struct {
	Token     candy_token.Token
	Value     string
	IsPointer bool
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
