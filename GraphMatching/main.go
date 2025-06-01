package main

import (
	"GraphMatching/contrastExperiment"
	"GraphMatching/mathingAlgorithm/KM"
	"GraphMatching/typeStruct"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// 可行性实验
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

	tasks = append(tasks, newTask("task6", "-1", map[int]int{1: 12, 2: 6}, 6, 0))
	tasks = append(tasks, newTask("task7", "-1", map[int]int{1: 7, 2: 8}, 5, 0))
	tasks = append(tasks, newTask("task8", "-1", map[int]int{1: 5, 2: 6}, 4, 0))
	tasks = append(tasks, newTask("task9", "-1", map[int]int{1: 3, 2: 5}, 3, 0))
	tasks = append(tasks, newTask("task10", "-1", map[int]int{1: 6, 2: 6}, 2, 0))

	fmt.Println("初始无人机")
	for uav := range uavs {
		printUav(*uavs[uav])
	}
	copyuavs := CopyUavs(uavs)
	//fmt.Println("拷贝无人机")
	//for uav := range uavs {
	//	printUav(*copyuavs[uav])
	//}

	result, unassigned := KM.GraphMatch(copyuavs, tasks)
	//result, unassigned := contrastExperiment.Dogreedy(copyuavs, tasks)
	//result, unassigned := contrastExperiment.SimpleCombination(copyuavs, tasks)
	//result, unassigned := contrastExperiment.DoAuction(copyuavs, tasks)

	for j, i := range result {
		fmt.Printf("任务%s->无人机%s\n", j, i)
	}
	for i, _ := range unassigned {
		fmt.Printf("任务%s未分配\n", unassigned[i].TaskID)
	}

	//fmt.Println("初始无人机")
	//for uav := range uavs {
	//	printUav(*uavs[uav])
	//}
	fmt.Println("拷贝无人机")
	for uav, _ := range uavs {
		printUav(*copyuavs[uav])
	}
	assess, f := Assess(uavs, tasks, result, copyuavs)
	fmt.Println("评估结果：", assess, f)
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
func printTask(task typeStruct.Task) {
	fmt.Printf("TaskID：%s\t", task.TaskID)
	fmt.Printf("\t任务优先级：%d", task.Priority)
	fmt.Printf("\t任务类型：%d", task.TaskType)
	fmt.Printf("\t上游任务：%s", task.Prev)

	fmt.Printf("\t所需资源：")
	for i, _ := range task.NeedResources {
		fmt.Printf("(%d:%d)", i, task.NeedResources)
	}

	fmt.Println()
}

func main() {
	test1()

	//GenerateTasksAndUavs(4, 10)

	//contrastTest(3, 10, 5)
	//contrastTestWithoutCombination(25, 100, 5)
	//contrastTestWithoutCombination(100, 300, 5)
}

//对比实验

// 深拷贝初始无人机
func CopyUavs(uavs []*typeStruct.Uav) []*typeStruct.Uav {
	copiedUavs := make([]*typeStruct.Uav, len(uavs))
	for i, uav := range uavs {
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

// 随机生成m台无人机n个任务，无人机资源在80~100，任务有三个规模按70、20、10
func GenerateTasksAndUavs(m int, n int) ([]*typeStruct.Uav, []typeStruct.Task) {
	rand.Seed(time.Now().UnixNano())

	// 生成无人机
	uavs := make([]*typeStruct.Uav, m)
	for i := 0; i < m; i++ {
		id := fmt.Sprintf("uav%d", i+1)
		resources := map[int]int{
			1: rand.Intn(21) + 80, // 资源1: 80-100
			2: rand.Intn(21) + 80, // 资源2: 80-100
		}
		nextUavs := make(map[string]int)
		// 每个无人机随机选择2-3个通信目标
		//numConnections := rand.Intn(2) + 2
		numConnections := rand.Intn(max(1, m/10)) + 1
		for j := 0; j < numConnections; j++ {
			otherID := fmt.Sprintf("uav%d", rand.Intn(m)+1)
			if otherID != id {
				nextUavs[otherID] = rand.Intn(10) + 1 // 通信能力1-10
			}
		}
		uavs[i] = &typeStruct.Uav{
			Uid:       id,
			Resources: resources,
			NextUavs:  nextUavs,
		}
	}

	// 生成任务
	//第一阶段，生成id、优先级、资源需求
	tasks := make([]typeStruct.Task, n)
	alltasks := make(map[string]typeStruct.Task)
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("task%d", i+1)

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

		tasks[i] = typeStruct.Task{
			TaskID:        id,
			TaskType:      0,
			Priority:      priority,
			Prev:          "-1",
			NeedResources: needResources,
		}
		alltasks[id] = tasks[i]
	}
	//第二阶段，生成时序依赖
	type0 := int(float64(n) * 0.4)
	type1 := int(float64(n) * 0.3)
	type2 := int(float64(n) * 0.2)
	//type3 := int(float64(n)*0.1)
	index := 0
	taskpool1 := make([]string, 0)
	taskpool2 := make([]string, 0)
	//处理type1
	for i := type0; i < type0+type1; i++ {
		tasks[i].TaskType = 1
		taskpool1 = append(taskpool1, tasks[i].TaskID)
	}
	//处理type2
	for i := type0 + type1; i < type0+type1+type2; i++ {
		tasks[i].TaskType = 2
		tasks[i].Prev = taskpool1[index]
		index++
		taskpool2 = append(taskpool2, tasks[i].TaskID)
	}

	index = 0
	//处理type3
	for i := type0 + type1 + type2; i < n; i++ {
		tasks[i].TaskType = 3
		tasks[i].Prev = taskpool2[index]
		index++
	}

	//TODO 搞一个超级无人机
	uavs[0].Resources[1] = 500 //int(math.Pow(10, 3))
	uavs[0].Resources[2] = 500 //int(math.Pow(10, 3))
	for _, uav := range uavs {
		uavs[0].NextUavs[uav.Uid] = 100
		uav.NextUavs[uavs[0].Uid] = 100
	}

	fmt.Println("生成无人机与任务成功！")
	for _, uav := range uavs {
		printUav(*uav)
	}
	for _, task := range tasks {
		printTask(task)
	}

	return uavs, tasks
}

// 评价一个集群的负载均衡度以及任务完成率
func Assess(olduavs []*typeStruct.Uav, tasks []typeStruct.Task, res map[string]string, newuavs []*typeStruct.Uav) (float64, float64) {
	//负载均衡度，方差表示
	clb := 0.0
	oldUavPool := make(map[string]*typeStruct.Uav)
	for _, uav := range olduavs {
		oldUavPool[uav.Uid] = uav
	}
	/*
		utils := make([]float64, len(newuavs))
		for _, newuav := range newuavs {
			utilx := 0.0
			for resType, resSum := range newuav.Resources {
				utilx += float64((oldUavPool[newuav.Uid].Resources[resType] - resSum)) / float64(oldUavPool[newuav.Uid].Resources[resType])
			}
			utilx /= float64(len(newuav.Resources))

			//utilx *= float64(len(newuav.LoadedTasks))

			utils = append(utils, utilx)
		}
	*/
	utils := make(map[int]float64, len(newuavs))
	for i, newuav := range newuavs {
		utilx := 0.0
		for resType, resSum := range newuav.Resources {
			utilx += float64((oldUavPool[newuav.Uid].Resources[resType] - resSum)) / float64(oldUavPool[newuav.Uid].Resources[resType])
		}
		utilx /= float64(len(newuav.Resources))

		utilx *= math.Abs(float64(len(newuav.LoadedTasks)) - float64(len(tasks)/len(olduavs)))

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
	assignedRate := float64(len(res)) / float64(len(tasks))
	return clb, assignedRate
}

// 对比实验入口
func contrastTest(m int, n int, count int) {
	//TODO 不带组合的
	/*
		fmt.Println("对比试验开始----------")
		var kmTime, greedyTime float64
		kmTime = 0.0
		greedyTime = 0.0
		//combinationTime = 0.0

		kmClb := 0.0
		greedyClb := 0.0
		//combinationClb := 0.0

		kmAssignedRate := 0.0
		greedyAssignedRate := 0.0
		//combinationAssignedRate := 0.0
		for i := 1; i <= count; i++ {
			fmt.Println("第", i, "次实验")
			uavs, tasks := GenerateTasksAndUavs(m, n)
			fmt.Println("初始无人机状态：")
			printUavs(uavs)

			copyuavs1 := CopyUavs(uavs)
			start := time.Now()
			result1, unassigned1 := KM.GraphMatch(copyuavs1, tasks)
			time1 := time.Since(start)
			clb1, assignedRate1 := Assess(uavs, tasks, result1, copyuavs1)

			copyuavs2 := CopyUavs(uavs)
			start = time.Now()
			result2, unassigned2 := contrastExperiment.Dogreedy(copyuavs2, tasks)
			time2 := time.Since(start)
			clb2, assignedRate2 := Assess(uavs, tasks, result2, copyuavs2)

			//copyuavs3 := CopyUavs(uavs)
			//start = time.Now()
			//result3, unassigned3 := contrastExperiment.SimpleCombination(copyuavs3, tasks)
			//time3 := time.Since(start)
			//clb3, assignedRate3 := Assess(uavs, tasks, result3, copyuavs3)

			fmt.Println("图匹配算法结果：")
			printResult(result1, unassigned1)
			printUavs(copyuavs1)

			fmt.Println("贪心算法结果：")
			printResult(result2, unassigned2)
			printUavs(copyuavs2)

			//fmt.Println("组合算法结果：")
			//printResult(result3, unassigned3)
			//printUavs(copyuavs3)

			fmt.Println("图匹配耗时、负载均衡度、任务完成率为：", time1, clb1, assignedRate1)
			fmt.Println("贪心算法耗时、负载均衡度、任务完成率为：", time2, clb2, assignedRate2)
			//fmt.Println("组合算法耗时、负载均衡度、任务完成率为：", time3, clb3, assignedRate3)

			kmTime += float64(time1.Milliseconds())
			kmClb += clb1
			kmAssignedRate += assignedRate1

			greedyTime += float64(time2.Milliseconds())
			greedyClb += clb2
			greedyAssignedRate += assignedRate2

			//combinationTime += float64(time3.Milliseconds())
			//combinationClb += clb3
			//combinationAssignedRate += assignedRate3

		}
		fmt.Println("\n\n对比实验结束，结果如下：")
		fmt.Println("图匹配平均耗时（毫秒）、负载均衡度、任务完成率为：", kmTime/float64(count), kmClb/float64(count), kmAssignedRate/float64(count))
		fmt.Println("贪心平均耗时（毫秒）、负载均衡度、任务完成率为：", greedyTime/float64(count), greedyClb/float64(count), greedyAssignedRate/float64(count))
		//fmt.Println("组合平均耗时（毫秒）、负载均衡度、任务完成率为：", combinationTime/float64(count), combinationClb/float64(count), combinationAssignedRate/float64(count))
	*/

	//TODO 带组合的

	fmt.Println("对比试验开始----------")
	var kmTime, greedyTime, combinationTime, auctionTime float64
	kmTime = 0.0
	greedyTime = 0.0
	combinationTime = 0.0
	auctionTime = 0.0

	kmClb := 0.0
	greedyClb := 0.0
	combinationClb := 0.0
	auctionClb := 0.0

	kmAssignedRate := 0.0
	greedyAssignedRate := 0.0
	combinationAssignedRate := 0.0
	auctionAssignedRate := 0.0

	for i := 1; i <= count; i++ {
		fmt.Println("第", i, "次实验")
		uavs, tasks := GenerateTasksAndUavs(m, n)
		fmt.Println("初始无人机状态：")
		printUavs(uavs)

		copyuavs1 := CopyUavs(uavs)
		start := time.Now()
		result1, unassigned1 := KM.GraphMatch(copyuavs1, tasks)
		time1 := time.Since(start)
		clb1, assignedRate1 := Assess(uavs, tasks, result1, copyuavs1)

		copyuavs2 := CopyUavs(uavs)
		start = time.Now()
		result2, unassigned2 := contrastExperiment.Dogreedy(copyuavs2, tasks)
		time2 := time.Since(start)
		clb2, assignedRate2 := Assess(uavs, tasks, result2, copyuavs2)

		copyuavs3 := CopyUavs(uavs)
		start = time.Now()
		result3, unassigned3 := contrastExperiment.SimpleCombination(copyuavs3, tasks)
		time3 := time.Since(start)
		clb3, assignedRate3 := Assess(uavs, tasks, result3, copyuavs3)

		copyuavs4 := CopyUavs(uavs)
		start = time.Now()
		result4, unassigned4 := contrastExperiment.DoAuction(copyuavs4, tasks)
		time4 := time.Since(start)
		clb4, assignedRate4 := Assess(uavs, tasks, result4, copyuavs4)

		fmt.Println("图匹配算法结果：")
		printResult(result1, unassigned1)
		printUavs(copyuavs1)

		fmt.Println("贪心算法结果：")
		printResult(result2, unassigned2)
		printUavs(copyuavs2)

		fmt.Println("组合算法结果：")
		printResult(result3, unassigned3)
		printUavs(copyuavs3)

		fmt.Println("拍卖算法结果：")
		printResult(result4, unassigned4)
		printUavs(copyuavs4)

		fmt.Println("图匹配耗时、负载均衡度、任务完成率为：", time1, clb1, assignedRate1)
		fmt.Println("贪心算法耗时、负载均衡度、任务完成率为：", time2, clb2, assignedRate2)
		fmt.Println("组合算法耗时、负载均衡度、任务完成率为：", time3, clb3, assignedRate3)
		fmt.Println("拍卖算法耗时、负载均衡度、任务完成率为：", time4, clb4, assignedRate4)

		kmTime += float64(time1.Milliseconds())
		kmClb += clb1
		kmAssignedRate += assignedRate1

		greedyTime += float64(time2.Milliseconds())
		greedyClb += clb2
		greedyAssignedRate += assignedRate2

		combinationTime += float64(time3.Milliseconds())
		combinationClb += clb3
		combinationAssignedRate += assignedRate3

		auctionTime += float64(time4.Milliseconds())
		auctionClb += clb4
		auctionAssignedRate += assignedRate4
	}
	fmt.Println("\n\n对比实验结束，结果如下：")
	fmt.Println("图匹配平均耗时（毫秒）、负载均衡度、任务完成率为：", kmTime/float64(count), kmClb/float64(count), kmAssignedRate/float64(count))
	fmt.Println("贪心平均耗时（毫秒）、负载均衡度、任务完成率为：", greedyTime/float64(count), greedyClb/float64(count), greedyAssignedRate/float64(count))
	fmt.Println("组合平均耗时（毫秒）、负载均衡度、任务完成率为：", combinationTime/float64(count), combinationClb/float64(count), combinationAssignedRate/float64(count))
	fmt.Println("拍卖平均耗时（毫秒）、负载均衡度、任务完成率为：", auctionTime/float64(count), auctionClb/float64(count), auctionAssignedRate/float64(count))

}

// 对比实验入口
func contrastTestWithoutCombination(m int, n int, count int) {
	//TODO 不带组合的

	fmt.Println("对比试验开始----------")
	var kmTime, greedyTime, auctionTime float64
	kmTime = 0.0
	greedyTime = 0.0
	//combinationTime = 0.0
	auctionTime = 0.0

	kmClb := 0.0
	greedyClb := 0.0
	//combinationClb := 0.0
	auctionClb := 0.0

	kmAssignedRate := 0.0
	greedyAssignedRate := 0.0
	//combinationAssignedRate := 0.0
	auctionAssignedRate := 0.0
	for i := 1; i <= count; i++ {
		fmt.Println("第", i, "次实验")
		uavs, tasks := GenerateTasksAndUavs(m, n)
		fmt.Println("初始无人机状态：")
		printUavs(uavs)

		copyuavs1 := CopyUavs(uavs)
		start := time.Now()
		result1, unassigned1 := KM.GraphMatch(copyuavs1, tasks)
		time1 := time.Since(start)
		clb1, assignedRate1 := Assess(uavs, tasks, result1, copyuavs1)

		copyuavs2 := CopyUavs(uavs)
		start = time.Now()
		result2, unassigned2 := contrastExperiment.Dogreedy(copyuavs2, tasks)
		time2 := time.Since(start)
		clb2, assignedRate2 := Assess(uavs, tasks, result2, copyuavs2)

		//copyuavs3 := CopyUavs(uavs)
		//start = time.Now()
		//result3, unassigned3 := contrastExperiment.SimpleCombination(copyuavs3, tasks)
		//time3 := time.Since(start)
		//clb3, assignedRate3 := Assess(uavs, tasks, result3, copyuavs3)

		copyuavs4 := CopyUavs(uavs)
		start = time.Now()
		result4, unassigned4 := contrastExperiment.DoAuction(copyuavs4, tasks)
		time4 := time.Since(start)
		clb4, assignedRate4 := Assess(uavs, tasks, result4, copyuavs4)

		fmt.Println("图匹配算法结果：")
		printResult(result1, unassigned1)
		printUavs(copyuavs1)

		fmt.Println("贪心算法结果：")
		printResult(result2, unassigned2)
		printUavs(copyuavs2)

		//fmt.Println("组合算法结果：")
		//printResult(result3, unassigned3)
		//printUavs(copyuavs3)

		fmt.Println("拍卖算法结果：")
		printResult(result4, unassigned4)
		printUavs(copyuavs4)

		fmt.Println("图匹配耗时、负载均衡度、任务完成率为：", time1, clb1, assignedRate1)
		fmt.Println("贪心算法耗时、负载均衡度、任务完成率为：", time2, clb2, assignedRate2)
		//fmt.Println("组合算法耗时、负载均衡度、任务完成率为：", time3, clb3, assignedRate3)
		fmt.Println("拍卖算法耗时、负载均衡度、任务完成率为：", time4, clb4, assignedRate4)

		kmTime += float64(time1.Milliseconds())
		kmClb += clb1
		kmAssignedRate += assignedRate1

		greedyTime += float64(time2.Milliseconds())
		greedyClb += clb2
		greedyAssignedRate += assignedRate2

		//combinationTime += float64(time3.Milliseconds())
		//combinationClb += clb3
		//combinationAssignedRate += assignedRate3

		auctionTime += float64(time4.Milliseconds())
		auctionClb += clb4
		auctionAssignedRate += assignedRate4
	}
	fmt.Println("\n\n对比实验结束，结果如下：")
	fmt.Println("图匹配平均耗时（毫秒）、负载均衡度、任务完成率为：", kmTime/float64(count), kmClb/float64(count), kmAssignedRate/float64(count))
	fmt.Println("贪心平均耗时（毫秒）、负载均衡度、任务完成率为：", greedyTime/float64(count), greedyClb/float64(count), greedyAssignedRate/float64(count))
	//fmt.Println("组合平均耗时（毫秒）、负载均衡度、任务完成率为：", combinationTime/float64(count), combinationClb/float64(count), combinationAssignedRate/float64(count))
	fmt.Println("拍卖平均耗时（毫秒）、负载均衡度、任务完成率为：", auctionTime/float64(count), auctionClb/float64(count), auctionAssignedRate/float64(count))

	//TODO 带组合的
	/*
		fmt.Println("对比试验开始----------")
		var kmTime, greedyTime, combinationTime float64
		kmTime = 0.0
		greedyTime = 0.0
		combinationTime = 0.0

		kmClb := 0.0
		greedyClb := 0.0
		combinationClb := 0.0

		kmAssignedRate := 0.0
		greedyAssignedRate := 0.0
		combinationAssignedRate := 0.0
		for i := 1; i <= count; i++ {
			fmt.Println("第", i, "次实验")
			uavs, tasks := GenerateTasksAndUavs(m, n)
			fmt.Println("初始无人机状态：")
			printUavs(uavs)

			copyuavs1 := CopyUavs(uavs)
			start := time.Now()
			result1, unassigned1 := KM.GraphMatch(copyuavs1, tasks)
			time1 := time.Since(start)
			clb1, assignedRate1 := Assess(uavs, tasks, result1, copyuavs1)

			copyuavs2 := CopyUavs(uavs)
			start = time.Now()
			result2, unassigned2 := contrastExperiment.Dogreedy(copyuavs2, tasks)
			time2 := time.Since(start)
			clb2, assignedRate2 := Assess(uavs, tasks, result2, copyuavs2)

			copyuavs3 := CopyUavs(uavs)
			start = time.Now()
			result3, unassigned3 := contrastExperiment.SimpleCombination(copyuavs3, tasks)
			time3 := time.Since(start)
			clb3, assignedRate3 := Assess(uavs, tasks, result3, copyuavs3)

			fmt.Println("图匹配算法结果：")
			printResult(result1, unassigned1)
			printUavs(copyuavs1)

			fmt.Println("贪心算法结果：")
			printResult(result2, unassigned2)
			printUavs(copyuavs2)

			fmt.Println("组合算法结果：")
			printResult(result3, unassigned3)
			printUavs(copyuavs3)

			fmt.Println("图匹配耗时、负载均衡度、任务完成率为：", time1, clb1, assignedRate1)
			fmt.Println("贪心算法耗时、负载均衡度、任务完成率为：", time2, clb2, assignedRate2)
			fmt.Println("组合算法耗时、负载均衡度、任务完成率为：", time3, clb3, assignedRate3)

			kmTime += float64(time1.Milliseconds())
			kmClb += clb1
			kmAssignedRate += assignedRate1

			greedyTime += float64(time2.Milliseconds())
			greedyClb += clb2
			greedyAssignedRate += assignedRate2

			combinationTime += float64(time3.Milliseconds())
			combinationClb += clb3
			combinationAssignedRate += assignedRate3

		}
		fmt.Println("\n\n对比实验结束，结果如下：")
		fmt.Println("图匹配平均耗时（毫秒）、负载均衡度、任务完成率为：", kmTime/float64(count), kmClb/float64(count), kmAssignedRate/float64(count))
		fmt.Println("贪心平均耗时（毫秒）、负载均衡度、任务完成率为：", greedyTime/float64(count), greedyClb/float64(count), greedyAssignedRate/float64(count))
		fmt.Println("组合平均耗时（毫秒）、负载均衡度、任务完成率为：", combinationTime/float64(count), combinationClb/float64(count), combinationAssignedRate/float64(count))
	*/
}

func printResult(result map[string]string, unassigned []typeStruct.Task) {
	for j, i := range result {
		fmt.Printf("任务%s->无人机%s\n", j, i)
	}
	for i, _ := range unassigned {
		fmt.Printf("任务%s未分配\n", unassigned[i].TaskID)
	}
}
