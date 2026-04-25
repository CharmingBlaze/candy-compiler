package candy_ast

import (
	"candy/candy_token"
)

type AssignExpression struct {
	Token    candy_token.Token // assignment token
	Operator string
	Left     Expression // Identifier, DotExpression, etc.
	Value    Expression
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) String() string {
	op := ae.Operator
	if op == "" {
		op = "="
	}
	return StringExpr(ae.Left) + " " + op + " " + StringExpr(ae.Value)
}
