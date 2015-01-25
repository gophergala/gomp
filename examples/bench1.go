package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Be sure you launch this benchmark via ./run_bench.sh")
	runtime.GOMAXPROCS(runtime.NumCPU())
	const N = 100000000
	c := make([]float64, N+1)
	var h float64
	h = math.Pi / N

	//Sequential execution
	beg1 := time.Now()
	for i := 0; i <= N; i++ {
		c[i] = math.Exp(math.Sin(float64(i)*h) + math.Cos(math.Pi+float64(i)*h))
	}
	fmt.Println("Sequential execution took: ", time.Since(beg1))

	//Parallel execution
	beg2 := time.Now()
	//gomp
	for i := 0; i <= N; i++ {
		c[i] = math.Exp(math.Sin(float64(i)*h) + math.Cos(math.Pi+float64(i)*h))
	}
	end2 := time.Now()
	fmt.Println("Parallel execution took: ", end2.Sub(beg2))

}
