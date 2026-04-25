package candy_ast

import (
	"candy/candy_token"
	"strings"
)

type DotExpression struct {
	Token  candy_token.Token // The . token or ?. token
	Left   Expression
	Right  *Identifier
	IsSafe bool
}

func (de *DotExpression) expressionNode()      {}
func (de *DotExpression) TokenLiteral() string { return de.Token.Literal }
func (de *DotExpression) String() string {
	var out strings.Builder
	out.WriteString(StringExpr(de.Left))
	out.WriteString(".")
	out.WriteString(de.Right.String())
	return out.String()
}
