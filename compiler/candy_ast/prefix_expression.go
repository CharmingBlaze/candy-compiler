package candy_ast

import "candy/candy_token"

type PrefixExpression struct {
	Token    candy_token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	return "(" + p.Operator + StringExpr(p.Right) + ")"
}
