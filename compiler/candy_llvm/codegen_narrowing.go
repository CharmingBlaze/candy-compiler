package candy_llvm

import (
	"candy/candy_ast"
	"fmt"
)

func (c *Compiler) ensureTypeTag(typeName string) int {
	if tag, ok := c.typeTags[typeName]; ok {
		return tag
	}
	c.tagCount++
	c.typeTags[typeName] = c.tagCount
	return c.tagCount
}

func (c *Compiler) guardNarrowing(expr candy_ast.Expression) (*candy_ast.IsExpression, string, string, bool) {
	isExpr, ok := expr.(*candy_ast.IsExpression)
	if !ok || isExpr == nil {
		return nil, "", "", false
	}
	id, ok := isExpr.Left.(*candy_ast.Identifier)
	if !ok || id == nil || isExpr.TypeName == nil {
		return nil, "", "", false
	}
	return isExpr, id.Value, candy_ast.ExprAsSimpleTypeName(isExpr.TypeName), true
}

func (c *Compiler) pushNarrowing(m map[string]string) {
	c.narrowed = append(c.narrowed, m)
}

func (c *Compiler) popNarrowing() {
	if len(c.narrowed) == 0 {
		return
	}
	c.narrowed = c.narrowed[:len(c.narrowed)-1]
}

func (c *Compiler) lookupNarrowed(name string) (string, bool) {
	for i := len(c.narrowed) - 1; i >= 0; i-- {
		if t, ok := c.narrowed[i][name]; ok {
			return t, true
		}
	}
	return "", false
}

func (c *Compiler) emitCString(s string) string {
	c.strCount++
	id := fmt.Sprintf("dynname%d", c.strCount)
	ptr := c.addGlobalStr(id, s+"\\00")
	reg := c.nextReg()
	l := len(s) + 1
	c.emit("%s = getelementptr inbounds [%d x i8], [%d x i8]* %s, i64 0, i64 0", reg, l, l, ptr)
	return reg
}
