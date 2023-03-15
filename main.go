package main

import (
	"fmt"
	"golangrehabilitation/mypackage"
)

func main() {

	// matrix := [][]int{
	// 	{5, 1, 9, 11},
	// 	{2, 4, 8, 10},
	// 	{13, 3, 6, 7},
	// 	{15, 14, 12, 16},
	// }

	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	mypackage.Rotate(matrix)

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			fmt.Println(matrix[i][j])
		}
	}
}
