package contrastExperiment

import (
	"GraphMatching/typeStruct"
	"fmt"
	"math"
)

var cishu = 0 //统计现在到哪了
var originalUavs []*typeStruct.Uav

// 全局的匹配结果
var globalResult2 = make(map[string]string)

// 记录入参的无人机，并每次匹配结束后修改其状态
var allUavs2 = make(map[string]*typeStruct.Uav)

// 全局变量：未分配的任务
var globalUnassigned2 []typeStruct.Task
var resUavs []*typeStruct.Uav
var preClb = math.MaxFloat64
var preAssignedRate = 0.0

// 复制路径 map
func copyMap(m map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
func copyUavs(oldUavs []*typeStruct.Uav) []*typeStruct.Uav {
	copiedUavs := make([]*typeStruct.Uav, len(oldUavs))
	for i, uav := range oldUavs {
		copiedUavs[i] = &typeStruct.Uav{
			Uid:         uav.Uid,
			Resources:   make(map[int]int),
			NextUavs:    make(map[string]int),
			LoadedTasks: make([]typeStruct.Task, len(uav.LoadedTasks)),
		}
		for k, v := range uav.Resources {
			copiedUavs[i].Resources[k] = v
		}
		for k, v := range uav.NextUavs {
			copiedUavs[i].NextUavs[k] = v
		}
		copy(copiedUavs[i].LoadedTasks, uav.LoadedTasks)
	}
	return copiedUavs
}

// 评估组合的负载均衡指标（资源利用率标准差）
func evaluateCombination(tasks []typeStruct.Task, path map[string]string, currentUavs []*typeStruct.Uav) {
	//负载均衡度，方差表示
	clb := 0.0
	originalUavPool := make(map[string]*typeStruct.Uav)
	for _, uav := range originalUavs {
		originalUavPool[uav.Uid] = uav
	}

	utils := make(map[int]float64, len(currentUavs))
	for i, newuav := range currentUavs {
		utilx := 0.0
		for resType, resSum := range newuav.Resources {
			utilx += float64(originalUavPool[newuav.Uid].Resources[resType]-resSum) / float64(originalUavPool[newuav.Uid].Resources[resType])
		}
		utilx /= float64(len(newuav.Resources))

		utilx *= math.Abs(float64(len(newuav.LoadedTasks)) - float64(len(tasks)/len(originalUavs)))

		utils[i] = utilx
	}

	utilAvg := 0.0
	for i, _ := range utils {
		utilAvg += utils[i]
	}
	utilAvg /= float64(len(utils))

	for i, _ := range utils {
		clb += math.Pow(utils[i]-utilAvg, 2)
	}
	clb /= float64(len(utils))

	//任务完成率
	assignedRate := float64(len(path)) / float64(len(tasks))
	if assignedRate < preAssignedRate {
		return
	} else if assignedRate > preAssignedRate {
		//fmt.Println("新的任务完成率：", assignedRate)

		preClb = clb
		preAssignedRate = assignedRate

		globalResult2 = copyMap(path)
		resUavs = copyUavs(currentUavs)
	} else {
		if clb < preClb {
			//fmt.Println("新的负载均衡度：", clb)

			preClb = clb
			preAssignedRate = assignedRate

			globalResult2 = copyMap(path)
			resUavs = copyUavs(currentUavs)
		}
	}
}

func dfss(index int, path map[string]string, uavs []*typeStruct.Uav, tasks []typeStruct.Task) {
	if index == len(tasks) {
		//fmt.Println("还活着，第", cishu, "次")
		//cishu++
		evaluateCombination(tasks, path, uavs)
		return
	}

	currentTask := tasks[index]
	for i := -1; i < len(uavs); i++ {
		//该任务放空
		if i == -1 {
			dfss(index+1, path, uavs, tasks)
			continue
		}

		// 检查资源约束
		resourceCheck := true
		for resType, resNeed := range currentTask.NeedResources {
			if uavs[i].Resources[resType] < resNeed {
				resourceCheck = false
				break
			}
		}
		if !resourceCheck {
			continue
		}

		// 检查通信约束
		commCheck := true
		if currentTask.TaskType >= 2 {
			prevUavID, ok := path[currentTask.Prev]
			//找不到上游
			if !ok {
				commCheck = false
			} else if prevUavID != uavs[i].Uid { //上游无人机非本无人机
				prevUav := allUavs2[prevUavID]
				if prevUav == nil || prevUav.NextUavs[uavs[i].Uid] == 0 {
					commCheck = false
				}
			}
			if !commCheck {
				continue
			}
		}

		//满足条件，尝试装载
		for resType, resNeed := range currentTask.NeedResources {
			uavs[i].Resources[resType] -= resNeed
		}
		path[currentTask.TaskID] = uavs[i].Uid
		uavs[i].LoadedTasks = append(uavs[i].LoadedTasks, currentTask)

		dfss(index+1, path, uavs, tasks)

		for resType, resNeed := range currentTask.NeedResources {
			uavs[i].Resources[resType] += resNeed
		}
		delete(path, currentTask.TaskID)
		uavs[i].LoadedTasks = uavs[i].LoadedTasks[:len(uavs[i].LoadedTasks)-1]
	}

}

func SimpleCombination(uavs []*typeStruct.Uav, tasks []typeStruct.Task) (map[string]string, []typeStruct.Task) {
	fmt.Println("---------------------开始执行组合算法---------------------------")
	//
	//cishu = 0

	originalUavs = uavs

	usedUavs := copyUavs(uavs)
	// 初始化全局变量：全无人机
	for _, uav := range usedUavs {
		allUavs2[uav.Uid] = uav
	}
	path := make(map[string]string)
	dfss(0, path, usedUavs, tasks)
	fmt.Println(globalResult2)
	for _, task := range tasks {
		//如果不在globalResult2内，则存入globalUnassigned2
		_, ok := globalResult2[task.TaskID]
		if !ok {
			globalUnassigned2 = append(globalUnassigned2, task)
		}
	}

	reslUavPool := make(map[string]*typeStruct.Uav)
	for _, uav := range resUavs {
		reslUavPool[uav.Uid] = uav
	}
	for _, uav := range uavs {
		resUav := reslUavPool[uav.Uid]
		uav.Resources = resUav.Resources
		uav.LoadedTasks = resUav.LoadedTasks
		uav.NextUavs = resUav.NextUavs
	}
	fmt.Println("---------------------组合算法执行结束---------------------------")
	return globalResult2, globalUnassigned2
}

/*
package contrastExperiment

import (
	"GraphMatching/typeStruct"
	"math"
)

// 全局的匹配结果
var globalResult2 = make(map[string]string)

// 记录入参的无人机，并每次匹配结束后修改其状态
var allUavs2 = make(map[string]*typeStruct.Uav)

// 全局变量：未分配的任务
var globalUnassigned2 []typeStruct.Task

var preSurplus = math.MaxInt32
var prejicha = math.MaxInt32

var totalNeed = make(map[int]int)

// 评估组合的负载均衡指标（资源利用率标准差）
func evaluateCombination(uavs []*typeStruct.Uav, path map[string]string) bool {
	// 剩余资源数越少越好
	surplus := make(map[int]int)
	for _, uav := range uavs {
		for resType, res := range uav.Resources {
			surplus[resType] += res
		}
	}

	totalSurplus := 0
	for _, res := range surplus {
		totalSurplus += res
	}

	// 负载均衡极差越小越好
	maxlen := 0
	minlen := 100000
	for _, uav := range uavs {
		maxlen = max(len(uav.LoadedTasks), maxlen)
		minlen = min(len(uav.LoadedTasks), minlen)
	}
	jicha := maxlen - minlen
	jicha = 0 - jicha
	jicha *= 100

	if totalSurplus < preSurplus {
		prejicha = jicha
		return true
	}
	if totalSurplus > preSurplus {
		return false
	}

	if jicha < prejicha {
		prejicha = jicha
		return true
	}
	return false
}

// 深拷贝无人机对象
func deepCopyUav(uav *typeStruct.Uav) *typeStruct.Uav {
	newUav := &typeStruct.Uav{
		Uid:         uav.Uid,
		Resources:   make(map[int]int),
		NextUavs:    make(map[string]int),
		LoadedTasks: make([]typeStruct.Task, len(uav.LoadedTasks)),
	}
	for k, v := range uav.Resources {
		newUav.Resources[k] = v
	}
	for k, v := range uav.NextUavs {
		newUav.NextUavs[k] = v
	}
	copy(newUav.LoadedTasks, uav.LoadedTasks)
	return newUav
}

// 复制路径 map
func copyMap(m map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// DFS 递归函数
func dfss(index int, path map[string]string, uavs []*typeStruct.Uav, tasks []typeStruct.Task) {
	if index == len(tasks) {
		if evaluateCombination(uavs, path) {
			// 复制当前路径作为全局结果
			globalResult2 = copyMap(path)
		}
		return
	}

	currentTask := tasks[index]
	for i := -1; i < len(uavs); i++ {
		// 该任务放空（未分配）
		if i == -1 {
			path[currentTask.TaskID] = "" // 显式标记未分配
			dfss(index+1, path, uavs, tasks)
			continue
		}

		// 检查资源约束
		resourceCheck := true
		for resType, resNeed := range currentTask.NeedResources {
			if uavs[i].Resources[resType] < resNeed {
				resourceCheck = false
				break
			}
		}
		if !resourceCheck {
			continue
		}

		// 检查通信约束
		commCheck := true
		if currentTask.TaskType >= 2 {
			prevUavID, ok := path[currentTask.Prev]
			if !ok {
				commCheck = false
			} else if prevUavID != uavs[i].Uid {
				prevUav := allUavs2[prevUavID]
				if prevUav == nil || prevUav.NextUavs[uavs[i].Uid] == 0 {
					commCheck = false
				}
			}
			if !commCheck {
				continue
			}
		}

		// 满足条件，尝试装载
		for resType, resNeed := range currentTask.NeedResources {
			uavs[i].Resources[resType] -= resNeed
		}
		path[currentTask.TaskID] = uavs[i].Uid
		uavs[i].LoadedTasks = append(uavs[i].LoadedTasks, currentTask)

		// 递归
		dfss(index+1, path, uavs, tasks)

		// 回溯
		for resType, resNeed := range currentTask.NeedResources {
			uavs[i].Resources[resType] += resNeed
		}
		delete(path, currentTask.TaskID)
		uavs[i].LoadedTasks = uavs[i].LoadedTasks[:len(uavs[i].LoadedTasks)-1]
	}
}

// 主函数
func SimpleCombination(uavs []*typeStruct.Uav, tasks []typeStruct.Task) (map[string]string, []typeStruct.Task) {
	// 重置全局变量
	preSurplus = math.MaxInt32
	prejicha = math.MaxInt32
	globalResult2 = make(map[string]string)
	globalUnassigned2 = make([]typeStruct.Task, 0)

	// 深拷贝无人机数组
	uavsCopy := make([]*typeStruct.Uav, len(uavs))
	for i := range uavs {
		uavsCopy[i] = deepCopyUav(uavs[i])
	}

	// 初始化路径
	initialPath := make(map[string]string)

	// 执行 DFS
	dfss(0, initialPath, uavsCopy, tasks)

	// 收集未分配任务
	for _, task := range tasks {
		uavid, ok := globalResult2[task.TaskID]
		if !ok || uavid == "" {
			globalUnassigned2 = append(globalUnassigned2, task)
		}
		//} else {
		//	for _, uav := range uavs {
		//		if uav.Uid == uavid {
		//			updateUav(uav, task)
		//		}
		//	}
		//}
	}

	return globalResult2, globalUnassigned2
}
*/
