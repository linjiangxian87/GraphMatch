package test

import (
	"fmt"
	"math/rand"
	"time"
)

/*
// KuhnMunkres 结构体，包含 KM 算法所需的数据
type KuhnMunkres struct {
	n        int     // 图的大小（n×n）
	graph    [][]int // 权值矩阵（邻接矩阵）
	A        []int   // 左部顶点的顶标
	B        []int   // 右部顶点的顶标
	matchU   []int   // 左部顶点的匹配结果（matchU[i] 表示左顶点i匹配的右顶点）
	matchV   []int   // 右部顶点的匹配结果
	visitedU []bool  // 左顶点是否被访问过（用于找增广路径）
	visitedV []bool  // 右顶点是否被访问过
	slack    []int   // 用于记录顶标调整时的最小差值
	prev     []int   // 记录路径的前驱节点（用于调整顶标时的交错树）
}

// NewKuhnMunkres 初始化 KM 算法结构体
func NewKuhnMunkres(graph [][]int) *KuhnMunkres {
	n := len(graph)
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

	// 初始化顶标 A[i] 为左部顶点的最大边权
	for i := 0; i < n; i++ {
		maxVal := 0
		for j := 0; j < n; j++ {
			if graph[i][j] > maxVal {
				maxVal = graph[i][j]
			}
		}
		kmm.A[i] = maxVal
	}

	// 初始化 B 为 0
	for i := 0; i < n; i++ {
		kmm.B[i] = 0
	}

	// 初始化匹配数组为 -1（未匹配）
	for i := 0; i < n; i++ {
		kmm.matchU[i] = -1
		kmm.matchV[i] = -1
	}

	return kmm
}

// 执行 KM 算法，返回最大权匹配的总权值和匹配结果
func (kmm *KuhnMunkres) FindMaxMatching() (total int, pairs [][2]int) {
	for {
		// 重置访问标记
		for i := 0; i < kmm.n; i++ {
			kmm.visitedU[i] = false
			kmm.visitedV[i] = false
		}

		// 尝试通过匈牙利算法找到增广路径

		//假设左边所有点已有匹配
		found := true
		for i := 0; i < kmm.n; i++ { //遍历所有左节点
			if kmm.matchU[i] == -1 { //如果该左节点未匹配
				found = kmm.find(i) //尝试寻找增广路径
				if !found {         //没找到，调整顶标
					break
				}
			}
		}

		if found {
			break // 所有顶点都匹配，结束
		}

		// 如果未找到增广路径，调整顶标
		kmm.adjust()
	}

	// 计算总权值和匹配对
	total = 0
	for i := 0; i < kmm.n; i++ {
		v := kmm.matchU[i]
		if v != -1 {
			total += kmm.graph[i][v]
			pairs = append(pairs, [2]int{i, v})
		}
	}

	return total, pairs
}

// find 使用递归尝试为当前顶点找到增广路径
func (kmm *KuhnMunkres) find(u int) bool {
	kmm.visitedU[u] = true       //标记该左节点已被访问
	for v := 0; v < kmm.n; v++ { //遍历所有右节点
		if !kmm.visitedV[v] && (kmm.A[u]+kmm.B[v] == kmm.graph[u][v]) { //如果该右节点未被访问，且顶标差值等于权值
			kmm.visitedV[v] = true                              //标记该右节点已被访问
			if kmm.matchV[v] == -1 || kmm.find(kmm.matchV[v]) { //如果该右节点没有被匹配，或者和该右节点匹配的左节点可以换人（找到增广路径）
				kmm.matchU[u] = v //该左节点和右节点配对
				kmm.matchV[v] = u
				return true
			}
		}
	}
	return false
}

// adjust 调整顶标，使相等子图扩大
func (kmm *KuhnMunkres) adjust() {
	// 计算交错树中的顶点
	//已访问过的左节点和右节点
	var S, T []int
	for u := 0; u < kmm.n; u++ {
		if kmm.visitedU[u] {
			S = append(S, u)
		}
	}
	for v := 0; v < kmm.n; v++ {
		if kmm.visitedV[v] {
			T = append(T, v)
		}
	}

	// 计算 d：所有不在 T 中的右顶点的最小 (A[u]+B[v] - w[u][v])
	//计算最小调整量
	minD := 1<<31 - 1 // 初始化为极大值
	for v := 0; v < kmm.n; v++ {
		if !kmm.visitedV[v] { //如果该右节点不在T中（未被访问过）
			for _, u := range S { //遍历所有已访问的左节点
				diff := kmm.A[u] + kmm.B[v] - kmm.graph[u][v]
				if diff < minD {
					minD = diff
				}
			}
		}
	}

	// 调整顶标
	for _, u := range S {
		kmm.A[u] -= minD
	}
	for _, v := range T {
		kmm.B[v] += minD
	}
}

func Testkm() {
	// 示例权值矩阵（3x3）
	graph := [][]int{
		{3, 5, 2},
		{1, 4, 6},
		{2, 0, 3},
	}

	kmm := NewKuhnMunkres(graph)
	total, pairs := kmm.FindMaxMatching()

	fmt.Println("最大权值总和:", total)
	fmt.Println("匹配对:", pairs)
}
*/
/*
type KuhnMunkres struct {
    n        int          // 左右顶点数量（比如男生和女生的数量）
    graph    [][]int      // 好感度矩阵：graph[u][v]表示男生u对女生v的好感度
    A        []int        // 男生的期望值（顶标）
    B        []int        // 女生的期望值（顶标）
    matchU   []int        // matchU[u]是男生u匹配的女生
    matchV   []int        // matchV[v]是女生v匹配的男生
    visitedU []bool       // 记录访问过的男生
    visitedV []bool       // 记录访问过的女生
    slack    []int        // 记录每个女生的最小"差距值"
    prev     []int        // 记录路径（可选，用于调试）
}

// 初始化KM算法
func NewKuhnMunkres(n int, graph [][]int) *KuhnMunkres {
    kmm := &KuhnMunkres{
        n:      n,
        graph:  graph,
        A:      make([]int, n),
        B:      make([]int, n),
        matchU: make([]int, n),
        matchV: make([]int, n),
        visitedU: make([]bool, n),
        visitedV: make([]bool, n),
        slack:   make([]int, n),
        prev:    make([]int, n),
    }
    // 初始化男生的期望值A为最大好感度
    for u := 0; u < n; u++ {
        maxA := 0
        for v := 0; v < n; v++ {
            if graph[u][v] > maxA {
                maxA = graph[u][v]
            }
        }
        kmm.A[u] = maxA
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

    // 计算总好感度
    total := 0
    for u := 0; u < kmm.n; u++ {
        if kmm.matchU[u] != -1 {
            total += kmm.graph[u][kmm.matchU[u]]
        }
    }
    return total
}

// find函数：尝试为男生u找到匹配的女生
func (kmm *KuhnMunkres) find(u int) bool {
    kmm.visitedU[u] = true // 标记男生u已访问
    for v := 0; v < kmm.n; v++ {
        if !kmm.visitedV[v] { // 女生v未被访问过
            // 计算当前边的"差距值"
            diff := kmm.A[u] + kmm.B[v] - kmm.graph[u][v]
            if diff < kmm.slack[v] { // 更新slack[v]的最小值
                kmm.slack[v] = diff
                kmm.prev[v] = u // 记录路径（可选）
            }
            if diff == 0 { // 这对情侣的期望值刚好匹配
                kmm.visitedV[v] = true // 标记女生v已访问
                // 如果女生v单身，或者她的现任对象能换人
                if kmm.matchV[v] == -1 || kmm.find(kmm.matchV[v]) {
                    kmm.matchU[u] = v // 男生u和女生v配对
                    kmm.matchV[v] = u
                    return true
                }
            }
        }
    }
    return false // 找不到匹配
}

// adjust函数：调整顶标，让匹配更容易
func (kmm *KuhnMunkres) adjust() {
    // 找到所有未被访问的女生中的最小slack值
    minD := math.MaxInt32
    for v := 0; v < kmm.n; v++ {
        if !kmm.visitedV[v] && kmm.slack[v] < minD {
            minD = kmm.slack[v]
        }
    }

    // 调整所有已访问的男生的期望值A[u] -= minD
    for u := 0; u < kmm.n; u++ {
        if kmm.visitedU[u] {
            kmm.A[u] -= minD
        }
    }
    // 调整所有已访问的女生的期望值B[v] += minD
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
*/
/*
package main

import (
    "fmt"
    "math"
)

type KuhnMunkres struct {
    n        int          // 左右顶点数量
    graph    [][]int      // 好感度矩阵（0表示不存在，正数/负数表示存在）
    A        []int        // 男生的期望值（顶标）
    B        []int        // 女生的期望值（顶标）
    matchU   []int        // matchU[u]是男生u匹配的女生
    matchV   []int        // matchV[v]是女生v匹配的男生
    visitedU []bool       // 记录访问过的男生
    visitedV []bool       // 记录访问过的女生
    slack    []int        // 记录每个女生的最小"差距值"
    prev     []int        // 记录路径（可选）
}

// 初始化KM算法
func NewKuhnMunkres(n int, graph [][]int) *KuhnMunkres {
    kmm := &KuhnMunkres{
        n:      n,
        graph:  graph,
        A:      make([]int, n),
        B:      make([]int, n),
        matchU: make([]int, n),
        matchV: make([]int, n),
        visitedU: make([]bool, n),
        visitedV: make([]bool, n),
        slack:   make([]int, n),
        prev:    make([]int, n),
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

func main() {
    // 示例输入：3个男生和3个女生
    graph := [][]int{
        {0, -5, 2}, // 男生0：女生0不存在，女生1好感度-5，女生2好感度2
        {4, 0, -3}, // 男生1：女生0好感度4，女生1不存在，女生2好感度-3
        {2, 7, 0},  // 男生2：女生0好感度2，女生1好感度7，女生2不存在
    }
    kmm := NewKuhnMunkres(3, graph)
    total := kmm.MaxWeightMatching()
    fmt.Println("最大好感度总和：", total) // 输出应为 2 + 4 +7 = 13
}
*/

/*
package KM

import (
"GraphMatching/typeStruct"
"fmt"
)

// 全局的匹配结果
var globalResult = make(map[string]string)
// 记录入参的无人机，并每次匹配结束后修改其状态
var allUavs = make(map[string]*typeStruct.Uav)
// 全局变量：未分配的任务
var globalUnassigned []typeStruct.Task

// 该程序包的入口，进行初匹配和次匹配
func GraphMatch(uavs []*typeStruct.Uav, tasks []typeStruct.Task) {
	// 初始化全局变量：全无人机
	allUavs = make(map[string]*typeStruct.Uav) // 每次调用前重置全局变量
	for _, uav := range uavs {
		allUavs[uav.Uid] = uav // 确保存储的是指针
	}
	// 初分配：非时序任务和起始时序任务（type 0或1）
	initialTasks := filterTasks(tasks, 0, 1)
	doKM(uavs, initialTasks)

	// 处理时序任务（type >=2）
	maxType := 0
	for _, task := range tasks {
		if task.TaskType > maxType {
			maxType = task.TaskType
		}
	}
	for taskType := 2; taskType <= maxType; taskType++ {
		currentTasks := filterTasks(tasks, taskType)
		doKM(uavs, currentTasks)
	}
}

// 根据任务类型划分任务
func filterTasks(tasks []typeStruct.Task, types ...int) []typeStruct.Task {
	filtered := make([]typeStruct.Task, 0)
	for _, task := range tasks {
		for _, t := range types {
			if task.TaskType == t {
				filtered = append(filtered, task)
				break
			}
		}
	}
	return filtered
}

// 调用KM入口
func doKM(uavs []*typeStruct.Uav, tasks []typeStruct.Task) {
	// 1. 使用KM算法进行初步分配
	onceKM(uavs, tasks)

	// 2. 收集未分配的任务
	var remainingTasks []typeStruct.Task
	for _, task := range tasks {
		if globalResult[task.TaskID] == "" {
			remainingTasks = append(remainingTasks, task)
		}
	}

	// 3. 使用兜底策略处理剩余任务
	ultimateMethod(uavs, remainingTasks)
}

// 最小一次KM，记得更新全局变量
func onceKM(uavs []*typeStruct.Uav, tasks []typeStruct.Task) {
	fmt.Println("---------------------一次KM开始---------------------------")
	// 构建权重矩阵
	graph := buildGraph(uavs, tasks)

	// 调用KM算法（注意：KM要求左右顶点数相同，需处理数量差异）
	n := len(tasks)
	m := len(uavs)
	// 若任务数多于无人机，需补全无人机到n个（虚拟节点权值为0）
	// 若无人机多于任务，需补全任务到m个（虚拟节点权值为0）
	// 这里假设任务数<=无人机数，否则需要调整
	// 为简单起见，这里直接使用较小的n和m
	kmm := NewKuhnMunkresZero(n, m, graph)
	total := kmm.MaxWeightMatching()

	// 处理匹配结果
	for taskIndex := 0; taskIndex < n; taskIndex++ {
		task := tasks[taskIndex]
		uavIndex := kmm.matchU[taskIndex]
		if uavIndex == -1 {
			// 未匹配
			globalUnassigned = append(globalUnassigned, task)
			continue
		}
		if uavIndex >= len(uavs) {
			// 超出无人机数量（虚拟节点）
			globalUnassigned = append(globalUnassigned, task)
			continue
		}
		// 更新匹配结果
		uav := uavs[uavIndex]
		globalResult[task.TaskID] = uav.Uid
		// 更新无人机状态
		updateUav(uav, task)
	}
}

// 构建二部图（任务作为左节点，无人机作为右节点）
func buildGraph(uavs []*typeStruct.Uav, tasks []typeStruct.Task) [][]int {
	n := len(tasks)
	m := len(uavs)
	graph := make([][]int, n)
	for i := range graph {
		graph[i] = make([]int, m)
	}
	for taskIndex, task := range tasks {
		for uavIndex, uav := range uavs {
			weight := calculateWeight(task, *uav)
			graph[taskIndex][uavIndex] = weight
		}
	}
	return graph
}

// 计算边权值（考虑资源、优先级、通信约束、负载惩罚）
func calculateWeight(task typeStruct.Task, uav typeStruct.Uav) int {
	resourceVal := calculateResourceValue(task, uav)
	if resourceVal == 0 {
		return 0
	}
	commVal := calculateCommValue(task, uav)
	if commVal == 0 {
		return 0
	}
	// 计算惩罚因子（无人机已加载任务数+1的倒数）
	punishFactor := 1.0 / float64(len(uav.LoadedTasks)+1)
	// 权值公式：资源差 * 优先级 * 通信值 * 惩罚因子，放大1000倍
	weight := int(float64(resourceVal) * float64(task.Priority) * commVal * punishFactor * 1000)
	return weight
}

// 计算资源约束值（最小剩余资源）
func calculateResourceValue(task typeStruct.Task, uav typeStruct.Uav) int {
	minDiff := 1<<31 - 1 // 最大整数
	for resType, resNeed := range task.NeedResources {
		resAvailable, ok := uav.Resources[resType]
		if !ok || resAvailable < resNeed {
			return 0
		}
		diff := resAvailable - resNeed
		if diff < minDiff {
			minDiff = diff
		}
	}
	return minDiff
}

// 计算通信约束值（依赖上游任务是否在可通信无人机上）
func calculateCommValue(task typeStruct.Task, uav typeStruct.Uav) float64 {
	if task.TaskType == 0 || task.TaskType == 1 {
		return 1.0
	}
	// 获取上游任务的分配无人机
	prevTaskID := task.Prev
	prevUavID, ok := globalResult[prevTaskID]
	if !ok {
		// 上游未分配
		return 0.0
	}
	prevUav, exists := allUavs[prevUavID]
	if !exists {
		return 0.0
	}
	// 检查当前无人机是否在上游无人机的通信列表中
	comm, exists := prevUav.NextUavs[uav.Uid]
	if !exists {
		return 0.0
	}
	return float64(comm%10) / 10.0 // 假设通信能力取个位数作为权重
}

// 更新无人机状态（减少资源、添加任务）
func updateUav(uav *typeStruct.Uav, task typeStruct.Task) {
	// 减少资源
	for resType, resNeed := range task.NeedResources {
		uav.Resources[resType] -= resNeed
	}
	// 添加任务到已装载任务列表
	uav.LoadedTasks = append(uav.LoadedTasks, task)
}

// 兜底策略：贪心分配剩余任务
func ultimateMethod(uavs []*typeStruct.Uav, tasks []typeStruct.Task) {
	for _, task := range tasks {
		maxWeight := 0
		var selectedUav *typeStruct.Uav
		for _, uav := range uavs {
			weight := calculateWeight(task, *uav)
			if weight > maxWeight {
				maxWeight = weight
				selectedUav = uav
			}
		}
		if maxWeight == 0 {
			// 无法分配，加入全局未分配列表
			globalUnassigned = append(globalUnassigned, task)
		} else {
			// 更新匹配结果和无人机状态
			globalResult[task.TaskID] = selectedUav.Uid
			updateUav(selectedUav, task)
		}
	}
}

*/

// 定义任务结构体
type Task struct {
	TaskID        string
	TaskType      int
	Priority      int
	Prev          string
	Next          string
	NeedResources map[int]int
}

// 定义无人机结构体
type Uav struct {
	Uid         string
	Resources   map[int]int // 资源id->资源数量
	NextUavs    map[string]int
	LoadedTasks []Task
}

// 随机生成任务和无人机
func generateTasksAndUavs(n int, m int) ([]*Uav, []Task) {
	rand.Seed(time.Now().UnixNano())

	// 生成无人机
	uavs := make([]*Uav, m)
	for i := 0; i < m; i++ {
		id := fmt.Sprintf("uav%d", i+1)
		resources := map[int]int{
			1: rand.Intn(21) + 80, // 资源1: 80-100
			2: rand.Intn(21) + 80, // 资源2: 80-100
		}
		nextUavs := make(map[string]int)
		// 每个无人机随机选择2-3个通信目标
		numConnections := rand.Intn(2) + 2
		for j := 0; j < numConnections; j++ {
			otherID := fmt.Sprintf("uav%d", rand.Intn(m)+1)
			if otherID != id {
				nextUavs[otherID] = rand.Intn(10) + 1 // 通信能力1-10
			}
		}
		uavs[i] = &Uav{
			Uid:       id,
			Resources: resources,
			NextUavs:  nextUavs,
		}
	}

	// 生成任务
	tasks := make([]Task, n)
	taskCountByType := map[int]int{
		0: (n * 40) / 100, // type 0: 40%
		1: (n * 30) / 100, // type 1: 30%
		2: (n * 20) / 100, // type 2: 20%
		3: (n * 10) / 100, // type 3: 10%
	}

	// 生成任务ID池（用于type≥2的prev任务）
	taskIDPool := make([]string, 0)
	currentTaskIndex := 0

	for taskType, count := range taskCountByType {
		for i := 0; i < count; i++ {
			id := fmt.Sprintf("task%d", currentTaskIndex+1)
			currentTaskIndex++

			// 生成优先级（1-10）
			priority := rand.Intn(9) + 1

			// 生成资源需求（小/中/大）
			resourceSize := ""
			switch {
			case rand.Float32() < 0.7: // 70% 小任务
				resourceSize = "small"
			case rand.Float32() < 0.9: // 20% 中任务
				resourceSize = "medium"
			default: // 10% 大任务
				resourceSize = "large"
			}

			// 生成资源需求值
			needResources := make(map[int]int)
			switch resourceSize {
			case "small":
				needResources[1] = rand.Intn(11) + 10 // 10-20
				needResources[2] = rand.Intn(11) + 10
			case "medium":
				needResources[1] = rand.Intn(11) + 30 // 30-40
				needResources[2] = rand.Intn(11) + 30
			case "large":
				needResources[1] = rand.Intn(11) + 50 // 50-60
				needResources[2] = rand.Intn(11) + 50
			}

			// 处理任务类型
			prev := ""
			if taskType >= 2 { // type≥2的任务需要设置prev
				// 从已生成的任务中随机选择一个prev
				if len(taskIDPool) > 0 {
					prevIndex := rand.Intn(len(taskIDPool))
					prev = taskIDPool[prevIndex]
				} else {
					prev = "-1" // 如果没有可用任务，则设为-1
				}
			}

			tasks[currentTaskIndex-1] = Task{
				TaskID:        id,
				TaskType:      taskType,
				Priority:      priority,
				Prev:          prev,
				NeedResources: needResources,
			}

			// 将当前任务ID加入池中（供后续type≥2任务使用）
			taskIDPool = append(taskIDPool, id)
		}
	}

	return uavs, tasks
}

// 打印生成结果（调试用）
func PrintGeneratedData(uavs []*Uav, tasks []Task) {
	fmt.Println("=== 无人机信息 ===")
	for _, uav := range uavs {
		fmt.Printf("ID: %s, 资源: %v, 通信: %v\n", uav.Uid, uav.Resources, uav.NextUavs)
	}

	fmt.Println("\n=== 任务信息 ===")
	for _, task := range tasks {
		fmt.Printf("ID: %s, 类型: %d, 优先级: %d, Prev: %s, 资源需求: %v\n",
			task.TaskID, task.TaskType, task.Priority, task.Prev, task.NeedResources)
	}
}

func main() {
	// 示例调用
	uavs, tasks := generateTasksAndUavs(10, 4) // 生成10个任务，4个无人机
	PrintGeneratedData(uavs, tasks)
}
