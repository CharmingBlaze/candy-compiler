package candy_ast

import (
	"candy/candy_token"
)

type ForStatement struct {
	Token    candy_token.Token // The FOR token
	Var      *Identifier // The loop variable (key in map)
	ValueVar *Identifier // Optional second loop variable (value in map)
	Iterable Expression
	Start    Expression
	End      Expression
	Step     Expression // Optional
	Body     *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	if fs.Iterable != nil {
		return "FOR " + fs.Var.String() + " IN " + StringExpr(fs.Iterable) + "\n" + fs.Body.String() + "\nNEXT"
	}
	out := "FOR " + fs.Var.String() + " = " + StringExpr(fs.Start) + " TO " + StringExpr(fs.End)
	if fs.Step != nil {
		out += " STEP " + StringExpr(fs.Step)
	}
	out += "\n" + fs.Body.String() + "\nNEXT"
	return out
}
