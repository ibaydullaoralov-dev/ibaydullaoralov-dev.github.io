package main

import "fmt"

func sum(n int) int {
	if n <= 0 {
		return 0
	}
	return n + sum(n-1)
}
func main() {
	fmt.Print("n=")
	var n int
	fmt.Scan(&n)
	fmt.Println(sum(n))
}
