package main

import (
	"fmt"
	"golangrehabilitation/mypackage"
)

func main() {

	// nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
	nums1 := []int{4, 0, 0, 0, 0, 0}
	m := 1
	nums2 := []int{1, 2, 3, 5, 6}
	n := 5
	mypackage.Merge(nums1, m, nums2, n)
	fmt.Println(nums1)

}
