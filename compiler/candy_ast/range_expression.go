package candy_ast

import "candy/candy_token"

type RangeExpression struct {
	Token candy_token.Token // The .. token
	Left  Expression
	Right Expression
}

func (re *RangeExpression) expressionNode()      {}
func (re *RangeExpression) TokenLiteral() string { return re.Token.Literal }
func (re *RangeExpression) String() string {
	return StringExpr(re.Left) + ".." + StringExpr(re.Right)
}
