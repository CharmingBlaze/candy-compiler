package candy_ast

import (
	"candy/candy_token"
	"fmt"
)

type GroupedExpression struct {
	Token candy_token.Token
	Expr  Expression
}

func (g *GroupedExpression) expressionNode()      {}
func (g *GroupedExpression) TokenLiteral() string { return g.Token.Literal }
func (g *GroupedExpression) String() string {
	return fmt.Sprintf("(%s)", StringExpr(g.Expr))
}
