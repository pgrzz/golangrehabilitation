package mypackage

const threshold = 16

func SortArray(nums []int) []int {
	return quickSort(nums)
}
func quickSort(nums []int) []int {

	begin := 0
	end := len(nums) - 1
	partition(nums, begin, end)
	return nums
}
func partition(nums []int, begin int, end int) {
	if begin >= end {
		return
	}
	if end-begin+1 <= threshold {
		insertSort(nums, begin, end)
		return
	}

	lt := begin // 小于等于区间的右边界
	i := begin + 1
	gt := end // 大于等于区间的左边界
	pivot := nums[begin]

	for i <= gt {
		if nums[i] < pivot { //
			swap(nums, i, lt)
			i++
			lt++
		} else if nums[i] > pivot { //在pivot 右边的
			swap(nums, i, gt)
			gt--
		} else { //相等的
			i++
		}
	}
	partition(nums, begin, lt-1)
	partition(nums, gt+1, end)
}

func insertSort(nums []int, begin int, end int) {
	for i := begin; i <= end; i++ {
		key := nums[i]
		j := i - 1
		for j >= begin && key < nums[j] {
			nums[j+1] = nums[j]
			j--
		}
		nums[j+1] = key
	}
}

// func swap(nums []int, i int, j int) {
// 	temp := nums[i]
// 	nums[i] = nums[j]
// 	nums[j] = temp
// }
