//go:build ignore

package main

import "fmt"

func main() {
	var n, a, b int
	fmt.Scan(&n)
	for range n {
		fmt.Scan(&a, &b)
		fmt.Println(a + b)
	}
}
