package candy_ast

import "candy/candy_token"

// MapPair is one k: v in a map literal.
type MapPair struct {
	Key   Expression
	Value Expression
}

// MapLiteral is `map { "k": v, ... }` (parsed after keyword `new` optional — here `map` + `{` ... `}`).
type MapLiteral struct {
	Token candy_token.Token
	Pairs []MapPair
}

func (m *MapLiteral) expressionNode()      {}
func (m *MapLiteral) TokenLiteral() string { return m.Token.Literal }
func (m *MapLiteral) String() string       { return "map{…}" }
