package candy_ast

import "candy/candy_token"

type ExpressionStatement struct {
	Token      candy_token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return StringExpr(es.Expression) + ";"
	}
	return ""
}
