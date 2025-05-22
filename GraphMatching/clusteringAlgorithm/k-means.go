package clusteringAlgorithm

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// 定义任务与无人机结构体
type Task struct {
	ID            string
	Prev, Next    string          // 依赖关系
	Priority      float64         // 优先级
	NeedResources map[int]float64 // 资源需求（资源类型->数量）
}
type Uav struct {
	ID        string
	Resources map[int]float64    // 资源能力（资源类型->数量）
	NextUavs  map[string]float64 // 与其他无人机的通信能力（无人机ID->连接度）
}

// 构建任务依赖图并获取连通组件
func getTaskComponents(tasks []Task) [][]Task {
	// 建立任务ID到索引的映射
	idIndex := make(map[string]int)
	for i, t := range tasks {
		idIndex[t.ID] = i
	}
	// 建立无向邻接列表
	adj := make(map[int][]int)
	for i, t := range tasks {
		if t.Prev != "" {
			if j, ok := idIndex[t.Prev]; ok {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
		if t.Next != "" {
			if j, ok := idIndex[t.Next]; ok {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}
	// DFS寻找连通组件
	visited := make([]bool, len(tasks))
	var components [][]Task
	for i := range tasks {
		if !visited[i] {
			stack := []int{i}
			visited[i] = true
			compIdx := []int{i}
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, v := range adj[u] {
					if !visited[v] {
						visited[v] = true
						stack = append(stack, v)
						compIdx = append(compIdx, v)
					}
				}
			}
			// 将索引转换为Task列表
			compTasks := make([]Task, len(compIdx))
			for k, idx := range compIdx {
				compTasks[k] = tasks[idx]
			}
			components = append(components, compTasks)
		}
	}
	return components
}

// 计算每个组件的特征向量（资源需求总和、优先级平均值）
func compFeatures(comps [][]Task) [][]float64 {
	// 先统计所有可能的资源类型
	resourceSet := make(map[int]bool)
	for _, comp := range comps {
		for _, t := range comp {
			for r := range t.NeedResources {
				resourceSet[r] = true
			}
		}
	}
	resKeys := make([]int, 0, len(resourceSet))
	for r := range resourceSet {
		resKeys = append(resKeys, r)
	}
	// 对每个组件计算特征向量
	features := make([][]float64, len(comps))
	for i, comp := range comps {
		feat := make([]float64, len(resKeys)+1) // 最后一维为平均优先级
		// 资源需求总和
		for _, t := range comp {
			for j, r := range resKeys {
				feat[j] += t.NeedResources[r]
			}
			feat[len(resKeys)] += t.Priority
		}
		// 计算平均优先级
		if len(comp) > 0 {
			feat[len(resKeys)] /= float64(len(comp))
		}
		features[i] = feat
	}
	return features
}

// 对无人机计算特征向量（资源能力、通信度）
func uavFeatures(uavs []*Uav) [][]float64 {
	// 统计所有资源类型
	resourceSet := make(map[int]bool)
	for _, u := range uavs {
		for r := range u.Resources {
			resourceSet[r] = true
		}
	}
	resKeys := make([]int, 0, len(resourceSet))
	for r := range resourceSet {
		resKeys = append(resKeys, r)
	}
	// 计算特征：资源能力和通信度
	features := make([][]float64, len(uavs))
	for i, u := range uavs {
		feat := make([]float64, len(resKeys)+1) // 最后一维为通信度
		// 资源能力
		for j, r := range resKeys {
			feat[j] = u.Resources[r]
		}
		// 计算通信度（与所有邻接无人机通信能力之和）
		var commSum float64
		for _, cap := range u.NextUavs {
			commSum += cap
		}
		feat[len(resKeys)] = commSum
		features[i] = feat
	}
	return features
}

// K-means 聚类算法（数据点为行向量）
func kMeans(data [][]float64, K int, maxIter int) []int {
	n := len(data)
	if n == 0 || K <= 0 {
		return nil
	}
	dim := len(data[0])
	// 初始化质心：随机选择K个数据点
	centroids := make([][]float64, K)
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(n)
	for i := 0; i < K; i++ {
		centroids[i] = make([]float64, dim)
		copy(centroids[i], data[perm[i%len(perm)]])
	}
	labels := make([]int, n)
	for iter := 0; iter < maxIter; iter++ {
		changed := false
		// 为每个点分配最近的质心
		for i, x := range data {
			minDist := math.MaxFloat64
			var best int
			for j := 0; j < K; j++ {
				// 计算欧氏距离的平方
				dist := 0.0
				for d := 0; d < dim; d++ {
					diff := x[d] - centroids[j][d]
					dist += diff * diff
				}
				if dist < minDist {
					minDist = dist
					best = j
				}
			}
			if labels[i] != best {
				labels[i] = best
				changed = true
			}
		}
		// 如果无点改变簇标签，则提前退出
		if !changed && iter > 0 {
			break
		}
		// 更新质心：对每个簇内点取平均
		counts := make([]int, K)
		newCentroids := make([][]float64, K)
		for j := 0; j < K; j++ {
			newCentroids[j] = make([]float64, dim)
		}
		for i, x := range data {
			lab := labels[i]
			counts[lab]++
			for d := 0; d < dim; d++ {
				newCentroids[lab][d] += x[d]
			}
		}
		for j := 0; j < K; j++ {
			if counts[j] > 0 {
				for d := 0; d < dim; d++ {
					newCentroids[j][d] /= float64(counts[j])
				}
				centroids[j] = newCentroids[j]
			}
		}
	}
	return labels
}

// 对任务列表进行聚类，返回每个簇对应的任务列表
func clusterTasks(tasks []Task, K int) map[int][]Task {
	// 构建连通组件
	comps := getTaskComponents(tasks)
	// 计算组件特征
	feats := compFeatures(comps)
	// K-means聚类组件
	labels := kMeans(feats, K, 100)
	// 将任务分配到对应簇
	clusters := make(map[int][]Task)
	for i, comp := range comps {
		lab := labels[i]
		clusters[lab] = append(clusters[lab], comp...)
	}
	return clusters
}

// 对无人机列表进行聚类，返回每个簇对应的无人机列表
func clusterUavs(uavs []*Uav, K int) map[int][]*Uav {
	feats := uavFeatures(uavs)
	labels := kMeans(feats, K, 100)
	clusters := make(map[int][]*Uav)
	for i, lab := range labels {
		clusters[lab] = append(clusters[lab], uavs[i])
	}
	return clusters
}

// 示例：KM匹配算法（此处为占位，实际应用可调用现成KM库）
func KMGraphMatch(uavs []*Uav, tasks []Task) {
	// 这里应实现或调用二分图匹配算法，将任务分配给无人机
	// 省略详细实现
	fmt.Printf("执行KM算法：%d个无人机与%d个任务匹配\n", len(uavs), len(tasks))
}

func Main() {
	// 示例数据：初始化无人机和任务列表
	uavs := []*Uav{
		{ID: "uav1", Resources: map[int]float64{1: 15, 2: 25}, NextUavs: map[string]float64{"uav3": 12, "uav4": 11}},
		{ID: "uav2", Resources: map[int]float64{1: 24, 2: 16}, NextUavs: map[string]float64{"uav3": 13}},
		{ID: "uav3", Resources: map[int]float64{1: 8, 2: 25}, NextUavs: map[string]float64{"uav1": 12, "uav2": 13}},
		{ID: "uav4", Resources: map[int]float64{1: 31, 2: 20}, NextUavs: map[string]float64{"uav1": 11}},
	}
	tasks := []Task{
		{ID: "task1", Prev: "", Next: "task2", NeedResources: map[int]float64{1: 5, 2: 7}, Priority: 2},
		{ID: "task2", Prev: "task1", Next: "", NeedResources: map[int]float64{1: 8, 2: 13}, Priority: 2},
		{ID: "task3", Prev: "", Next: "task4", NeedResources: map[int]float64{1: 2, 2: 3}, Priority: 5},
		{ID: "task4", Prev: "task3", Next: "task5", NeedResources: map[int]float64{1: 13, 2: 5}, Priority: 5},
		{ID: "task5", Prev: "task4", Next: "", NeedResources: map[int]float64{1: 6, 2: 10}, Priority: 5},
		{ID: "task6", Prev: "", Next: "", NeedResources: map[int]float64{1: 30, 2: 15}, Priority: 6},
		{ID: "task7", Prev: "", Next: "", NeedResources: map[int]float64{1: 7, 2: 8}, Priority: 5},
	}

	K := 2 // 设置分组数，可根据需求调整

	// 对任务和无人机分别聚类
	taskClusters := clusterTasks(tasks, K)
	uavClusters := clusterUavs(uavs, K)

	// 对每一簇执行组内匹配
	for i := 0; i < K; i++ {
		tgroup := taskClusters[i]
		ugroup := uavClusters[i]
		fmt.Printf("第%d组：%d个任务，%d架无人机\n", i+1, len(tgroup), len(ugroup))
		// 调用KM算法进行组内任务分配
		KMGraphMatch(ugroup, tgroup)
	}
}
