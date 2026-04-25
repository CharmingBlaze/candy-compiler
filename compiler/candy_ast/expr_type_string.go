package candy_ast

import "candy/candy_token"

// ExprToken returns the source token for a type expression, if known.
func ExprToken(e Expression) candy_token.Token {
	if e == nil {
		return candy_token.Token{}
	}
	if id, ok := e.(*Identifier); ok {
		return id.Token
	}
	if te, ok := e.(*TypeExpression); ok {
		return te.Token
	}
	return candy_token.Token{}
}

// ExprAsSimpleTypeName returns the identifier text for a type written as a plain name
// (e.g. `int`, `Vector2`). Returns "" for nil or non-identifier type expressions.
func ExprAsSimpleTypeName(e Expression) string {
	if e == nil {
		return ""
	}
	if id, ok := e.(*Identifier); ok {
		return id.Value
	}
	if te, ok := e.(*TypeExpression); ok {
		if te.ResolvedName != "" {
			return te.ResolvedName
		}
		return te.String()
	}
	return ""
}
