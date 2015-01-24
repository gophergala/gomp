package gensym

import (
	"bytes"
	"fmt"
	"strconv"
)

func MkGen(s string) func() string {
	var count int = 0
	return func() string {
		buffer := bytes.NewBufferString("")
		fmt.Fprint(buffer, "__sym", strconv.Itoa(count))
		count++
		return buffer.String()
	}
}
