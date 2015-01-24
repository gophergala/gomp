package gensym_test

import (
	. "github.com/gophergala/gomp/gensym"
	"testing"
)

func check(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Errorf("Fail: %s != %s\n", expected, actual)
	}
}

func TestMkGen(t *testing.T) {
	f := MkGen("")
	if VocabOn {
		for i := 0; i < 108; i++ {
			f()
		}
	}
	check(t, "gompsym0", f())
	check(t, "gompsym1", f())
}

func TestSeveralFuncs(t *testing.T) {
	const n = 500
	const nfuncs = 10
	fs := make([](func() string), nfuncs)
	for i := 0; i < nfuncs; i++ {
		fs[i] = MkGen("")
	}
	out := make([][]string, nfuncs)
	for i := 0; i < nfuncs; i++ {
		if out[i] == nil {
			out[i] = make([]string, n)
			for j := 0; j < n; j++ {
				out[i][j] = fs[i]()
			}
		}
		for j := 0; j < n; j++ {
			if out[i][j] != out[0][j] {
				t.Fatal("Two invocations of MkGen gave different results")
			}
		}
	}
}

func TestTokens(t *testing.T) {
	var src = `package main

import (
	"fmt"
)

func main() {
	var gompsym0 int32 = 5
	gompsym1 := "asd"
	fmt.Printf("Hello, %d, %s and gompsym2\n", gompsym0, gompsym1)
}
`

	f := MkGen(src)
	if VocabOn {
		for i := 0; i < 108; i++ {
			f()
		}
	}
	check(t, "gompsym2", f())
}
