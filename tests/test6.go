package main

import "fmt"

func main() {

	p := fmt.Println

	for i := 22650; i > 2134; i -= 123 {
		p(10 - i)
	}

}
