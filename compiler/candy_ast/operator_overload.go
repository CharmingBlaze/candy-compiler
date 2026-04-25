package candy_ast

import "candy/candy_token"

type OperatorOverloadStatement struct {
	Token      candy_token.Token // the 'operator' token
	Operator   string            // e.g. "+", "-", "[]"
	ReturnType Expression
	Parameters []Parameter
	Body       *BlockStatement
}

func (oos *OperatorOverloadStatement) statementNode()       {}
func (oos *OperatorOverloadStatement) TokenLiteral() string { return oos.Token.Literal }
func (oos *OperatorOverloadStatement) String() string {
	return "operator" + oos.Operator + " (...) { ... }"
}
