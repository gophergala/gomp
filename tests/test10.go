package main

import "fmt"

func main() {

	p := fmt.Println

	for i := 5; i <= 10; i += 3 {
		p(10 - i)
	}

}
