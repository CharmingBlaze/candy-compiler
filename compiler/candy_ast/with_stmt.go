package candy_ast

import "candy/candy_token"

// WithStatement is resource-scope sugar:
// with x = open() { ... }
// Evaluator runs body then auto-cleans x.
type WithStatement struct {
	Token candy_token.Token
	Name  *Identifier
	Value Expression
	Body  *BlockStatement
}

func (s *WithStatement) statementNode()       {}
func (s *WithStatement) TokenLiteral() string { return s.Token.Literal }
func (s *WithStatement) String() string {
	return "with " + s.Name.String() + " = " + StringExpr(s.Value) + " " + s.Body.String()
}
