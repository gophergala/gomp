package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gophergala/gomp/preproc"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // Read one line
	result, err := preproc.PreprocFile(scanner.Text(), "stdin")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Gompp error while using preproc.PreprocFile:\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdin, result)
	return
}
