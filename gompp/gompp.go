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
		result, err := preproc.PreprocFile(scanner.Text(), "Stdin")
		if (err != nil) {
			fmt.Fprintln(os.Stderr, "Gompp error while using preproc.PreprocFile: ", err.Error())
			os.Exit(-1)
		} else {
			fmt.Fprintf(os.Stdin, result)
		}
		return
}
