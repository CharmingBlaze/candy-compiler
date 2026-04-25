package candy_ast

import "candy/candy_token"

// PostfixExpression is `x++` or `x--`.
type PostfixExpression struct {
	Token    candy_token.Token
	Left     Expression
	Operator string
}

func (p *PostfixExpression) expressionNode()      {}
func (p *PostfixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PostfixExpression) String() string {
	return "(" + StringExpr(p.Left) + p.Operator + ")"
}
