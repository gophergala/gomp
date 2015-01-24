package gensym_test

import (
	. "github.com/gophergala/gomp/gensym"
	"testing"
)

func TestMkGen(t *testing.T) {
	check := func(expected, actual string) {
		if actual != expected {
			t.Errorf("Fail: %s != %s\n", expected, actual)
		}
	}

	f := MkGen("")
	if VocabOn {
		for i := 0; i < 108; i++ {
			f()
		}
	}
	check("gompsym0", f())
	check("gompsym1", f())
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
