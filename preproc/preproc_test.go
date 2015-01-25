package preproc

import "testing"

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
		gompsym7 := (gompsym1-gompsym0)/(gompsym3*gompsym2) + 1
		gompsym6 := make(chan struct {
		}, gompsym7)
		for gompsym4 := 0; gompsym0+gompsym4*(gompsym3*gompsym2) <= gompsym1; gompsym4++ {
			go func(gompsym4 int) {
				for i, gompsym5 := gompsym0+gompsym4*(gompsym3*gompsym2), 0; i <= gompsym1 && gompsym5 < gompsym3; i, gompsym5 = i+gompsym2, gompsym5+1 {
					fmt.Println(i)
				}
				gompsym6 <- struct {
				}{}
			}(int(gompsym4))
		}
		for gompsym8 := 0; gompsym8 < gompsym7; gompsym8++ {
			<-gompsym6
		}
	}
	{
		gompsym9, gompsym10, gompsym11 := 31337, -10, 1
		for j := gompsym9; j > gompsym10; j -= gompsym11 {
		}
	}
	for f0, f1 := 0, 1; f0 < f1; f0, f1 = f1, f0+f1 {
	}
}
func Bar() {
	if true {
		{
			gompsym12, gompsym13, gompsym14 := 99, -10, 1
			for i := gompsym12; i >= gompsym13; i -= gompsym14 {
			}
		}
	}
}
func Baz() {
	g := func() {
		for i := gompsym15; i < gompsym16; i += gompsym17 {
		}
	}
	g()
}
func main() {
	fmt.Println("Hello, World!")
}
`
	in01 = `package main

import "fmt"

func main() {

	p := fmt.Println

	for i := 0; i < 134; i += 123 {
		p(10 - i)
	}

}
`
	out01 = `package main

import "fmt"

func main() {
	p := fmt.Println
	{
		gompsym0, gompsym1, gompsym2 := 0, 134, 123
		for i := gompsym0; i < gompsym1; i += gompsym2 {
			p(10 - i)
		}
	}
}
`
)

func TestPreprocFile(t *testing.T) {
	check := func(input, output, name string) {
		t.Log("Running on", name, "...")
		result, err := PreprocFile(input, name)
		if err != nil {
			t.Error(err.Error())
		}
		if result != output {
			t.Errorf("Output mismatch:\n%s\n", result)
		}
	}
	check(in00, out00, "test00")
	check(in01, out01, "test01")
}
