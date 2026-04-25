package candy_ast

import "strings"

// MonomorphStructName matches candy_typecheck specialisation naming: base + "_" + each type arg
// (e.g. List<int> -> List_int). Used by the typechecker and LLVM for generic struct literals.
func MonomorphStructName(base string, typeArgs []Expression) string {
	s := base
	for _, a := range typeArgs {
		arg := ExprAsSimpleTypeName(a)
		s += "_" + strings.ReplaceAll(strings.ReplaceAll(arg, " ", "_"), "*", "Ptr")
	}
	return s
}
