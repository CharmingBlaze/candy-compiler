package candy_llvm

import "strings"

func (c *Compiler) functionAttributes(name string) string {
	lower := strings.ToLower(name)
	attrs := []string{"nounwind", "willreturn"}

	if c.isHotFunctionName(lower) {
		attrs = append(attrs, "hot", "inlinehint")
	}
	if c.isLikelyPureMathName(lower) {
		attrs = append(attrs, "readnone", "nosync", "speculatable")
	}
	return strings.Join(attrs, " ")
}

func (c *Compiler) isHotFunctionName(name string) bool {
	return name == "main" ||
		strings.Contains(name, "update") ||
		strings.Contains(name, "render") ||
		strings.Contains(name, "tick") ||
		strings.Contains(name, "step")
}

func (c *Compiler) isLikelyPureMathName(name string) bool {
	return strings.HasPrefix(name, "math_") ||
		name == "sqrt" ||
		strings.Contains(name, "dot") ||
		strings.Contains(name, "length") ||
		strings.Contains(name, "normalize") ||
		strings.Contains(name, "abs") ||
		strings.Contains(name, "min") ||
		strings.Contains(name, "max")
}
