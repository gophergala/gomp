package preproc

import "testing"

const (
	file00 = `
package main

import "fmt"

func Foo() {
   fmt.Println("Hello, World!")
}

func main() {
	fmt.Println("Hello, World!")
}
`
)

func TestPreprocFile(t *testing.T) {
	_, err := PreprocFileImpl(file00, "file00")
	if err != nil {
		t.Error(err.Error())
	}
}
