package candy_ast

import "candy/candy_token"

type PackageStatement struct {
	Token candy_token.Token
	Name  string
}

func (s *PackageStatement) statementNode()       {}
func (s *PackageStatement) TokenLiteral() string { return s.Token.Literal }

type ClassStatement struct {
	Token          candy_token.Token
	Sealed         bool
	Name           *Identifier
	TypeParameters []*Identifier
	Parameters     []Parameter // Primary constructor
	Base           *Identifier
	Traits         []*Identifier
	Members        []Statement // Fields and Methods
}

func (s *ClassStatement) statementNode()       {}
func (s *ClassStatement) TokenLiteral() string { return s.Token.Literal }

type ObjectStatement struct {
	Token   candy_token.Token
	Name    *Identifier
	Base    *Identifier
	Members []Statement
}

func (s *ObjectStatement) statementNode()       {}
func (s *ObjectStatement) TokenLiteral() string { return s.Token.Literal }

type InterfaceMethod struct {
	Name       *Identifier
	Parameters []Parameter
	ReturnType Expression
}

type InterfaceStatement struct {
	Token          candy_token.Token
	Name           *Identifier
	TypeParameters []*Identifier
	Methods        []*InterfaceMethod
}

func (s *InterfaceStatement) statementNode()       {}
func (s *InterfaceStatement) TokenLiteral() string { return s.Token.Literal }

type TraitStatement struct {
	Token candy_token.Token
	Name  *Identifier
}

func (s *TraitStatement) statementNode()       {}
func (s *TraitStatement) TokenLiteral() string { return s.Token.Literal }

type ExternFunctionStatement struct {
	Token    candy_token.Token
	Function *FunctionStatement
}

func (s *ExternFunctionStatement) statementNode()       {}
func (s *ExternFunctionStatement) TokenLiteral() string { return s.Token.Literal }
