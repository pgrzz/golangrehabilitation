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

// 1、判断 i中2, i+1中的1 是否有交集
// https://leetcode.cn/problems/merge-intervals/
func Merge22(intervals [][]int) [][]int {
	n := len(intervals)

	if n <= 1 {
		return intervals
	}

	result := [][]int{}

	i := 0
	j := 1

	//sort一把 直接插入排序好了
	for t := 1; t < len(intervals); t++ {
		tempArray := intervals[t]
		q := t - 1
		for q >= 0 && tempArray[0] < intervals[q][0] {
			intervals[q+1] = intervals[q]
			q--
		}
		intervals[q+1] = tempArray
	}

	for i < n && j < n {
		array1 := intervals[i]
		array2 := intervals[j]
		//有交集则做合并，无交集时则加入到结果集合
		//三种相交情况都需要做判断
		//case1、 右交 array1[1]>array2[0] case 2、全包裹  array1[0]<array2[0] && array1[1]<array2[1]
		// case3 左交 array1[0]>array2[1]
		//同时镜像case
		if (array2[0] <= array1[1] && array2[1] >= array1[0]) ||
			(array1[0] <= array2[1] && array1[1] >= array2[0]) {
			//2、缩小数组长度
			//copy 就好了 i不动，j+1
			maxa1 := Max(array1[1], array2[1])
			array1[1] = maxa1

			mina1 := Min(array1[0], array2[0])
			array1[0] = mina1

			//当 i+1 ==lens 如果合并完说明没有继续的了直接apend
			if j == n-1 {
				result = append(result, array1)
			}

			//让j向后移动，说明i已经合并了
			j++
		} else {
			if j == n-1 {
				//说明后面两个都需要合并
				result = append(result, array1)
				result = append(result, array2)
			} else {
				result = append(result, array1)
			}
			//i的对应数组 加入到结果中

			i = j
			j++
		}

	}
	return result

}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
