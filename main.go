package main

import (
	"golangrehabilitation/mypackage"
)

func main() {

	// matrix := [][]int{
	// 	{5, 1, 9, 11},
	// 	{2, 4, 8, 10},
	// 	{13, 3, 6, 7},
	// 	{15, 14, 12, 16},
	// }

	// matrix := [][]int{
	// 	{1, 2, 3},
	// 	{4, 5, 6},
	// 	{7, 8, 9},
	// }

	arrays := []int{3, 2, 3, 1, 2, 4, 5, 5, 6}

	mypackage.FindKthLargest(arrays, 4)

}
