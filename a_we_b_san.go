package main

import "fmt"

func main() {
	var a, b int

	fmt.Print("a snay giriz=")
	fmt.Scan(&a)

	fmt.Print("b sany giriz=")
	fmt.Scan(&b)
	a = a - b
	b = a + b
	a = b - a
	fmt.Println("a=", a, "b=", b)
}
