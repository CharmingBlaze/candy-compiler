package candy_ast

import "candy/candy_token"

// IndexExpression is base [ index ].
type IndexExpression struct {
	Token candy_token.Token
	Base  Expression
	Index Expression
	IsSafe bool
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return StringExpr(ie.Base) + "[" + StringExpr(ie.Index) + "]"
}
