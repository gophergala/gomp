package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func main() {
	p := fmt.Println
	p("Be sure you launch this benchmark via ./run_bench.sh")
	runtime.GOMAXPROCS(4)
	const N = 100000000
	c := make([]float64, N+1)
	var h float64
	h = math.Pi / float64(N)

	//Sequential exeqution
	beg1 := time.Now()
	for i := 0; i <= N; i++ {
		c[i] = math.Exp(math.Sin(float64(i)*h) + math.Cos(math.Pi+float64(i)*h))
	}
	end1 := time.Now()
	p("Sequential execution took: ", end1.Sub(beg1))

	beg2 := time.Now()
	//gomp
	for i := 0; i <= N; i++ {
		c[i] = math.Exp(math.Sin(float64(i)*h) + math.Cos(math.Pi+float64(i)*h))
	}
	end2 := time.Now()
	p("Parallel execution took: ", end2.Sub(beg2))

}
