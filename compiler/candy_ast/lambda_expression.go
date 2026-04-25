package candy_ast

import (
	"candy/candy_token"
)

type LambdaExpression struct {
	Token      candy_token.Token // The ( token
	Parameters []Parameter
	Body       Expression
}

func (le *LambdaExpression) expressionNode()      {}
func (le *LambdaExpression) statementNode()       {}
func (le *LambdaExpression) TokenLiteral() string { return le.Token.Literal }
func (le *LambdaExpression) String() string       { return "lambda" }
