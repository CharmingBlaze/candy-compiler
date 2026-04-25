package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"candy/candy_evaluator"
	"candy/candy_lexer"
	"candy/candy_parser"

	"github.com/chzyer/readline"
)

// runREPL is a simple read-eval-print loop with history and basic shortcuts.
func runREPL(stdin io.Reader, stdout, stderr io.Writer) int {
	home, _ := os.UserHomeDir()
	historyFile := filepath.Join(home, ".candy_history")

	rc := io.NopCloser(stdin)
	if stdin == nil {
		rc = os.Stdin
	}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "candy> ",
		HistoryFile: historyFile,
		Stdin:       rc,
		Stdout:      stdout,
		Stderr:      stderr,
	})
	if err != nil {
		fmt.Fprintln(stderr, "failed to initialize readline:", err)
		return 1
	}
	defer rl.Close()

	env := candy_evaluator.ReplEnv()
	fmt.Fprintln(stdout, "Candy REPL. Type :help, :vars, or :exit. Ctrl+D to end.")

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF (Ctrl+D) or Interrupt (Ctrl+C)
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, ":") {
			switch {
			case line == ":exit" || line == ":q" || line == ":quit":
				fmt.Fprintln(stdout, "Bye.")
				return 0
			case line == ":help" || line == ":h" || line == "?":
				fmt.Fprintln(stdout, "  :help, :h   — this text")
				fmt.Fprintln(stdout, "  :vars       — list names bound in this session")
				fmt.Fprintln(stdout, "  :exit, :q   — leave the REPL")
			case line == ":vars" || line == ":v":
				names := env.AllNameBindings()
				sort.Strings(names)
				if len(names) == 0 {
					fmt.Fprintln(stdout, "(no bindings yet)")
				} else {
					fmt.Fprintln(stdout, strings.Join(names, " "))
				}
			default:
				fmt.Fprintln(stderr, "unknown command (use :help). :exit to quit.")
			}
			continue
		}

		l := candy_lexer.New(line + "\n")
		p := candy_parser.New(l)
		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, d := range p.Errors() {
				fmt.Fprintln(stderr, d.Message)
			}
			continue
		}
		v, err := candy_evaluator.Eval(prog, env)
		if err != nil {
			fmt.Fprintln(stderr, err)
			continue
		}
		if v != nil && v.Kind != candy_evaluator.ValNull {
			fmt.Fprintln(stdout, v.String())
		}
	}
	fmt.Fprintln(stdout, "Bye.")
	return 0
}
