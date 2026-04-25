package main

import (
	"fmt"
	"os"

	"candy/candy_lsp"
)

func main() {
	if err := candy_lsp.Run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
