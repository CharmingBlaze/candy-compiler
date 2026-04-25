package candy_ast

import "candy/candy_token"

type ModuleStatement struct {
	Token candy_token.Token // the 'module' token
	Name  *Identifier
	Body  *BlockStatement
}

func (ms *ModuleStatement) statementNode()       {}
func (ms *ModuleStatement) TokenLiteral() string { return ms.Token.Literal }
func (ms *ModuleStatement) String() string {
	return "module " + ms.Name.String() + " " + ms.Body.String()
}
