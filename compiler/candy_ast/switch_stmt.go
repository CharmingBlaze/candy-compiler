package candy_ast

import "candy/candy_token"

// SwitchCase is a single `case ...:` arm in a `switch`.
// Default is represented with IsDefault=true (Patterns may be nil).
type SwitchCase struct {
	Token     candy_token.Token
	IsDefault bool
	Patterns  []Expression
	Body      Statement
}

func (sc SwitchCase) String() string {
	if sc.IsDefault {
		return "default"
	}
	return "case"
}

// SwitchStatement parses C/BASIC-style switch with block bodies.
type SwitchStatement struct {
	Token   candy_token.Token
	Subject Expression
	Cases   []SwitchCase
}

func (s *SwitchStatement) statementNode()       {}
func (s *SwitchStatement) TokenLiteral() string { return s.Token.Literal }
