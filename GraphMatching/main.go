package main

import (
	"GraphMatching/mathingAlgorithm/KM"
	"GraphMatching/typeStruct"
	"fmt"
	"time"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func newUav(id string, nex map[string]int, res map[int]int) *typeStruct.Uav {
	uav1 := &typeStruct.Uav{
		Uid:         id,
		Resources:   res,
		NextUavs:    nex,
		LoadedTasks: []typeStruct.Task{}, // 初始无任务
	}
	return uav1
}
func newTask(id string, prev string, need map[int]int, priority int, ttype int) typeStruct.Task {
	task := typeStruct.Task{
		TaskID:        id,
		TaskType:      ttype,
		Priority:      priority,
		Prev:          prev,
		NeedResources: need,
	}
	return task
}

func test1() {
	uavs := make([]*typeStruct.Uav, 0)
	tasks := make([]typeStruct.Task, 0)

	uavs = append(uavs, newUav("uav1", map[string]int{"uav3": 12, "uav4": 11}, map[int]int{1: 15, 2: 25}))
	uavs = append(uavs, newUav("uav2", map[string]int{"uav3": 13}, map[int]int{1: 24, 2: 16}))
	uavs = append(uavs, newUav("uav3", map[string]int{"uav1": 12, "uav2": 13}, map[int]int{1: 8, 2: 25}))
	uavs = append(uavs, newUav("uav4", map[string]int{"uav1": 11}, map[int]int{1: 31, 2: 20}))
	/*
		//测试1：正常测试
		tasks = append(tasks, newTask("task1", "-1", map[int]int{1: 5, 2: 7}, 2, 1))
		tasks = append(tasks, newTask("task2", "task1", map[int]int{1: 8, 2: 13}, 2, 2))
		tasks = append(tasks, newTask("task3", "-1", map[int]int{1: 2, 2: 3}, 5, 1))
		tasks = append(tasks, newTask("task4", "task3", map[int]int{1: 13, 2: 5}, 5, 2))
		tasks = append(tasks, newTask("task5", "task4", map[int]int{1: 6, 2: 10}, 5, 3))

		tasks = append(tasks, newTask("task6", "-1", map[int]int{1: 30, 2: 15}, 6, 0))
		tasks = append(tasks, newTask("task7", "-1", map[int]int{1: 7, 2: 8}, 5, 0))
		tasks = append(tasks, newTask("task8", "-1", map[int]int{1: 5, 2: 6}, 4, 0))
		tasks = append(tasks, newTask("task9", "-1", map[int]int{1: 3, 2: 5}, 3, 0))
		tasks = append(tasks, newTask("task10", "-1", map[int]int{1: 6, 2: 6}, 2, 0))
	*/
	// 测试2：测试时序任务是否正常被分配
	tasks = append(tasks, newTask("task1", "-1", map[int]int{1: 5, 2: 7}, 2, 1))
	tasks = append(tasks, newTask("task2", "task1", map[int]int{1: 3, 2: 3}, 2, 2))
	tasks = append(tasks, newTask("task3", "-1", map[int]int{1: 2, 2: 3}, 5, 1))
	tasks = append(tasks, newTask("task4", "task3", map[int]int{1: 3, 2: 3}, 5, 2))
	tasks = append(tasks, newTask("task5", "task4", map[int]int{1: 2, 2: 3}, 5, 3))

	tasks = append(tasks, newTask("task6", "-1", map[int]int{1: 30, 2: 15}, 6, 0))
	tasks = append(tasks, newTask("task7", "-1", map[int]int{1: 7, 2: 8}, 5, 0))
	tasks = append(tasks, newTask("task8", "-1", map[int]int{1: 5, 2: 6}, 4, 0))
	tasks = append(tasks, newTask("task9", "-1", map[int]int{1: 3, 2: 5}, 3, 0))
	tasks = append(tasks, newTask("task10", "-1", map[int]int{1: 6, 2: 6}, 2, 0))

	result, unassigned := KM.GraphMatch(uavs, tasks)

	for j, i := range result {
		fmt.Printf("任务%s->无人机%s\n", j, i)
	}
	for i, _ := range unassigned {
		fmt.Printf("任务%s未分配\n", unassigned[i].TaskID)
	}

	for uav := range uavs {
		printUav(*uavs[uav])
	}
}
func printUav(uav typeStruct.Uav) {
	fmt.Printf("UID：%s", uav.Uid)
	fmt.Printf("资源：%s", uav.Resources)
	fmt.Printf("下一跳：%s", uav.NextUavs)
	fmt.Printf("装载任务：%s\n", uav.LoadedTasks)

}

func main() {
	start := time.Now()
	test1()

	fmt.Println("耗时：", time.Since(start))

}
