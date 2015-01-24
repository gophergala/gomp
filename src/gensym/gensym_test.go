package gensym

import "testing"

func TestMkGen(t *testing.T) {
	check := func(expected, actual string) {
		if actual != expected {
			t.Errorf("Fail: %s != %s\n", expected, actual)
		}
	}

	{
		gen := MkGen("")
		check("__sym0", gen())
		check("__sym1", gen())
		check("__sym2", gen())
	}

	{
		gen := MkGen("")
		check("__sym0", gen())
		check("__sym1", gen())
	}
}
