package test

import (
	"fmt"
	"math"
	"strings"
)

type KuhnMunkres struct {
	n        int
	m        int
	graph    [][]int
	A        []int
	B        []int
	matchU   []int
	matchV   []int
	visitedU []bool
	visitedV []bool
	slack    []int
	prev     []int // 记录路径前驱
}

func NewKuhnMunkres(n, m int, graph [][]int) *KuhnMunkres {
	kmm := &KuhnMunkres{
		n:        n,
		m:        m,
		graph:    graph,
		A:        make([]int, n),
		B:        make([]int, m),
		matchU:   make([]int, n),
		matchV:   make([]int, m),
		visitedU: make([]bool, n),
		visitedV: make([]bool, m),
		slack:    make([]int, m),
		prev:     make([]int, m),
	}

	// 初始化顶标A
	for u := 0; u < n; u++ {
		maxA := -1
		for v := 0; v < m; v++ {
			if graph[u][v] > maxA {
				maxA = graph[u][v]
			}
		}
		kmm.A[u] = maxA
	}

	// 初始化匹配数组
	for v := 0; v < m; v++ {
		kmm.matchV[v] = -1
	}
	for u := 0; u < n; u++ {
		kmm.matchU[u] = -1
	}

	return kmm
}

func (kmm *KuhnMunkres) MaxWeightMatching() int {
	for {
		// 重置标记和slack数组
		for u := 0; u < kmm.n; u++ {
			kmm.visitedU[u] = false
		}
		for v := 0; v < kmm.m; v++ {
			kmm.visitedV[v] = false
			kmm.slack[v] = math.MaxInt32
			kmm.prev[v] = -1 // 初始化前驱为-1
		}

		found := false
		// 尝试为每个未匹配的左顶点寻找增广路径
		for u := 0; u < kmm.n; u++ {
			if kmm.matchU[u] == -1 {
				if kmm.find(u) {
					found = true
				}
			}
		}
		if !found {
			// 调整顶标
			minD := kmm.adjust()
			if minD == 0 {
				break
			}
		}
	}

	// 计算总权值
	total := 0
	for u := 0; u < kmm.n; u++ {
		v := kmm.matchU[u]
		if v != -1 {
			total += kmm.graph[u][v]
		}
	}
	return total
}

// 迭代版的find函数
func (kmm *KuhnMunkres) find(startU int) bool {
	stack := make([]struct{ u, vIndex int }, 0)
	stack = append(stack, struct{ u, vIndex int }{startU, 0})

	// 标记初始节点为已访问
	kmm.visitedU[startU] = true

	for len(stack) > 0 {
		elem := stack[len(stack)-1]
		currentU := elem.u
		vIndex := elem.vIndex
		stack = stack[:len(stack)-1]

		// 遍历所有v，从当前vIndex开始
		found := false
		for v := vIndex; v < kmm.m; v++ {
			if kmm.graph[currentU][v] == 0 {
				continue
			}
			if !kmm.visitedV[v] {
				diff := kmm.A[currentU] + kmm.B[v] - kmm.graph[currentU][v]
				if diff < kmm.slack[v] {
					kmm.slack[v] = diff
					kmm.prev[v] = currentU
				}
				if diff == 0 {
					kmm.visitedV[v] = true
					if kmm.matchV[v] == -1 {
						// 找到增广路径的终点，沿路径更新匹配
						kmm.updateMatching(v)
						return true
					} else {
						nextU := kmm.matchV[v]
						// 将当前节点的处理进度保存，并压入下一个节点
						stack = append(stack, struct{ u, vIndex int }{currentU, v + 1})
						// 标记下一个节点为未访问
						kmm.visitedU[nextU] = true
						stack = append(stack, struct{ u, vIndex int }{nextU, 0})
						found = true
						break
					}
				}
			}
		}
		if !found {
			// 没有找到有效边，回溯
			stack = append(stack, struct{ u, vIndex int }{currentU, vIndex})
		}
	}
	return false
}

// 沿路径更新匹配
func (kmm *KuhnMunkres) updateMatching(v int) {
	for v != -1 {
		u := kmm.prev[v]
		prevV := kmm.matchU[u]
		kmm.matchU[u] = v
		kmm.matchV[v] = u
		v = prevV
	}
}

// 调整顶标
func (kmm *KuhnMunkres) adjust() int {
	minD := math.MaxInt32
	// 找到未被访问的右顶点中最小的slack值
	for v := 0; v < kmm.m; v++ {
		if !kmm.visitedV[v] && kmm.slack[v] < minD {
			minD = kmm.slack[v]
		}
	}
	if minD == math.MaxInt32 {
		return 0
	}

	// 调整顶标
	for u := 0; u < kmm.n; u++ {
		if kmm.visitedU[u] {
			kmm.A[u] -= minD
		}
	}
	for v := 0; v < kmm.m; v++ {
		if kmm.visitedV[v] {
			kmm.B[v] += minD
		} else {
			kmm.slack[v] -= minD
		}
	}
	return minD
}

// 测试函数
func Main() {
	graph := [][]int{
		{1, 2, 0},
		{4, 1, 0},
		{2, 7, 0},
	}
	n := 3
	m := 3
	kmm := NewKuhnMunkres(n, m, graph)
	total := kmm.MaxWeightMatching()
	fmt.Println("最大权值总和：", total) // 应输出11

	fmt.Print("匹配对（左→右）：")
	for u := 0; u < n; u++ {
		fmt.Printf("%d→%d ", u, kmm.matchU[u])
	}
	fmt.Println()

	fmt.Print("匹配对（右→左）：")
	for v := 0; v < m; v++ {
		fmt.Printf("%d→%d ", v, kmm.matchV[v])
	}
}

func Main1() {
	var n int
	var m int
	fmt.Scan(&n)
	fmt.Scan(&m)
	//fmt.Println(n, m)
	graph := make([][]int, n)
	for i := range graph {
		graph[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		var a int
		var b int
		var weight int
		fmt.Scan(&a)
		fmt.Scan(&b)
		fmt.Scan(&weight)
		//fmt.Println(a, b, weight)
		graph[a-1][b-1] = weight
	}
	//fmt.Println(graph)
	kmm := NewKuhnMunkres(n, n, graph)
	total := kmm.MaxWeightMatching()

	//fmt.Println("最大权值总和:", total)
	//fmt.Println("匹配对:", pairs)

	fmt.Println(total)
	ans := make([]int, n)
	for i := 0; i < n; i++ {
		ans[i] = kmm.matchV[i] + 1
	}
	s := fmt.Sprint(ans)      // 输出 "[5 4 1 3 2]"
	s = strings.Trim(s, "[]") // 去掉方括号
	fmt.Println(s)            // 输出 "5 4 1 3 2"

}

/*
remainingTasks := make([]typeStruct.Task, len(tasks))
	copy(remainingTasks, tasks)
	sort.Slice(remainingTasks, func(i, j int) bool {
		return remainingTasks[i].Priority > remainingTasks[j].Priority ||
			(taskResourceSum(remainingTasks[i]) > taskResourceSum(remainingTasks[j]))
	})
	m := len(uavs)
	for len(remainingTasks) > 0 {
		batchSize := m
		if len(remainingTasks) < batchSize {
			batchSize = len(remainingTasks)
		}
		currentTasks := remainingTasks[:batchSize]

		// 调用一次KM处理当前批次
		onceKM(uavs, currentTasks)

		// 更新剩余任务列表（移除已处理的任务）
		remainingTasks = remainingTasks[batchSize:]
	}
*/
