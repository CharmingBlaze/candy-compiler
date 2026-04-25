package candy_report

import (
	"fmt"
	"strings"
)

type Level int

const (
	Error Level = iota
	Warning
	Note
)

func (l Level) String() string {
	switch l {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Note:
		return "note"
	default:
		return "info"
	}
}

type Diagnostic struct {
	Level   Level
	Message string
	File    string
	Line    int
	Col     int
	Offset  int
	Length  int
}

func Report(src string, diagnostics []Diagnostic) {
	lines := strings.Split(src, "\n")

	for _, d := range diagnostics {
		levelStr := d.Level.String()
		// Basic color codes (Red for error, Yellow for warning, Cyan for note)
		// Only if not on windows or if we want to be safe, but let's assume modern terminal.
		// Windows terminal supports ANSI now.
		color := ""
		reset := "\033[0m"
		switch d.Level {
		case Error:
			color = "\033[1;31m" // Bold Red
		case Warning:
			color = "\033[1;33m" // Bold Yellow
		case Note:
			color = "\033[1;36m" // Bold Cyan
		}

		fmt.Printf("%s%s: %s%s\n", color, levelStr, d.Message, reset)
		if d.File != "" {
			fmt.Printf("  --> %s:%d:%d\n", d.File, d.Line, d.Col)
		} else {
			fmt.Printf("  --> %d:%d\n", d.Line, d.Col)
		}

		if d.Line > 0 && d.Line <= len(lines) {
			lineIdx := d.Line - 1
			lineText := lines[lineIdx]

			// Print the line with line number
			fmt.Printf("%5d | %s\n", d.Line, lineText)

			// Print the caret
			padding := strings.Repeat(" ", d.Col-1)
			caretLen := d.Length
			if caretLen <= 0 {
				caretLen = 1
			}
			carets := strings.Repeat("^", caretLen)
			fmt.Printf("      | %s%s%s%s\n", padding, color, carets, reset)
		}
		fmt.Println()
	}
}
