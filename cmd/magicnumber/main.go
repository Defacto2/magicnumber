package main

import (
	"fmt"
	"os"

	"github.com/Defacto2/magicnumber"
)

func main() {
	const minArgs = 2
	if len(os.Args) < minArgs {
		fmt.Fprintln(os.Stderr, "Usage: magicnumber <path-to-file>")
		os.Exit(1)
	}

	name := os.Args[1]
	file, err := os.Open(name) //nolint:gosec
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	result := magicnumber.Find(file)
	fmt.Fprintf(os.Stdout, "%s : %s %q\n", name, result, result.Title())
}
