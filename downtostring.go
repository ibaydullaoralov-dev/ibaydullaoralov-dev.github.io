package main

import "fmt"

func downto(s string) string {

	if len(s) <= 1 {
		return s
	}
	return downto(s[1:]) + string(s[0])
}

func main() {
	n := "Spiderman"
	fmt.Println(downto(n))
}
