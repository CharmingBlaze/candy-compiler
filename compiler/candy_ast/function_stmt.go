package candy_ast

import "candy/candy_token"

type FunctionStatement struct {
	Token      candy_token.Token
	Suspend    bool // Kotlin-style `suspend fun` (distinct from `async function`)
	IsAsync    bool
	IsOverride bool
	// Exported: `export` at module level (C-style `export T name(...) { }`); parser-only.
	Exported       bool
	Attributes     []*Attribute
	Receiver       *Parameter
	Name           *Identifier
	TypeParameters []*Identifier
	Parameters     []Parameter
	Variadic       bool
	ReturnType     Expression
	Body           *BlockStatement
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
