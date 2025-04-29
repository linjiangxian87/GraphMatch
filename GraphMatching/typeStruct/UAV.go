package typeStruct

type Uav struct {
	Uid string
	// 资源id->资源数量
	Resources map[int]int
	/*
		目前暂定的资源：
		battery int
		cpu     int
		memory  int
	*/
	// 可通信无人机id->通信能力
	NextUavs    map[string]int
	LoadedTasks []Task
}
