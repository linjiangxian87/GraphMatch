package contrastExperiment

/*
import (
	"GraphMatching/typeStruct"
)

// 拍卖算法实现
func AuctionAlgorithm(uavs []*typeStruct.Uav, tasks []typeStruct.Task) (map[string]string, []typeStruct.Task) {
	globalResult := make(map[string]string)
	globalUnassigned := make([]Task, 0)

	// 初始化任务价格
	taskPrices := make(map[string]int)
	for _, task := range tasks {
		taskPrices[task.TaskID] = 0
	}

	// 初始化无人机可用资源
	uavAvailable := make([]*Uav, len(uavs))
	for i := range uavs {
		uavAvailable[i] = copyUav(uavs[i])
	}

	// 迭代拍卖过程
	for {
		assignedTasks := 0
		for _, task := range tasks {
			if _, exists := globalResult[task.TaskID]; exists {
				continue
			}

			maxBid := 0
			bestUavIndex := -1

			// 无人机竞标
			for i, uav := range uavAvailable {
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

				// 计算出价（权值 + 价格）
				bid := calculateWeight(task, *uav) + taskPrices[task.TaskID]
				if bid > maxBid {
					maxBid = bid
					bestUavIndex = i
				}
			}

			// 分配任务
			if bestUavIndex != -1 {
				uav := uavAvailable[bestUavIndex]
				globalResult[task.TaskID] = uav.Uid

				// 更新无人机资源
				for resType, resNeed := range task.NeedResources {
					uav.Resources[resType] -= resNeed
				}
				uav.LoadedTasks = append(uav.LoadedTasks, task)
				assignedTasks++
			} else {
				globalUnassigned = append(globalUnassigned, task)
			}
		}

		// 如果无法分配更多任务，终止
		if assignedTasks == 0 {
			break
		}
	}

	return globalResult, globalUnassigned
}

// 复制无人机对象
func copyUav(uav *Uav) *Uav {
	newUav := &Uav{
		Uid:         uav.Uid,
		Resources:   make(map[int]int),
		NextUavs:    make(map[string]int),
		LoadedTasks: make([]Task, len(uav.LoadedTasks)),
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
*/
