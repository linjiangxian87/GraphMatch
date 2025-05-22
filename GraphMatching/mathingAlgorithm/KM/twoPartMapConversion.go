package KM

import (
	"GraphMatching/typeStruct"
	"fmt"
	"math"
	"sort"
)

// 全局的匹配结果
var globalResult = make(map[string]string)

// 记录入参的无人机，并每次匹配结束后修改其状态
var allUavs = make(map[string]*typeStruct.Uav)

// 全局变量：未分配的任务
var globalUnassigned []typeStruct.Task

// 该程序包的入口，进行初匹配和次匹配
func GraphMatch(uavs []*typeStruct.Uav, tasks []typeStruct.Task) (map[string]string, []typeStruct.Task) {
	//TODO 每次调用该函数是否需要重置全局变量？
	//allUavs = make(map[string]*typeStruct.Uav) // 每次调用前重置全局变量

	fmt.Println("---------------------初分配开始---------------------------")
	// 初始化全局变量：全无人机
	for _, uav := range uavs {
		allUavs[uav.Uid] = uav
	}

	// 初分配：非时序任务和起始时序任务（type 0或1）
	initialTasks := filterTasks(tasks, 0, 1)
	doKM(uavs, initialTasks)
	//fmt.Println("初分配结果如下：", globalResult)
	//fmt.Println("初分配成功匹配数：", len(globalResult))
	//fmt.Println("匹配结果：", globalResult)
	//fmt.Println("初分配中未分配任务：", globalUnassigned)
	//fmt.Println("初分配后无人机状态：")
	//printUavs(uavs)
	fmt.Println("---------------------初分配结束---------------------------")
	fmt.Println()

	fmt.Println("---------------------次分配开始---------------------------")
	// 处理时序任务（type >=2）
	maxType := 0
	for _, task := range tasks {
		if task.TaskType > maxType {
			maxType = task.TaskType
		}
	}

	// 次分配：处理时序任务（type >=2）
	for taskType := 2; taskType <= maxType; taskType++ {
		currentTasks := filterTasks(tasks, taskType)
		doKM(uavs, currentTasks)
	}
	fmt.Println("---------------------次分配结束---------------------------")

	return globalResult, globalUnassigned
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
// TODO 完善该函数
func doKM(uavs []*typeStruct.Uav, tasks []typeStruct.Task) {
	//将待匹配任务排序：优先级越高越前，相同则需求越高越前
	sort.Slice(tasks, func(i, j int) bool {

		if tasks[i].Priority == tasks[j].Priority {
			a := taskResourceSum(tasks[i])
			b := taskResourceSum(tasks[j])
			return a > b
		}

		return tasks[i].Priority > tasks[j].Priority
	})
	remainingTasks := make([]typeStruct.Task, 0)
	//每次取m个任务分配，将每次未分配出去的任务记录
	//TODO 校对一下
	m := len(uavs)
	now := 0
	for now < len(tasks) {
		badTasks := onceKM(uavs, tasks[now:min(now+m, len(tasks))])
		now += m
		//TODO 这能将一个数组塞到一个数组里头嘛？
		remainingTasks = append(remainingTasks, badTasks...)
	}

	//调用兜底策略，将这些未分配的任务再尝试分配一次，仍有任务未被分配则记录进全局的未分配表中，等待下次调度
	fmt.Println("贪心前未分配任务：", globalUnassigned)
	ultimateMethod(uavs, remainingTasks)
	fmt.Println("贪心后未分配任务：", globalUnassigned)
}

// 最小一次KM
func onceKM(uavs []*typeStruct.Uav, tasks []typeStruct.Task) []typeStruct.Task {
	fmt.Println("---------------------一次KM开始---------------------------")

	// 构建权重矩阵
	graph := buildGraph(uavs, tasks)
	fmt.Println("构建权重矩阵：", graph)

	// 调用KM算法
	kmm := NewKuhnMunkresZero(len(uavs), len(tasks), graph)
	kmm.MaxWeightMatching()

	remainingTasks := make([]typeStruct.Task, 0)
	//第j个任务匹配到第i个无人机上
	fmt.Println("匹配结果：", kmm.matchV)
	for j, i := range kmm.matchV {
		if i == -1 {
			remainingTasks = append(remainingTasks, tasks[j])
		} else {
			updateUav(uavs[i], tasks[j])
			globalResult[tasks[j].TaskID] = uavs[i].Uid
		}
	}
	fmt.Println("---------------------一次KM结束---------------------------")
	return remainingTasks
}

// 构建二部图
func buildGraph(uavs []*typeStruct.Uav, tasks []typeStruct.Task) [][]int {
	n := len(tasks)
	m := len(uavs)
	graph := make([][]int, m)
	for i := range graph {
		graph[i] = make([]int, n)
	}

	for i, uav := range uavs {
		for j, task := range tasks {
			graph[i][j] = calculateWeight(task, *uav)
		}
	}
	return graph
}

// 计算二部图边权值(暂定为整数)
func calculateWeight(task typeStruct.Task, uav typeStruct.Uav) int {
	//资源约束值
	resourceVal := calculateResourceValue(task, uav)
	if resourceVal == 0 {
		return 0
	}

	//任务优先级
	priority := task.Priority

	//通信约束值
	commVal := calculateCommValue(task, uav)
	if commVal == 0 {
		return 0
	}

	//无人机负载惩罚因子
	//punishFactor := 1.0 / float64((len(uav.LoadedTasks) + 1))
	punishFactor := 1.0 / (float64((len(uav.LoadedTasks) + 1))) * 100

	//权值计算公式：四者相乘，并转为int
	// 转换为整数权重（放大1000倍防止浮点损失）
	weight := int(float64(resourceVal) * float64(priority) * commVal * punishFactor)
	return weight
}

// 计算资源约束值
func calculateResourceValue(task typeStruct.Task, uav typeStruct.Uav) int {
	minDiff := math.MaxInt32
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
	return minDiff + 1 //+1防止无人机资源和任务所需资源相同时导致0
}

// 计算通信约束值
func calculateCommValue(task typeStruct.Task, uav typeStruct.Uav) float64 {
	if task.TaskType == 0 || task.TaskType == 1 {
		return 1.0
	}

	// 获取上游任务所在的无人机
	prevTaskID := task.Prev
	prevUavID, ok := globalResult[prevTaskID]
	if !ok {
		return 0.0 // 上游未分配，无法通信
	}

	// 如果是自身的上游任务，则直接返回10
	if prevUavID == uav.Uid {
		return 10.0
	}

	prevUav := allUavs[prevUavID]
	if prevUav == nil {
		return 0.0
	}

	// 检查当前无人机是否在上游无人机的通信列表中
	comm, exists := prevUav.NextUavs[uav.Uid]
	if !exists {
		return 0.0
	}

	return float64(comm)
}

func updateUav(uav *typeStruct.Uav, task typeStruct.Task) {
	// 减少资源
	for resType, resNeed := range task.NeedResources {
		uav.Resources[resType] -= resNeed
	}

	// 添加任务到已装载任务列表
	uav.LoadedTasks = append(uav.LoadedTasks, task)
}
func uavResourceSum(uav typeStruct.Uav) int {
	sum := 0
	for _, v := range uav.Resources {
		sum += v
	}
	return sum
}
func taskResourceSum(task typeStruct.Task) int {
	sum := 0
	for _, v := range task.NeedResources {
		sum += v
	}
	return sum
}

func ultimateMethod(uavs []*typeStruct.Uav, tasks []typeStruct.Task) {
	for _, task := range tasks {
		maxWeight := 0
		var loadingUav *typeStruct.Uav
		for _, uav := range uavs {
			weight := calculateWeight(task, *uav)
			if weight > maxWeight {
				maxWeight = weight
				loadingUav = uav
			}
		}

		if maxWeight == 0 {
			//没救了，记录到全局表中
			globalUnassigned = append(globalUnassigned, task)
		} else {
			globalResult[task.TaskID] = loadingUav.Uid
			updateUav(loadingUav, task)
		}
	}
}

func printUav(uav typeStruct.Uav) {
	fmt.Printf("UID：%s\t", uav.Uid)
	fmt.Printf("\t资源：")
	for i, _ := range uav.Resources {
		fmt.Printf("(%d:%d)", i, uav.Resources[i])
	}
	fmt.Printf("\t下一跳：")
	for i, _ := range uav.NextUavs {
		fmt.Printf("(%s:%d)", i, uav.NextUavs[i])
	}
	fmt.Printf("\t装载任务：")
	for i, _ := range uav.LoadedTasks {
		fmt.Printf("%s\t", uav.LoadedTasks[i].TaskID)
	}
	fmt.Println()
}
func printUavs(uavs []*typeStruct.Uav) {
	for uav, _ := range uavs {
		printUav(*uavs[uav])
	}
}
