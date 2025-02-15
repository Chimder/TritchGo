package main

import (
	"fmt"
	"slices"
)

func main() {
	s1 := []int{1, 3, 2, 10, 9, 8}
	s2 := []int{4, 5, 6}
	s3 := append(s1, s2...)
	slices.Sort(s3)

	fmt.Println(s3)
}
