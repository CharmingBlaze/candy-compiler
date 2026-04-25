package main

import (
	"fmt"
	"os"

	"candy/candy_pkg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: candypm <init|lock>")
		return
	}
	switch os.Args[1] {
	case "init":
		_ = os.WriteFile("candy.pkg", []byte("name = app\nversion = 0.1.0\n"), 0o644)
		fmt.Println("initialized candy.pkg")
	case "lock":
		m, err := candy_pkg.LoadManifest("candy.pkg")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := candy_pkg.WriteLock("candy.lock", m); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("wrote candy.lock")
	default:
		fmt.Println("unknown command")
	}
}
