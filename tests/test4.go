package main

import "fmt"

func main() {

	p := fmt.Println

	for i := 100; i >= 0; i-- {
		p(10 - i)
	}

}
