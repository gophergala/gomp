package preproc

import (
	"fmt"
	"testing"
)

const (
	in00 = `package main
func Foo() {
   //gomp
   for i := 0; i <= 10; i++ {
       fmt.Println(i)
   }
}
`
	out00 = `package main

import "runtime"

func Foo() {
	{
		gompsym0, gompsym1, gompsym2 := 0, 10, 1
		gompsym3 := (gompsym1 - gompsym0 + gompsym2*runtime.NumCPU()) / (gompsym2 * runtime.NumCPU())
		gompsym5 := (gompsym1-gompsym0)/(gompsym3*gompsym2) + 1
		gompsym6 := make(chan struct {
		}, gompsym5)
		for gompsym4 := 0; gompsym0+gompsym4*(gompsym3*gompsym2) <= gompsym1; gompsym4++ {
			go func(gompsym4 int) {
				for i, gompsym7 := gompsym0+gompsym4*(gompsym3*gompsym2), 0; i <= gompsym1 && gompsym7 < gompsym3; i, gompsym7 = i+gompsym2, gompsym7+1 {
					fmt.Println(i)
				}
				gompsym6 <- struct {
				}{}
			}(int(gompsym4))
		}
		for gompsym8 := 0; gompsym8 < gompsym5; gompsym8++ {
			<-gompsym6
		}
	}
}
`
	in01 = `package main

import "fmt"

func main() {

	p := fmt.Println

    //gomp
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
			fmt.Println(result)
			t.Error("Failure, output mismatch!")
		}
	}
	check(in00, out00, "test00")
	check(in01, out01, "test01")
}
