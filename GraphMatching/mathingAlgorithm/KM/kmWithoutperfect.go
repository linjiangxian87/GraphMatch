package KM

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

// TODO 兼容了全零行和全零列和非方阵，但是会出现可能不为最大匹配（最优）的情况
// 本KM是实现完美匹配
type KuhnMunkresZero struct {
	n        int     // 左顶点数
	m        int     // 右顶点数
	graph    [][]int // 权重矩阵，graph[u][v]存在当0<=u<n, 0<=v<m
	A        []int   // 左顶点顶标
	B        []int   // 右顶点顶标
	matchU   []int   // 左顶点匹配的右顶点（初始-1）
	matchV   []int   // 右顶点匹配的左顶点（初始-1）
	visitedU []bool  // 左顶点访问标记
	visitedV []bool  // 右顶点访问标记
	slack    []int   // 松弛值
	prev     []int   // 路径记录（未使用）
}

// 初始化KM算法
func NewKuhnMunkresZero(n, m int, graph [][]int) *KuhnMunkresZero {
	kmm := &KuhnMunkresZero{
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

	// 初始化顶标A为每个左顶点u的有效边的最大权值
	for u := 0; u < n; u++ {
		maxA := math.MinInt32
		for v := 0; v < m; v++ {
			if graph[u][v] > maxA {
				maxA = graph[u][v]
			}
		}
		// 如果该左顶点没有有效边，顶标设为0
		//if maxA == math.MinInt32 {
		//	kmm.A[u] = 0
		//} else {
		//	kmm.A[u] = maxA
		//}
		kmm.A[u] = maxA

	}

	// 初始化匹配数组为未匹配状态
	for v := 0; v < m; v++ {
		kmm.matchV[v] = -1
	}
	for u := 0; u < n; u++ {
		kmm.matchU[u] = -1
	}
	return kmm
}

// 执行KM算法，返回最大权匹配的总权值
func (kmm *KuhnMunkresZero) MaxWeightMatching() int {
	for {
		// 重置访问标记和slack数组
		for u := 0; u < kmm.n; u++ {
			kmm.visitedU[u] = false
		}
		for v := 0; v < kmm.m; v++ {
			kmm.visitedV[v] = false
			kmm.slack[v] = math.MaxInt32
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
		if found {
			// 找到增广路径，继续循环寻找
			continue
		}

		// 调整顶标
		minD := kmm.adjust()
		if minD == 0 {
			break // 无法调整顶标，终止循环
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

// 寻找增广路径
// TODO 改为非递归形式
func (kmm *KuhnMunkresZero) find(u int) bool {
	kmm.visitedU[u] = true
	for v := 0; v < kmm.m; v++ {
		// 跳过不存在的边（权值为0的情况）
		if kmm.graph[u][v] == 0 {
			continue
		}
		if !kmm.visitedV[v] {
			diff := kmm.A[u] + kmm.B[v] - kmm.graph[u][v]
			if diff < kmm.slack[v] {
				kmm.slack[v] = diff
				kmm.prev[v] = u // 记录路径（未使用）
			}
			if diff == 0 {
				kmm.visitedV[v] = true
				// 如果右顶点未匹配或能通过递归找到增广路径
				if kmm.matchV[v] == -1 || kmm.find(kmm.matchV[v]) {
					kmm.matchU[u] = v
					kmm.matchV[v] = u
					return true
				}
			}
		}
	}
	return false
}

// 调整顶标
func (kmm *KuhnMunkresZero) adjust() int {
	minD := math.MaxInt32
	// 找到未被访问的右顶点中最小的slack值
	for v := 0; v < kmm.m; v++ {
		if !kmm.visitedV[v] && kmm.slack[v] < minD {
			minD = kmm.slack[v]
		}
	}
	// 如果minD未被更新，说明无法调整
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
func Mainn1() {
	graph := [][]int{
		{0, 40000, 65000, 40000},
		{0, 40000, 65000, 40000},
		{0, 5000, 30000, 12000},
	}
	n := 3 // 左顶点数
	m := 4 // 右顶点数（即使右2无边，仍需包含）
	kmm := NewKuhnMunkresZero(n, m, graph)
	total := kmm.MaxWeightMatching()
	fmt.Println("最大权值总和：", total)

	// 输出匹配对（左→右）
	fmt.Print("匹配对（左→右）：")
	for u := 0; u < n; u++ {
		fmt.Printf("%d→%d ", u, kmm.matchU[u])
	}
	fmt.Println()

	// 输出匹配对（右→左）
	fmt.Print("匹配对（右→左）：")
	for v := 0; v < m; v++ {
		fmt.Printf("%d→%d ", v, kmm.matchV[v])
	}
}

func MainFile1() {
	// 打开输入文件
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Printf("无法打开输入文件: %v\n", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var n, m int
	// 读取n和m
	_, err = fmt.Fscan(reader, &n)
	if err != nil {
		fmt.Printf("读取n时出错: %v\n", err)
		return
	}
	_, err = fmt.Fscan(reader, &m)
	if err != nil {
		fmt.Printf("读取m时出错: %v\n", err)
		return
	}

	// 初始化好感度矩阵
	graph := make([][]int, n)
	for i := range graph {
		graph[i] = make([]int, n)
	}

	// 读取m条边
	for i := 0; i < m; i++ {
		var y, c, h int
		_, err = fmt.Fscan(reader, &y, &c, &h)
		if err != nil {
			fmt.Printf("读取第%d条边时出错: %v\n", i+1, err)
			return
		}
		// 将节点号转换为0-based索引
		graph[y-1][c-1] = h
	}

	// 执行KM算法
	kmm := NewKuhnMunkresZero(n, n, graph)
	total := kmm.MaxWeightMatching()

	// 输出结果
	fmt.Println(total)
	ans := make([]int, n)
	for i := 0; i < n; i++ {
		ans[i] = kmm.matchV[i] + 1 // 转换为1-based输出
	}
	s := fmt.Sprint(ans)
	s = strings.Trim(s, "[]")
	fmt.Println(s)
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
	kmm := NewKuhnMunkresZero(n, n, graph)
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
