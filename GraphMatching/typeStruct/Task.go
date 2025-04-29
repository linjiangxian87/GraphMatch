package typeStruct

type Task struct {
	TaskID   string
	TaskType int
	Priority int
	Prev     string
	Next     string

	NeedResources map[int]int
}
