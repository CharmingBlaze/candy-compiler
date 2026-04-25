package candy_ast

import "candy/candy_token"

// WhenArm is one condition → body in a `when` expression.
type WhenArm struct {
	Cond Expression
	Body Expression
}

// WhenExpression is a simplified `when` (arms evaluated in order; first true cond runs body).
type WhenExpression struct {
	Token candy_token.Token
	Arms  []WhenArm
	ElseV Expression
}

func (w *WhenExpression) expressionNode()      {}
func (w *WhenExpression) TokenLiteral() string { return w.Token.Literal }
func (w *WhenExpression) String() string       { return "when" }
