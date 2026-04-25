package candy_ast

import (
	"candy/candy_token"
)

type IsExpression struct {
	Token    candy_token.Token // The 'is' token
	Left     Expression
	TypeName Expression
}

func (ie *IsExpression) expressionNode()      {}
func (ie *IsExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IsExpression) String() string {
	return "(" + StringExpr(ie.Left) + " is " + StringExpr(ie.TypeName) + ")"
}
