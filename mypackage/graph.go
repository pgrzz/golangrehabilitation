package mypackage

// https://leetcode.cn/problems/number-of-islands/
type Point struct {
	x int
	y int
}

type Queue []Point

func (q *Queue) Enqueue(value Point) {
	*q = append(*q, value)
}

func (q *Queue) Dequeue() Point {
	if len(*q) == 0 {
		panic("queue is empty")
	}
	value := (*q)[0]
	*q = (*q)[1:]
	return value
}

func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

func NumIslands(grid [][]byte) int {
	count := 0 //岛的数量

	m := len(grid)
	n := len(grid[0])
	queue := Queue{}
	visited := make([][]int, len(grid))
	for i := range grid {
		visited[i] = make([]int, len(grid[i]))
	}
	//遇到水出栈 bfs把上下左右的节点加入，如果不是水则继续走，
	//走过的地方用一个visit标记 4个指针

	// {'1', '1', '1'},
	// 	{'0', '1', '0'},
	// 	{'1', '1', '1'},
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if visited[i][j] != 1 && grid[i][j] == '1' {

				count++
				visited[i][j] = 1
				queue.Enqueue(Point{x: i, y: j})

				for len(queue) > 0 {
					p := queue.Dequeue()
					//看这个陆地有多宽多长走到头
					// 没到最上方并且没来过这块大陆并且大陆上是陆地
					for up := p.x - 1; up >= 0 && visited[up][p.y] != 1 && grid[up][p.y] == '1'; up-- {
						visited[up][p.y] = 1
						queue.Enqueue(Point{x: up, y: p.y})
					}
					for down := p.x + 1; down < m && visited[down][p.y] != 1 && grid[down][p.y] == '1'; down++ {
						visited[down][p.y] = 1
						queue.Enqueue(Point{x: down, y: p.y})

					}
					for left := p.y - 1; left >= 0 && visited[p.x][left] != 1 && grid[p.x][left] == '1'; left-- {
						visited[p.x][left] = 1
						queue.Enqueue(Point{x: p.x, y: left})
					}
					for right := p.y + 1; right < n && visited[p.x][right] != 1 && grid[p.x][right] == '1'; right++ {
						visited[p.x][right] = 1
						queue.Enqueue(Point{x: p.x, y: right})
					}

				}

			}
		}
	}

	return count
}
