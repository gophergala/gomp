package preproc

import (
	"fmt"
	"testing"
)

const (
	in00 = `package main

import "fmt"

func Foo() {
   for i := a * b + c; i < 10; i++ {
   }
   for j := 31337; j > -10; j-- {
   }
   for f0, f1 := 0, 1 ; f0 < f1; f0, f1 = f1, f0 + f1 {
   }
}

func main() {
	fmt.Println("Hello, World!")
}
`
	out00 = `package main

import "fmt"

func Foo() {
	{
		__sym0, __sym1, __sym2 := a*b+c, 10, 1
		for i := __sym0; i < __sym1; i += __sym2 {
		}
	}
	{
		__sym3, __sym4, __sym5 := 31337, -10, -1
		for j := __sym3; j > __sym4; j += __sym5 {
		}
	}
	for f0, f1 := 0, 1; f0 < f1; f0, f1 = f1, f0+f1 {
	}
}
func main() {
	fmt.Println("Hello, World!")
}
`
)

func TestPreprocFile(t *testing.T) {
	result, err := PreprocFile(in00, "in00")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(result)
	if result != out00 {
		t.Errorf("Failure")
	}
}
