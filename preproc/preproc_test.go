package preproc

import (
	"fmt"
	"testing"
)

const (
	in00 = `package main

import "runtime"
import (
  "fmt"
  "math"
)

func Foo() {
   fmt.Println("Hello")
   for i := 0; i <= 10; i++ {
       fmt.Println(i)
   }
   for j := 31337; j > -10; j-- {
   }
   for f0, f1 := 0, 1 ; f0 < f1; f0, f1 = f1, f0 + f1 {
   }
}
func Bar() {
  if true {
    for i := 99; i >= -10; i-- {
    }
  }
}

func Baz() {
	g := func() {
		for i := 0; i < 100; i++ {
		}
	}
	g()
}

func main() {
	fmt.Println("Hello, World!")
}
`
	out00 = `package main

import "runtime"
import (
	"fmt"
	"math"
)

func Foo() {
	fmt.Println("Hello")
	{
		gompsym0, gompsym1, gompsym2 := 0, 10, 1
		gompsym3 := (gompsym1 - gompsym0 + 1) / (gompsym2 * runtime.NumCPU())
		for i := gompsym0; i <= gompsym1; i += gompsym2 {
			fmt.Println(i)
		}
	}
	{
		gompsym4, gompsym5, gompsym6 := 31337, -10, -1
		for j := gompsym4; j > gompsym5; j += gompsym6 {
		}
	}
	for f0, f1 := 0, 1; f0 < f1; f0, f1 = f1, f0+f1 {
	}
}
func Bar() {
	if true {
		{
			gompsym7, gompsym8, gompsym9 := 99, -10, -1
			for i := gompsym7; i >= gompsym8; i += gompsym9 {
			}
		}
	}
}
func Baz() {
	g := func() {
		for i := gompsym10; i < gompsym11; i += gompsym12 {
		}
	}
	g()
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
