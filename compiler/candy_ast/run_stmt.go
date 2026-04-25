package candy_ast

import "candy/candy_token"

type RunStatement struct {
	Token candy_token.Token // the 'run' token
	Value Expression        // the expression being run (usually a call)
}

func (rs *RunStatement) statementNode()       {}
func (rs *RunStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *RunStatement) String() string       { return "run " + StringExpr(rs.Value) }
