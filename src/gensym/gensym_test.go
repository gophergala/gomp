package gensym

import "testing"

func check(expected, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Fail: %s != %s\n", expected, actual)
	}
}

func TestMkGen(t *testing.T) {
	gen := MkGen("")
	check("__sym0", gen(), t)
	check("__sym1", gen(), t)

}
