package candy_ast

import (
	"candy/candy_token"
	"testing"
)

func TestMonomorphStructName(t *testing.T) {
	floatT := &Identifier{Token: candy_token.Token{Type: candy_token.IDENT, Literal: "float"}, Value: "float"}
	if got := MonomorphStructName("Box", []Expression{floatT}); got != "Box_float" {
		t.Fatalf("MonomorphStructName = %q, want Box_float", got)
	}
}
