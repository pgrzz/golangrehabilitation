package mypackage

// https://leetcode.cn/problems/merge-sorted-array/description/
func Merge(nums1 []int, m int, nums2 []int, n int) {

	if m < 1 {
		for n > 0 {
			n--
			nums1[n] = nums2[n]
		}
		return
	}
	if n < 1 {
		return
	}

	i := len(nums1) - n - 1
	j := len(nums2) - 1
	t := len(nums1) - 1
	for t > 0 {
		if i >= 0 && nums1[i] > nums2[j] {
			nums1[t] = nums1[i]
			i--
		} else if j >= 0 && nums1[i] <= nums2[j] {

			nums1[t] = nums2[j]
			j--
		}
		t--

		for t >= 0 && i >= 0 && j < 0 {
			nums1[t] = nums1[i]
			t--
			i--
		}

		for t >= 0 && j >= 0 && i < 0 {
			nums1[t] = nums2[j]
			t--
			j--
		}
	}
}

func merge(intervals [][]int) [][]int {
	i := 1
	j := 0
	t := 0
	result := [][]int{}

	length := len(intervals)
	for i < length {
		//check
		next := intervals[i]
		pre := intervals[j]
		if pre[1] > next[0] {
			// 有交集
			result = append(result, []int{pre[0], next[1]})

			pre[0] = result[t][0]
			pre[1] = result[t][1]

			i++
			t++
		} else {
			//无交集
			result = append(result, pre)
			result = append(result, next)

			i++
			j++
			t++
			t++
		}

	}
	return result

}
