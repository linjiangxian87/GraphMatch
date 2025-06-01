package contrastExperiment

import (
	"GraphMatching/typeStruct"
	"math"
)

// DoAuction 实现拍卖算法进行多无人机任务分配。
func DoAuction(uavs []*typeStruct.Uav, tasks []typeStruct.Task) (map[string]string, []typeStruct.Task) {
	result := make(map[string]string) // 任务ID->无人机ID 的分配结果
	var unassigned []typeStruct.Task  // 未分配任务列表

	// 建立 UAV ID 到 UAV 指针的映射，便于查找前序任务所在 UAV
	uavMap := make(map[string]*typeStruct.Uav)
	for _, uav := range uavs {
		uavMap[uav.Uid] = uav
	}

	// 第一阶段：处理 type 0 和 1 的任务（无依赖或起始任务）
	var phaseTasks []typeStruct.Task
	for _, task := range tasks {
		if task.TaskType == 0 || task.TaskType == 1 {
			phaseTasks = append(phaseTasks, task)
		}
	}
	remaining := runAuction(uavs, phaseTasks, result, uavMap)
	// 将未分配任务加入结果列表
	for _, task := range remaining {
		unassigned = append(unassigned, task)
	}

	// 后续阶段：依次处理 有依赖
	// 找出任务中的最大类型
	maxType := 0
	for _, task := range tasks {
		if task.TaskType > maxType {
			maxType = task.TaskType
		}
	}
	for ttype := 2; ttype <= maxType; ttype++ {
		var tasksByType []typeStruct.Task
		for _, task := range tasks {
			if task.TaskType == ttype {
				tasksByType = append(tasksByType, task)
			}
		}
		if len(tasksByType) == 0 {
			continue
		}
		// 将任务按前序任务是否已分配分为可分配和不可分配
		var eligible []typeStruct.Task
		for _, task := range tasksByType {
			if task.Prev != "" && task.Prev != "-1" {
				// 如果前序任务已分配，则加入可竞标列表；否则视为不可分配
				if _, ok := result[task.Prev]; ok {
					eligible = append(eligible, task)
				} else {
					unassigned = append(unassigned, task)
				}
			} else {
				eligible = append(eligible, task)
			}
		}
		// 对当前阶段可分配任务执行拍卖
		remaining = runAuction(uavs, eligible, result, uavMap)
		for _, task := range remaining {
			unassigned = append(unassigned, task)
		}
	}

	return result, unassigned
}

// 每轮拍卖，做三件事：
// 1.每台 UAV对所有剩余任务出价，选一个“最适合自己的任务”；
// 2.每个任务挑选对它出价最高的无人机；
// 3.把任务分配给它选的 UAV，更新无人机资源与任务列表。
func runAuction(uavs []*typeStruct.Uav, tasks []typeStruct.Task, result map[string]string, uavMap map[string]*typeStruct.Uav) []typeStruct.Task {
	// 剩余未分配任务
	remaining := make([]typeStruct.Task, len(tasks))
	copy(remaining, tasks)
	if len(remaining) == 0 {
		return remaining
	}

	for {
		// 记录本轮每个任务的最高出价和投标 UAV 索引
		type bidInfo struct {
			bid    float64
			uavIdx int
		}
		bids := make(map[int]bidInfo) //键是任务在 remaining 切片里的索引 tIdx，值是该任务目前最高的 bidInfo

		// 每架 UAV 选出自己能竞标的最佳任务及出价
		for uIdx, uav := range uavs {
			bestIdx := -1       // 记录这台 UAV 想竞标的任务在 remaining 里的下标，-1 表示还没选
			var bestBid float64 // 记录它对该任务的最高出价
			for tIdx, task := range remaining {

				// 资源匹配度评估：计算该 UAV 完成该任务后最小资源余量
				minDiff := math.MaxInt32
				feasible := true // 标记是否满足资源需求
				for resType, need := range task.NeedResources {
					avail, ok := uav.Resources[resType]
					if !ok || avail < need { // 无人机上没有这个资源类型，或资源不足
						feasible = false
						break
					}
					diff := avail - need
					if diff < minDiff {
						minDiff = diff
					}
				}
				if !feasible { // 如果这台 UAV 资源不满足，就跳过这个任务
					continue
				}
				resourceVal := minDiff + 1

				// 通信约束评估
				commVal := 1.0
				if task.TaskType >= 2 {
					prevID := task.Prev                 //前驱任务id
					if prevID != "" && prevID != "-1" { //存在前驱任务
						prevUavID := result[prevID] //前驱任务所在UAV
						if prevUavID == uav.Uid {   //前驱任务由同一 UAV完成
							// 如果前序任务由同一 UAV 完成，则通信得分较高
							commVal = 10.0
						} else {
							// 检查前序 UAV 与当前 UAV 是否可通信
							if prevUav, ok := uavMap[prevUavID]; ok {
								if val, exist := prevUav.NextUavs[uav.Uid]; exist {
									commVal = float64(val%10) / 10.0
								} else {
									// 不可通信，跳过此任务
									continue
								}
							} else {
								continue
							}
						}
					}
				}

				// 无人机负载惩罚因子，负载越高惩罚越大
				punish := 1.0 / float64(len(uav.LoadedTasks)+1)

				// 计算综合出价
				bid := float64(resourceVal*task.Priority) * commVal * punish
				if bid > bestBid {
					bestBid = bid
					bestIdx = tIdx
				}
			}

			// 记录这台 UAV 对任务 bestIdx 的出价，如果比已有出价高就覆盖
			if bestIdx >= 0 {
				if cur, ok := bids[bestIdx]; !ok || bestBid > cur.bid {
					bids[bestIdx] = bidInfo{bid: bestBid, uavIdx: uIdx}
				}
			}

		}

		// 如果没有任何 UAV 出价，就跳出主循环，拍卖结束
		if len(bids) == 0 {
			break
		}

		// 执行本轮分配：出价最高的 UAV 获得对应任务
		assignedIndices := make([]int, 0)
		for tIdx, info := range bids {
			if info.bid <= 0 {
				continue
			}
			task := remaining[tIdx]
			uav := uavs[info.uavIdx]
			// 更新分配映射
			result[task.TaskID] = uav.Uid
			// 更新 UAV 资源与已装载任务列表
			for resType, need := range task.NeedResources {
				uav.Resources[resType] -= need
			}
			uav.LoadedTasks = append(uav.LoadedTasks, task)
			assignedIndices = append(assignedIndices, tIdx)
		}
		if len(assignedIndices) == 0 {
			break
		}
		// 从剩余任务中移除已分配的任务
		newRemaining := make([]typeStruct.Task, 0, len(remaining))
		for idx, task := range remaining {
			skip := false
			for _, aIdx := range assignedIndices {
				if idx == aIdx {
					skip = true
					break
				}
			}
			if !skip {
				newRemaining = append(newRemaining, task)
			}
		}
		remaining = newRemaining
	}

	return remaining
}
