package candy_ast

import "candy/candy_token"

type BlockStatement struct {
	Token      candy_token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	s := ""
	for _, st := range bs.Statements {
		if st != nil {
			s += StringStmt(st) + " "
		}
	}
	return s
}
