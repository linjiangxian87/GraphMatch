/*
package main

import (

	"fmt"
	"math"
	"math/rand"
	"time"

)

// 任务结构体

	type Task struct {
		ID            string
		Prev, Next    string          // 依赖关系
		Priority      float64         // 优先级
		NeedResources map[int]float64 // 资源需求（资源类型->数量）
	}

// 无人机结构体

	type Uav struct {
		ID        string
		Resources map[int]float64    // 资源能力（资源类型->数量）
		NextUavs  map[string]float64 // 与其他无人机的通信能力（无人机ID->连接度）
	}

// 获取任务连通组件

	func getTaskComponents(tasks []Task) [][]Task {
		idIndex := make(map[string]int)
		for i, t := range tasks {
			idIndex[t.ID] = i
		}

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

		visited := make([]bool, len(tasks))
		components := [][]Task{}

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

				compTasks := make([]Task, len(compIdx))
				for k, idx := range compIdx {
					compTasks[k] = tasks[idx]
				}
				components = append(components, compTasks)
			}
		}
		return components
	}

// 计算组件特征向量

	func compFeatures(comps [][]Task) [][]float64 {
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

		features := make([][]float64, len(comps))
		for i, comp := range comps {
			feat := make([]float64, len(resKeys)+1) // 最后一维为平均优先级
			for _, t := range comp {
				for j, r := range resKeys {
					feat[j] += t.NeedResources[r]
				}
				feat[len(resKeys)] += t.Priority
			}
			if len(comp) > 0 {
				feat[len(resKeys)] /= float64(len(comp))
			}
			features[i] = feat
		}
		return features
	}

// 无人机特征向量

	func uavFeatures(uavs []*Uav) [][]float64 {
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

		features := make([][]float64, len(uavs))
		for i, u := range uavs {
			feat := make([]float64, len(resKeys)+1) // 最后一维为通信度
			for j, r := range resKeys {
				feat[j] = u.Resources[r]
			}

			var commSum float64
			for _, next := range u.NextUavs {
				commSum += next
			}
			feat[len(resKeys)] = commSum

			features[i] = feat
		}
		return features
	}

// K-means++聚类

	func kMeans(data [][]float64, K int, maxIter int) []int {
		n := len(data)
		if n == 0 || K <= 0 {
			return nil
		}
		dim := len(data[0])

		// K-means++ 初始化
		centroids := make([][]float64, K)
		selected := make([]bool, n)

		// 第一个质心随机选择
		rand.Seed(time.Now().UnixNano())
		firstIdx := rand.Intn(n)
		centroids[0] = make([]float64, dim)
		copy(centroids[0], data[firstIdx])
		selected[firstIdx] = true

		for i := 1; i < K; i++ {
			distances := make([]float64, n)
			totalDist := 0.0

			for j := 0; j < n; j++ {
				if selected[j] {
					distances[j] = 0
					continue
				}

				minDist := math.MaxFloat64
				for c := 0; c < i; c++ {
					dist := 0.0
					for d := 0; d < dim; d++ {
						diff := data[j][d] - centroids[c][d]
						dist += diff * diff
					}
					if dist < minDist {
						minDist = dist
					}
				}
				distances[j] = minDist
				totalDist += minDist
			}

			if totalDist == 0 {
				break
			}

			// 按距离概率选择下一个质心
			target := rand.Float64() * totalDist
			for j := 0; j < n; j++ {
				if !selected[j] {
					target -= distances[j]
					if target <= 0 {
						centroids[i] = make([]float64, dim)
						copy(centroids[i], data[j])
						selected[j] = true
						break
					}
				}
			}
		}

		labels := make([]int, n)
		for iter := 0; iter < maxIter; iter++ {
			changed := false

			// 分配点到最近质心
			for i, x := range data {
				minDist := math.MaxFloat64
				best := 0

				for j := 0; j < K; j++ {
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

			if !changed && iter > 0 {
				break
			}

			// 更新质心
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

// 任务聚类（返回簇和簇特征）

	func clusterTasks(tasks []Task, K int) (map[int][]Task, [][]float64) {
		comps := getTaskComponents(tasks)
		feats := compFeatures(comps)
		labels := kMeans(feats, K, 100)

		clusters := make(map[int][]Task)
		for i, comp := range comps {
			lab := labels[i]
			clusters[lab] = append(clusters[lab], comp...)
		}

		// 计算每个簇的特征（组件特征的平均）
		clusterFeats := make([][]float64, K)
		for lab := range clusters {
			var sum []float64
			count := 0

			for i, comp := range comps {
				if labels[i] == lab {
					if sum == nil {
						sum = make([]float64, len(feats[i]))
					}
					for j, v := range feats[i] {
						sum[j] += v
					}
					count++
				}
			}

			if count > 0 {
				for j := range sum {
					sum[j] /= float64(count)
				}
			}
			clusterFeats[lab] = sum
		}

		return clusters, clusterFeats
	}

// 无人机聚类（返回簇和簇特征）

	func clusterUavs(uavs []*Uav, K int) (map[int][]*Uav, [][]float64) {
		feats := uavFeatures(uavs)
		labels := kMeans(feats, K, 100)

		clusters := make(map[int][]*Uav)
		for i, lab := range labels {
			clusters[lab] = append(clusters[lab], uavs[i])
		}

		// 验证通信约束
		for {
			valid := true
			for _, ugroup := range clusters {
				for _, uav := range ugroup {
					hasConnection := false
					for _, other := range ugroup {
						if uav.ID != other.ID && uav.NextUavs[other.ID] > 0 {
							hasConnection = true
							break
						}
					}
					if !hasConnection {
						valid = false
						break
					}
				}
				if !valid {
					break
				}
			}
			if valid {
				break
			}

			// 重新聚类
			labels = kMeans(feats, K, 100)
			clusters = make(map[int][]*Uav)
			for i, lab := range labels {
				clusters[lab] = append(clusters[lab], uavs[i])
			}
		}

		// 计算每个簇的特征（无人机特征的平均）
		clusterFeats := make([][]float64, K)
		for lab := range clusters {
			var sum []float64
			count := 0

			for i, lab2 := range labels {
				if lab2 == lab {
					if sum == nil {
						sum = make([]float64, len(feats[i]))
					}
					for j, v := range feats[i] {
						sum[j] += v
					}
					count++
				}
			}

			if count > 0 {
				for j := range sum {
					sum[j] /= float64(count)
				}
			}
			clusterFeats[lab] = sum
		}

		return clusters, clusterFeats
	}

// 余弦相似度

	func cosineSimilarity(a, b []float64) float64 {
		dot := 0.0
		normA := 0.0
		normB := 0.0

		for i := 0; i < len(a); i++ {
			dot += a[i] * b[i]
			normA += a[i] * a[i]
			normB += b[i] * b[i]
		}

		if normA == 0 || normB == 0 {
			return 0
		}

		return dot / (math.Sqrt(normA) * math.Sqrt(normB))
	}

// 任务-无人机簇配对

	func matchClusters(taskFeats, uavFeats [][]float64) []int {
		// 创建相似度矩阵
		simMatrix := make([][]float64, len(taskFeats))
		for i := range simMatrix {
			simMatrix[i] = make([]float64, len(uavFeats))
		}

		for i, tf := range taskFeats {
			for j, uf := range uavFeats {
				simMatrix[i][j] = cosineSimilarity(tf, uf)
			}
		}

		// 实际应用中应使用KM算法进行匹配
		// 此处简化为顺序配对
		matches := make([]int, len(taskFeats))
		for i := range matches {
			matches[i] = i % len(uavFeats)
		}

		return matches
	}

// 示例：KM匹配算法（占位）

	func KMGraphMatch(uavs []*Uav, tasks []Task) {
		fmt.Printf("执行KM算法：%d个无人机与%d个任务匹配\n", len(uavs), len(tasks))
	}

	func main() {
		// 示例数据
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

		K := 2

		// 任务聚类
		taskClusters, taskClusterFeats := clusterTasks(tasks, K)

		// 无人机聚类
		uavClusters, uavClusterFeats := clusterUavs(uavs, K)

		// 簇配对
		matches := matchClusters(taskClusterFeats, uavClusterFeats)

		// 组内匹配
		for i, match := range matches {
			tgroup := taskClusters[i]
			ugroup := uavClusters[match]
			fmt.Printf("第%d组配对：%d个任务，%d架无人机\n", i+1, len(tgroup), len(ugroup))
			KMGraphMatch(ugroup, tgroup)
		}
	}
*/
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Drone 结构体定义
type Drone struct {
	ID       int
	Resource float64
	NextUavs []int
}

// Task 结构体定义
type Task struct {
	ID          int
	ResourceReq float64
	Priority    int
	TaskType    int
	Prev        []int
}

// K-means++ 初始化, 返回初始质心索引
func initCentroidsKMeansPP(points [][]float64, K int) []int {
	n := len(points)
	if K <= 0 || K > n {
		return nil
	}
	centroids := make([]int, 0, K)
	rand.Seed(time.Now().UnixNano())
	// 随机选择第一个质心
	first := rand.Intn(n)
	centroids = append(centroids, first)

	dist := make([]float64, n)
	for len(centroids) < K {
		var sum float64
		// 计算每个点到最近质心的距离平方
		for i := 0; i < n; i++ {
			minDist := math.Inf(1)
			for _, c := range centroids {
				d := euclideanDistance(points[i], points[c])
				if d < minDist {
					minDist = d
				}
			}
			dist[i] = minDist * minDist
			sum += dist[i]
		}
		// 按照概率选择新的质心
		r := rand.Float64() * sum
		acc := 0.0
		next := 0
		for i := 0; i < n; i++ {
			acc += dist[i]
			if acc >= r {
				next = i
				break
			}
		}
		centroids = append(centroids, next)
	}
	return centroids
}

// 计算欧氏距离
func euclideanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

// KMeans 聚类 (返回每个点的簇标签)
func kMeansCluster(points [][]float64, K int, maxIter int) []int {
	n := len(points)
	centroidIndices := initCentroidsKMeansPP(points, K)
	if centroidIndices == nil {
		return nil
	}
	// 初始化质心坐标
	centroids := make([][]float64, K)
	for i, idx := range centroidIndices {
		centroids[i] = make([]float64, len(points[idx]))
		copy(centroids[i], points[idx])
	}
	labels := make([]int, n)
	for iter := 0; iter < maxIter; iter++ {
		// Assignment: 分配每个点到最近质心
		for i := 0; i < n; i++ {
			minDist := math.Inf(1)
			label := 0
			for k := 0; k < K; k++ {
				d := euclideanDistance(points[i], centroids[k])
				if d < minDist {
					minDist = d
					label = k
				}
			}
			labels[i] = label
		}
		// Update: 重新计算质心
		newCentroids := make([][]float64, K)
		counts := make([]int, K)
		dim := len(points[0])
		for k := 0; k < K; k++ {
			newCentroids[k] = make([]float64, dim)
		}
		for i := 0; i < n; i++ {
			k := labels[i]
			for d := 0; d < dim; d++ {
				newCentroids[k][d] += points[i][d]
			}
			counts[k]++
		}
		for k := 0; k < K; k++ {
			if counts[k] > 0 {
				for d := 0; d < dim; d++ {
					newCentroids[k][d] /= float64(counts[k])
				}
			}
		}
		// 检查质心是否收敛
		unchanged := true
		for k := 0; k < K; k++ {
			for d := 0; d < dim; d++ {
				if math.Abs(newCentroids[k][d]-centroids[k][d]) > 1e-6 {
					unchanged = false
					break
				}
			}
			if !unchanged {
				break
			}
		}
		centroids = newCentroids
		if unchanged {
			break
		}
	}
	return labels
}

// 判断 slice 是否包含某个值
func contains(slice []int, val int) bool {
	for _, x := range slice {
		if x == val {
			return true
		}
	}
	return false
}

// 检查簇内连通性（BFS）
func isClusterConnected(cluster []int, nextMap map[int][]int) bool {
	if len(cluster) == 0 {
		return true
	}
	visited := make(map[int]bool)
	queue := []int{cluster[0]}
	visited[cluster[0]] = true
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, nb := range nextMap[cur] {
			if contains(cluster, nb) && !visited[nb] {
				visited[nb] = true
				queue = append(queue, nb)
			}
		}
	}
	// 检查簇内所有节点是否被访问
	for _, v := range cluster {
		if !visited[v] {
			return false
		}
	}
	return true
}

// 根据任务依赖构建连通组件
func buildTaskComponents(tasks []Task) map[int][]int {
	parent := make(map[int]int)
	for _, t := range tasks {
		parent[t.ID] = t.ID
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(x, y int) {
		rx, ry := find(x), find(y)
		if rx != ry {
			parent[ry] = rx
		}
	}
	// 根据 Prev 关联任务
	for _, t := range tasks {
		for _, p := range t.Prev {
			union(t.ID, p)
		}
	}
	comp := make(map[int][]int)
	for _, t := range tasks {
		root := find(t.ID)
		comp[root] = append(comp[root], t.ID)
	}
	return comp
}

// KM 匹配算法接口占位
func matchDroneTask(droneIDs, taskIDs []int) {
	fmt.Printf("调用KM算法匹配 无人机组%s 与 任务组%s\n",
		fmt.Sprint(droneIDs), fmt.Sprint(taskIDs))
}

func main() {
	// 模拟数据
	drones := []Drone{
		{ID: 1, Resource: 10, NextUavs: []int{2, 3}},
		{ID: 2, Resource: 9, NextUavs: []int{1, 3}},
		{ID: 3, Resource: 11, NextUavs: []int{1, 2}},
		{ID: 4, Resource: 50, NextUavs: []int{5, 6}},
		{ID: 5, Resource: 45, NextUavs: []int{4, 6}},
		{ID: 6, Resource: 48, NextUavs: []int{4, 5}},
	}
	tasks := []Task{
		{ID: 101, ResourceReq: 7, Priority: 1, TaskType: 1, Prev: []int{}},
		{ID: 102, ResourceReq: 6, Priority: 2, TaskType: 1, Prev: []int{101}},
		{ID: 103, ResourceReq: 8, Priority: 1, TaskType: 2, Prev: []int{102}},
		{ID: 104, ResourceReq: 5, Priority: 2, TaskType: 3, Prev: []int{}},
		{ID: 105, ResourceReq: 6, Priority: 3, TaskType: 3, Prev: []int{104}},
		{ID: 106, ResourceReq: 10, Priority: 1, TaskType: 2, Prev: []int{}},
	}
	K := 2 // 簇数量

	// -------- 无人机聚类 --------
	// 构建无人机特征向量（这里只使用资源）
	dronePoints := [][]float64{}
	nextMap := make(map[int][]int)
	for _, d := range drones {
		dronePoints = append(dronePoints, []float64{d.Resource})
		nextMap[d.ID] = d.NextUavs
	}
	// 执行 K-means 聚类
	droneLabels := kMeansCluster(dronePoints, K, 100)
	// 分组结果
	droneClusters := make([][]int, K)
	for i, label := range droneLabels {
		droneClusters[label] = append(droneClusters[label], drones[i].ID)
	}
	// 簇内连通性检查
	for k := 0; k < K; k++ {
		if !isClusterConnected(droneClusters[k], nextMap) {
			fmt.Printf("警告: 无人机簇 %d 内部不连通\n", k)
		}
	}

	// -------- 任务聚类 --------
	// 构建任务特征向量（资源需求, 优先级, 类型）
	taskPoints := [][]float64{}
	for _, t := range tasks {
		taskPoints = append(taskPoints, []float64{t.ResourceReq, float64(t.Priority), float64(t.TaskType)})
	}
	taskLabels := kMeansCluster(taskPoints, K, 100)
	taskClusters := make([][]int, K)
	for i, label := range taskLabels {
		taskClusters[label] = append(taskClusters[label], tasks[i].ID)
	}
	// 合并依赖组件（确保Prev相关任务同簇）
	comps := buildTaskComponents(tasks)
	for _, comp := range comps {
		if len(comp) > 1 {
			mainCluster := -1
			// 找到组件主簇
			for _, tid := range comp {
				for k := 0; k < K; k++ {
					if contains(taskClusters[k], tid) {
						mainCluster = k
						break
					}
				}
				if mainCluster >= 0 {
					break
				}
			}
			// 将其他簇任务移动到主簇
			if mainCluster >= 0 {
				for _, tid := range comp {
					if !contains(taskClusters[mainCluster], tid) {
						for k := 0; k < K; k++ {
							if k != mainCluster {
								// 从簇 k 中移除 tid
								newList := []int{}
								for _, x := range taskClusters[k] {
									if x != tid {
										newList = append(newList, x)
									}
								}
								taskClusters[k] = newList
							}
						}
						taskClusters[mainCluster] = append(taskClusters[mainCluster], tid)
					}
				}
			}
		}
	}

	// -------- 输出聚类结果 --------
	fmt.Println("无人机分组结果:")
	for k := 0; k < K; k++ {
		fmt.Printf("Group %d: %v\n", k, droneClusters[k])
	}
	fmt.Println("任务分组结果:")
	for k := 0; k < K; k++ {
		fmt.Printf("Group %d: %v\n", k, taskClusters[k])
	}

	// -------- KM 匹配接口调用 --------
	for k1 := 0; k1 < K; k1++ {
		for k2 := 0; k2 < K; k2++ {
			matchDroneTask(droneClusters[k1], taskClusters[k2])
		}
	}
}
