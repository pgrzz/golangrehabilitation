package mypackage

// "fmt"
// 旋转矩阵
func Rotate(matrix [][]int) {
	//核心点在于它是nxn的

	var a, b, c, d, apoint, bpoint, cpoint, dpoint int
	//层数
	level := 0
	length := len(matrix)
	for level < length {
		//更新边界指针
		apoint = level
		bpoint = level
		cpoint = length - level - 1
		dpoint = length - level - 1
		outSide := length - level - 1

		for i := level; i < outSide; i++ {
			a = matrix[level][apoint]
			b = matrix[bpoint][outSide]
			c = matrix[outSide][cpoint]
			d = matrix[dpoint][level]

			matrix[level][apoint] = d
			matrix[bpoint][outSide] = a
			matrix[outSide][cpoint] = b
			matrix[dpoint][level] = c

			apoint++
			bpoint++
			cpoint--
			dpoint--

		}
		level++

	}

}
