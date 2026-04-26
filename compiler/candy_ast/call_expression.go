package candy_ast

import (
	"candy/candy_token"
	"fmt"
	"strings"
)

type CallExpression struct {
	Token         candy_token.Token
	Function      Expression
	TypeArguments []Expression
	Arguments     []Expression
	IsSafe        bool
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var typeArgs []string
	for _, a := range ce.TypeArguments {
		typeArgs = append(typeArgs, StringExpr(a))
	}
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, StringExpr(a))
	}
	if len(typeArgs) > 0 {
		return fmt.Sprintf("%s<%s>(%s)", StringExpr(ce.Function), strings.Join(typeArgs, ", "), strings.Join(args, ", "))
	}
	return fmt.Sprintf("%s(%s)", StringExpr(ce.Function), strings.Join(args, ", "))
}
