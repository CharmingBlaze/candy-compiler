package candy_ast

import "candy/candy_token"

type BreakStatement struct {
	Token candy_token.Token
}

func (s *BreakStatement) statementNode()       {}
func (s *BreakStatement) TokenLiteral() string { return s.Token.Literal }
func (s *BreakStatement) String() string       { return "break" }

type ContinueStatement struct {
	Token candy_token.Token
}

func (s *ContinueStatement) statementNode()       {}
func (s *ContinueStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ContinueStatement) String() string       { return "continue" }

type DeleteStatement struct {
	Token candy_token.Token
	Value Expression
}

func (s *DeleteStatement) statementNode()       {}
func (s *DeleteStatement) TokenLiteral() string { return s.Token.Literal }
func (s *DeleteStatement) String() string       { return "delete(" + StringExpr(s.Value) + ")" }

type ForEachStatement struct {
	Token    candy_token.Token
	Var      *Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (s *ForEachStatement) statementNode()       {}
func (s *ForEachStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ForEachStatement) String() string {
	return "foreach (" + s.Var.String() + " in " + StringExpr(s.Iterable) + ") " + s.Body.String()
}

type RepeatStatement struct {
	Token candy_token.Token
	Count Expression
	Body  *BlockStatement
}

func (s *RepeatStatement) statementNode()       {}
func (s *RepeatStatement) TokenLiteral() string { return s.Token.Literal }
func (s *RepeatStatement) String() string {
	return "repeat " + StringExpr(s.Count) + " " + s.Body.String()
}

type LoopStatement struct {
	Token candy_token.Token
	Body  *BlockStatement
}

func (s *LoopStatement) statementNode()       {}
func (s *LoopStatement) TokenLiteral() string { return s.Token.Literal }
func (s *LoopStatement) String() string {
	return "loop " + s.Body.String()
}
