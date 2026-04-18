package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/sum", sum)
	http.ListenAndServe(":8000", nil)
}
func sum1(w http.ResponseWriter, n int) int {
	if n <= 0 {
		return 0
	}
	return n + sum1(w, n-1)
}
func sum(w http.ResponseWriter, r *http.Request) {
	x := r.FormValue("san")
	n, _ := strconv.Atoi(x)
	fmt.Fprintln(w, sum1(w, n))
}
