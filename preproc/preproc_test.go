package preproc

import "testing"

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
		for i := a*b + c; i < 10; i++ {
		}
	}
	{
		for j := 31337; j > -10; j-- {
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
	if result != out00 {
		t.Errorf("Output differ:\nExpected:\n%v\nActual:\n%v\n",
			out00, result)
	}
}
