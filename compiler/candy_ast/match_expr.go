package candy_ast

import "candy/candy_token"

// MatchBranch is one `pattern => body` (pattern is an expression in the minimal design).
type MatchBranch struct {
	Pat   Expression
	Guard Expression
	Body  Expression
}

// MatchExpression is `match (subject) { pat => e; ... }` (minimal).
type MatchExpression struct {
	Token    candy_token.Token
	Subject  Expression
	Branches []MatchBranch
	Default  Expression
}

func (m *MatchExpression) expressionNode()      {}
func (m *MatchExpression) TokenLiteral() string { return m.Token.Literal }
func (m *MatchExpression) String() string       { return "match" }
