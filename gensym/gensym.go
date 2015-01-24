package gensym

import (
	"bytes"
	"fmt"
	"strconv"
)

func MkGen(s string) func() string {
	var count int = 0
	return func() string {
		buffer := bytes.NewBufferString("__sym")
		fmt.Fprint(buffer, strconv.Itoa(count))
		count++
		return buffer.String()
	}
}
