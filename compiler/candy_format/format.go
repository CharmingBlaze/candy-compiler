package candy_format

import (
	"candy/candy_ast"
	"strings"
)

// Source formats a program into stable, newline-delimited source.
func Source(prog *candy_ast.Program) string {
	if prog == nil {
		return ""
	}
	var b strings.Builder
	for _, s := range prog.Statements {
		if s == nil {
			continue
		}
		writeStatement(&b, s, 0)
	}
	return b.String()
}

func writeIndent(b *strings.Builder, n int) {
	for i := 0; i < n; i++ {
		b.WriteString("    ")
	}
}

func writeStatement(b *strings.Builder, s candy_ast.Statement, ind int) {
	switch t := s.(type) {
	case *candy_ast.ValStatement:
		writeIndent(b, ind)
		b.WriteString("val ")
		b.WriteString(t.Name.Value)
		b.WriteString(" = ")
		b.WriteString(candy_ast.StringExpr(t.Value))
		b.WriteString(";\n")
	case *candy_ast.VarStatement:
		writeIndent(b, ind)
		b.WriteString("var ")
		b.WriteString(t.Name.Value)
		b.WriteString(" = ")
		b.WriteString(candy_ast.StringExpr(t.Value))
		b.WriteString(";\n")
	case *candy_ast.ReturnStatement:
		writeIndent(b, ind)
		b.WriteString("return")
		if t.ReturnValue != nil {
			b.WriteString(" ")
			b.WriteString(candy_ast.StringExpr(t.ReturnValue))
		}
		b.WriteString(";\n")
	case *candy_ast.ExpressionStatement:
		writeIndent(b, ind)
		b.WriteString(candy_ast.StringExpr(t.Expression))
		b.WriteString(";\n")
	case *candy_ast.ImportStatement:
		writeIndent(b, ind)
		b.WriteString(`import "`)
		b.WriteString(t.Path)
		b.WriteString("\";\n")
	case *candy_ast.StructStatement:
		writeIndent(b, ind)
		b.WriteString("struct ")
		b.WriteString(t.Name.Value)
		b.WriteString(" {\n")
		for _, f := range t.Fields {
			writeIndent(b, ind+1)
			b.WriteString(f.Name.Value)
			b.WriteString(": ")
			b.WriteString(candy_ast.StringExpr(f.TypeName))
			b.WriteString(";\n")
		}
		writeIndent(b, ind)
		b.WriteString("};\n")
	case *candy_ast.FunctionStatement:
		writeIndent(b, ind)
		b.WriteString("fun ")
		b.WriteString(t.Name.Value)
		b.WriteString("(")
		for i, p := range t.Parameters {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(p.Name.Value)
			b.WriteString(": ")
			b.WriteString(candy_ast.StringExpr(p.TypeName))
		}
		b.WriteString("): ")
		if t.ReturnType != nil {
			b.WriteString(candy_ast.StringExpr(t.ReturnType))
		} else {
			b.WriteString("Any")
		}
		b.WriteString(" {\n")
		if t.Body != nil {
			for _, bs := range t.Body.Statements {
				writeStatement(b, bs, ind+1)
			}
		}
		writeIndent(b, ind)
		b.WriteString("};\n")
	case *candy_ast.IfExpression:
		writeIndent(b, ind)
		b.WriteString("if (")
		b.WriteString(candy_ast.StringExpr(t.Condition))
		b.WriteString(") {\n")
		if t.Consequence != nil {
			for _, bs := range t.Consequence.Statements {
				writeStatement(b, bs, ind+1)
			}
		}
		writeIndent(b, ind)
		b.WriteString("}")
		if t.Alternative != nil {
			if altBlock, ok := t.Alternative.(*candy_ast.BlockStatement); ok {
				b.WriteString(" else {\n")
				for _, bs := range altBlock.Statements {
					writeStatement(b, bs, ind+1)
				}
				writeIndent(b, ind)
				b.WriteString("}")
			} else {
				b.WriteString(" else ")
				writeStatement(b, t.Alternative, ind)
				return
			}
		}
		b.WriteString(";\n")
	default:
		writeIndent(b, ind)
		b.WriteString(candy_ast.StringStmt(s))
		b.WriteString("\n")
	}
}
