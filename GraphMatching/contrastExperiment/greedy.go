package contrastExperiment

import (
	"GraphMatching/typeStruct"
	"fmt"
	"math"
	"sort"
)

// 全局的匹配结果
var globalResult1 = make(map[string]string)

// 记录入参的无人机，并每次匹配结束后修改其状态
var allUavs1 = make(map[string]*typeStruct.Uav)

// 全局变量：未分配的任务
var globalUnassigned1 []typeStruct.Task

func Dogreedy(uavs []*typeStruct.Uav, tasks []typeStruct.Task) (map[string]string, []typeStruct.Task) {
	// 初始化全局变量
	//globalResult = make(map[string]string)
	//globalUnassigned = make([]typeStruct.Task, 0)
	//allUavs = make(map[string]*typeStruct.Uav)

	fmt.Println("---------------------开始执行贪心算法---------------------------")
	// 初始化全局变量：全无人机
	for _, uav := range uavs {
		allUavs1[uav.Uid] = uav
	}
	initialTasks := filterTasks(tasks, 0, 1)
	greedyAlgorithm(uavs, initialTasks)

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
		greedyAlgorithm(uavs, currentTasks)
	}
	fmt.Println("---------------------贪心算法执行结束---------------------------")
	return globalResult1, globalUnassigned1
}

// 贪心算法实现
func greedyAlgorithm(uavs []*typeStruct.Uav, tasks []typeStruct.Task) (map[string]string, []typeStruct.Task) {

	// 按优先级降序排序任务
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})

	for _, task := range tasks {
		maxWeight := 0
		bestUavIndex := -1

		// 遍历所有无人机，选择权值最高的
		for i, uav := range uavs {
			/*
				// 检查资源约束
				resourceCheck := true
				for resType, resNeed := range task.NeedResources {
					if uav.Resources[resType] < resNeed {
						resourceCheck = false
						break
					}
				}
				if !resourceCheck {
					continue
				}

				// 检查通信约束
				commCheck := true
				if task.TaskType >= 2 {
					prevUavID, ok := globalResult[task.Prev]
					if !ok {
						commCheck = false
					} else if prevUavID != uav.Uid {
						prevUav := allUavs[prevUavID]
						if prevUav == nil || prevUav.NextUavs[uav.Uid] == 0 {
							commCheck = false
						}
					}
					if !commCheck {
						continue
					}
				}
			*/
			// 计算权值
			weight := calculateWeight(task, *uav)
			if weight > maxWeight {
				maxWeight = weight
				bestUavIndex = i
			}
		}

		// 分配任务
		if bestUavIndex != -1 {
			uav := uavs[bestUavIndex]
			globalResult1[task.TaskID] = uav.Uid
			/*
				// 更新无人机资源
				for resType, resNeed := range task.NeedResources {
					uav.Resources[resType] -= resNeed
				}
				uav.LoadedTasks = append(uav.LoadedTasks, task)
			*/
			updateUav(uav, task)
		} else {
			globalUnassigned1 = append(globalUnassigned1, task)
		}
	}

	return globalResult1, globalUnassigned1
}

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
	punishFactor := 1.0 / float64(len(uav.LoadedTasks)+1)

	//权值计算公式：四者相乘，并转为int
	// 转换为整数权重（放大1000倍防止浮点损失）
	weight := int(float64(resourceVal) * float64(priority) * commVal * punishFactor * 1000)
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
	prevUavID, ok := globalResult1[prevTaskID]
	if !ok {
		return 0.0 // 上游未分配，无法通信
	}

	// 如果是自身的上游任务，则直接返回10
	if prevUavID == uav.Uid {
		return 10.0
	}

	prevUav := allUavs1[prevUavID]
	if prevUav == nil {
		return 0.0
	}

	// 检查当前无人机是否在上游无人机的通信列表中
	comm, exists := prevUav.NextUavs[uav.Uid]
	if !exists {
		return 0.0
	}

	return float64(comm%10) / 10.0
}

func updateUav(uav *typeStruct.Uav, task typeStruct.Task) {
	// 减少资源
	for resType, resNeed := range task.NeedResources {
		uav.Resources[resType] -= resNeed
	}

	// 添加任务到已装载任务列表
	uav.LoadedTasks = append(uav.LoadedTasks, task)
}
