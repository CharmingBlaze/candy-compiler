package candy_ast

import "candy/candy_token"

type IfExpression struct {
	Token       candy_token.Token
	Condition   Expression
	Consequence Statement
	Alternative Statement
}

func (ie *IfExpression) statementNode()       {}
func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string       { return "if" }
