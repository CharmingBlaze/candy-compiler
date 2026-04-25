package candy_ast

import (
	"candy/candy_token"
	"strings"
)

type ReturnStatement struct {
	Token       candy_token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out strings.Builder
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(StringExpr(rs.ReturnValue))
	}
	out.WriteString(";")
	return out.String()
}
