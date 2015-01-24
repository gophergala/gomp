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
		gompsym0, gompsym1, gompsym2 := a*b+c, 10, 1
		for i := gompsym0; i < gompsym1; i += gompsym2 {
		}
	}
	{
		gompsym3, gompsym4, gompsym5 := 31337, -10, -1
		for j := gompsym3; j > gompsym4; j += gompsym5 {
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
	result, err := PreprocFileImpl(in00, "in00")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(result)
	if result != out00 {
		t.Errorf("Failure")
	}
}
