package candy_typecheck

import (
	"candy/candy_ast"
	"fmt"
	"strings"
)

func (c *Checker) resolveTypeExpression(te *candy_ast.TypeExpression) string {
	if te == nil {
		return "any"
	}
	baseName := canonType(te.Name.Value)

	// Check if it's a generic template
	if template, ok := c.genericStructs[baseName]; ok {
		spec := c.specializeStruct(template, te.Arguments)
		if spec != nil && spec.Name != nil {
			te.ResolvedName = spec.Name.Value
			return canonType(spec.Name.Value)
		}
	}

	return te.String()
}

// specializeStruct creates a concrete StructStatement from a generic template.
func (c *Checker) specializeStruct(template *candy_ast.StructStatement, args []candy_ast.Expression) *candy_ast.StructStatement {
	if len(template.TypeParameters) != len(args) {
		c.add(fmt.Sprintf("wrong number of type arguments for %s: expected %d, got %d",
			template.Name.Value, len(template.TypeParameters), len(args)), template)
		return nil
	}

	mapping := make(map[string]candy_ast.Expression)
	for i, param := range template.TypeParameters {
		mapping[param.Value] = args[i]
	}

	// Unique name for the specialized struct: List<int> -> List_int
	newName := candy_ast.MonomorphStructName(template.Name.Value, args)

	// If already specialized, return it
	if existing, ok := c.structs[canonType(newName)]; ok {
		return existing
	}

	// Deep copy and substitute
	spec := &candy_ast.StructStatement{
		Token: template.Token,
		Name:  &candy_ast.Identifier{Token: template.Name.Token, Value: newName},
	}

	for _, f := range template.Fields {
		newF := f
		newF.TypeName = c.substituteType(f.TypeName, mapping)
		spec.Fields = append(spec.Fields, newF)
	}

	// Register the specialized struct
	c.structs[canonType(newName)] = spec
	c.bind(newName, "type:"+canonType(newName))
	c.SpecializedStructs = append(c.SpecializedStructs, spec)

	// Now check the body of the specialized struct
	c.checkStructBody(spec)

	return spec
}

func (c *Checker) substituteType(e candy_ast.Expression, mapping map[string]candy_ast.Expression) candy_ast.Expression {
	if e == nil {
		return nil
	}

	// If it's a simple identifier, check if it's a generic parameter
	if id, ok := e.(*candy_ast.Identifier); ok {
		for k, v := range mapping {
			if strings.EqualFold(k, id.Value) {
				return v
			}
		}
		return id
	}

	// If it's a type expression (e.g. List<T>), substitute recursively
	if te, ok := e.(*candy_ast.TypeExpression); ok {
		newTe := &candy_ast.TypeExpression{
			Token: te.Token,
			Name:  te.Name,
		}
		for _, arg := range te.Arguments {
			newTe.Arguments = append(newTe.Arguments, c.substituteType(arg, mapping))
		}
		return newTe
	}

	return e
}
