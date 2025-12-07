package storage

type Task struct {
	ID              int    `json:"id"`
	ResponsibleID   int    `json:"responsible_id"`
	ResponsibleName string `json:"responsible_name"`
	Context         string `json:"context"`
	AssignedAt      int64  `json:"assigned_at"`
	DueDate         int64  `json:"due_date"`
}

type Interface interface {
	Tasks() ([]Task, error)
	AddTask(Task) error
	UpdateTask(Task) error
	DeleteTask(Task) error
}

