package candy_ast

import "candy/candy_token"

type PropertyStatement struct {
	Token        candy_token.Token // the type token or first token of property
	Name         *Identifier
	Type         Expression
	Getter       *BlockStatement
	Setter       *BlockStatement
	IsAuto       bool
	DefaultValue Expression
}

func (ps *PropertyStatement) statementNode()       {}
func (ps *PropertyStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PropertyStatement) String() string       { return ps.Name.String() + " { ... }" }
