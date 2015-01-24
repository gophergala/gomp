package main

import "fmt"

func main() {

	p := fmt.Println
	var a [10]int
	for i := 0; i < 10; i++ {
		a[i] = 1
	}

	for i := 0; i < 10; i += a[i] {
		p(10 - i)
	}

}
