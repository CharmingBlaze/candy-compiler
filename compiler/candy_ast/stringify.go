package candy_ast

import (
	"fmt"
)

// StringExpr formats an expression for debugging.
func StringExpr(e Expression) string {
	if e == nil {
		return "<nil>"
	}
	if s, ok := e.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("[%T: %s]", e, e.TokenLiteral())
}

// StringStmt formats a statement.
func StringStmt(s Statement) string {
	if s == nil {
		return "<nil>"
	}
	if s2, ok := s.(interface{ String() string }); ok {
		return s2.String()
	}
	return fmt.Sprintf("[%T: %s]", s, s.TokenLiteral())
}

// IntString returns decimal string for an int.
func IntString(v int64) string { return fmt.Sprintf("%d", v) }

// FloatString returns a float string.
func FloatString(f float64) string { return fmt.Sprintf("%g", f) }
