package candy_ast

import "candy/candy_token"

type TernaryExpression struct {
	Token       candy_token.Token // The '?' token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) expressionNode()      {}
func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TernaryExpression) String() string {
	return StringExpr(te.Condition) + " ? " + StringExpr(te.Consequence) + " : " + StringExpr(te.Alternative)
}
