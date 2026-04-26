package candy_ast

import "candy/candy_token"

// NamedArgumentExpression models `name: value` at call sites.
type NamedArgumentExpression struct {
	Token candy_token.Token
	Name  *Identifier
	Value Expression
}

func (n *NamedArgumentExpression) expressionNode()      {}
func (n *NamedArgumentExpression) TokenLiteral() string { return n.Token.Literal }
func (n *NamedArgumentExpression) String() string {
	if n == nil || n.Name == nil {
		return "<named-arg>"
	}
	return n.Name.Value + ": <expr>"
}
