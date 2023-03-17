package mypackage

func FindKthLargest(nums []int, k int) int {

	// 构建最大堆	[2n+1],[2n+2]
	length := len(nums)
	for i := length/2 - 1; i >= 0; i-- {
		compareAndSwap(nums, i, length)
	}
	//remove 前k个 通过伸缩size来非显式的的删除，然后每次重新维护一次堆性质
	size := length
	for i := 1; i < k; i++ {
		swap(nums, 0, size-1)
		//resort
		size--
		compareAndSwap(nums, 0, size)

	}
	return nums[0]
}

// 默认的肯定是左边会比右边大 构建堆的过程是 从n/2的节点开始判断每一个叶子节点最大值做交换，递归到底层
func compareAndSwap(nums []int, i int, length int) {
	largest := i
	subNode := 2*i + 1
	if subNode < length {
		if nums[subNode] > nums[largest] {
			largest = subNode
		}
	}
	subNode = 2*i + 2
	if subNode < length {
		if nums[subNode] > nums[largest] {
			largest = subNode
		}
	}
	if largest != i {
		swap(nums, i, largest)
		//在存在和右节点交换的场景,检查现在的subNode和另一个节点谁大，大的应该在左边
		compareAndSwap(nums, largest, length)
	}
}

func swap(nums []int, i int, j int) {
	temp := nums[i]
	nums[i] = nums[j]
	nums[j] = temp
}
