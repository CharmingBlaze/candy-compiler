package candy_ast

import (
	"candy/candy_token"
	"strings"
)

type CatchClause struct {
	Token      candy_token.Token // the 'catch' token
	Type       Expression        // The exception type
	Identifier *Identifier       // The variable name for the exception
	Body       *BlockStatement
}

type TryStatement struct {
	Token        candy_token.Token // the 'try' token
	TryBody      *BlockStatement
	CatchClauses []*CatchClause
	FinallyBody  *BlockStatement // Optional
}

func (ts *TryStatement) statementNode()       {}
func (ts *TryStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TryStatement) String() string {
	var sb strings.Builder
	sb.WriteString("try ")
	sb.WriteString(ts.TryBody.String())
	for _, cc := range ts.CatchClauses {
		sb.WriteString(" catch ")
		sb.WriteString(StringExpr(cc.Type))
		if cc.Identifier != nil {
			sb.WriteString(" ")
			sb.WriteString(cc.Identifier.String())
		}
		sb.WriteString(" ")
		sb.WriteString(cc.Body.String())
	}
	if ts.FinallyBody != nil {
		sb.WriteString(" finally ")
		sb.WriteString(ts.FinallyBody.String())
	}
	return sb.String()
}
