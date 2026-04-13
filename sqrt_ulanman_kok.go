package main

import (
	"fmt"
)

func Sqrt(x int) int {
	z := 1
	for i := 0; i <= 10; i++ {
		z -= (z*z - x) / (2 * z)
		fmt.Printf("Iteration %d: ",  z)
	}
}

func main() {
	var x int
	fmt.Scan(&x)
	Sqrt(int(x))
}
