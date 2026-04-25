package candy_ast

import "candy/candy_token"

type AwaitExpression struct {
	Token candy_token.Token // the 'await' token
	Value Expression        // the expression being awaited
}

func (ae *AwaitExpression) expressionNode()      {}
func (ae *AwaitExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AwaitExpression) String() string       { return "await " + StringExpr(ae.Value) }
