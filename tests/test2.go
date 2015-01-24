package main

import "fmt"

func main() {

	p := fmt.Println

	for i := 10; i > 0; i-- {
		p(10 - i)
	}

}
