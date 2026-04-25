package candy_ast

import (
	"candy/candy_token"
	"fmt"
	"strings"
)

// TypeExpression represents a type usage, possibly generic: List<int>
type TypeExpression struct {
	Token        candy_token.Token
	Name         *Identifier
	Arguments    []Expression // Could be TypeExpressions or Identifiers
	ResolvedName string
}

func (te *TypeExpression) expressionNode()      {}
func (te *TypeExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TypeExpression) String() string {
	if len(te.Arguments) == 0 {
		return te.Name.Value
	}
	args := []string{}
	for _, a := range te.Arguments {
		args = append(args, StringExpr(a))
	}
	return fmt.Sprintf("%s<%s>", te.Name.Value, strings.Join(args, ", "))
}

type TupleTypeExpression struct {
	Token candy_token.Token
	Types []Expression
}

func (te *TupleTypeExpression) expressionNode()      {}
func (te *TupleTypeExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TupleTypeExpression) String() string {
	types := []string{}
	for _, t := range te.Types {
		types = append(types, StringExpr(t))
	}
	return "(" + strings.Join(types, ", ") + ")"
}
