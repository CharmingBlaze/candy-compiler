package candy_ast

import (
	"candy/candy_token"
	"reflect"
)

// Node is the base for all AST nodes.
type Node interface {
	TokenLiteral() string
}

// GetToken tries to extract a candy_token.Token from any AST node via reflection.
func GetToken(n Node) candy_token.Token {
	if n == nil {
		return candy_token.Token{}
	}
	v := reflect.ValueOf(n)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return candy_token.Token{}
	}
	f := v.FieldByName("Token")
	if f.IsValid() && f.Type().Name() == "Token" {
		if t, ok := f.Interface().(candy_token.Token); ok {
			return t
		}
	}
	return candy_token.Token{}
}

// Statement is a statement node.
type Statement interface {
	Node
	statementNode()
}

// Expression is an expression node.
type Expression interface {
	Node
	expressionNode()
}
