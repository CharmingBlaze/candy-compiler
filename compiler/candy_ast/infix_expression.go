package candy_ast

import (
	"candy/candy_token"
	"fmt"
)

type InfixExpression struct {
	Token    candy_token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", StringExpr(ie.Left), ie.Operator, StringExpr(ie.Right))
}
