package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gophergala/gomp/preproc"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result, err := preproc.PreprocFile(scanner.Text(), "stdin")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Gompp error while using preproc.PreprocFile:\n", err.Error())
			os.Exit(-1)
		}
		fmt.Fprintln(os.Stdout, result)
	}
	return
}
