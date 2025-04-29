package KM

//实现km算法
import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

// todo 已经在洛谷验证过了，虽然有部分ttl，但是我本地执行快的一比
type KuhnMunkres struct {
	n        int     // 左右顶点数量
	graph    [][]int // 好感度矩阵（0表示不存在，正数/负数表示存在）
	A        []int   // 男生的期望值（顶标）
	B        []int   // 女生的期望值（顶标）
	matchU   []int   // matchU[u]是男生u匹配的女生
	matchV   []int   // matchV[v]是女生v匹配的男生
	visitedU []bool  // 记录访问过的男生
	visitedV []bool  // 记录访问过的女生
	slack    []int   // 记录每个女生的最小"差距值"
	prev     []int   // 记录路径（可选）
}

// 初始化KM算法
func NewKuhnMunkres(n int, graph [][]int) *KuhnMunkres {
	kmm := &KuhnMunkres{
		n:        n,
		graph:    graph,
		A:        make([]int, n),
		B:        make([]int, n),
		matchU:   make([]int, n),
		matchV:   make([]int, n),
		visitedU: make([]bool, n),
		visitedV: make([]bool, n),
		slack:    make([]int, n),
		prev:     make([]int, n),
	}

	// 初始化顶标A为每个男生u的最大有效好感度（排除0）
	for u := 0; u < n; u++ {
		maxA := -math.MaxInt32 // 初始为极小值
		for v := 0; v < n; v++ {
			if graph[u][v] != 0 { // 只要边存在（非0）
				if graph[u][v] > maxA {
					maxA = graph[u][v]
				}
			}
		}
		// 如果该男生没有有效边（所有边权为0），则顶标A[u]设为0
		if maxA == -math.MaxInt32 {
			kmm.A[u] = 0
		} else {
			kmm.A[u] = maxA
		}
	}

	// 初始时所有女生未匹配
	for v := 0; v < n; v++ {
		kmm.matchV[v] = -1
	}
	for u := 0; u < n; u++ {
		kmm.matchU[u] = -1
	}
	return kmm
}

// 执行KM算法，返回最大权匹配的总好感度
func (kmm *KuhnMunkres) MaxWeightMatching() int {
	for {
		// 重置访问标记
		for u := 0; u < kmm.n; u++ {
			kmm.visitedU[u] = false
		}
		for v := 0; v < kmm.n; v++ {
			kmm.visitedV[v] = false
		}

		// 初始化slack数组为极大值
		for v := 0; v < kmm.n; v++ {
			kmm.slack[v] = math.MaxInt32
		}

		found := true
		for u := 0; u < kmm.n; u++ {
			if kmm.matchU[u] == -1 {
				found = kmm.find(u)
				if !found {
					break
				}
			}
		}
		if found {
			break // 所有男生都找到对象，结束
		}
		kmm.adjust() // 调整顶标
	}

	// 计算总好感度（忽略不存在的边）
	total := 0
	for u := 0; u < kmm.n; u++ {
		v := kmm.matchU[u]
		if v != -1 && kmm.graph[u][v] != 0 { // 边存在（权值非0）
			total += kmm.graph[u][v]
		}
	}
	return total
}

// find函数：尝试为男生u找到匹配的女生
func (kmm *KuhnMunkres) find(u int) bool {
	kmm.visitedU[u] = true
	for v := 0; v < kmm.n; v++ {
		if kmm.graph[u][v] == 0 { // 边不存在，跳过
			continue
		}
		if !kmm.visitedV[v] {
			diff := kmm.A[u] + kmm.B[v] - kmm.graph[u][v]
			if diff < kmm.slack[v] {
				kmm.slack[v] = diff
				kmm.prev[v] = u // 记录路径
			}
			if diff == 0 { // 可以匹配
				kmm.visitedV[v] = true
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

// adjust函数：调整顶标
func (kmm *KuhnMunkres) adjust() {
	// 找到未被访问的女生中的最小slack值
	minD := math.MaxInt32
	for v := 0; v < kmm.n; v++ {
		if !kmm.visitedV[v] && kmm.slack[v] < minD {
			minD = kmm.slack[v]
		}
	}

	// 调整顶标
	for u := 0; u < kmm.n; u++ {
		if kmm.visitedU[u] {
			kmm.A[u] -= minD
		}
	}
	for v := 0; v < kmm.n; v++ {
		if kmm.visitedV[v] {
			kmm.B[v] += minD
		}
	}

	// 更新未被访问的女生的slack值
	for v := 0; v < kmm.n; v++ {
		if !kmm.visitedV[v] {
			kmm.slack[v] -= minD
		}
	}
}

func Mainn() {
	// 示例输入：3个男生和3个女生

	graph := [][]int{
		{1, 0, 0}, // 男生0：女生0不存在，女生1好感度-5，女生2好感度2
		{4, 1, 3}, // 男生1：女生0好感度4，女生1不存在，女生2好感度-3
		{2, 0, 0}, // 男生2：女生0好感度2，女生1好感度7，女生2不存在
	}

	kmm := NewKuhnMunkres(3, graph)
	total := kmm.MaxWeightMatching()
	fmt.Println("最大好感度总和：", total) // 输出应为 2 + 4 +7 = 13
	fmt.Printf("匹配对：", kmm.matchU)
}
func Main() {
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
	kmm := NewKuhnMunkres(n, graph)
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

func MainFile() {
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
	kmm := NewKuhnMunkres(n, graph)
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
