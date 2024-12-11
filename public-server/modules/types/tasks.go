package types

type TasksSvc interface {
	GetTask(id int) (Task, error)
	GetTasks(
		userId int,
		search string,
		onlyCompleted bool,
		startIndex, stopIndex int,
	) ([]Task, error)
	CreateTask(req CreateTaskReq) error
	UpdateTask(req Task) error
	DeleteTask(id int) error
}

type TasksRepo interface {
	FindOneTask(id int) (Task, error)
	// [startIndex, stopIndex)
	FindManyTasks(
		userId int,
		search string,
		onlyCompleted bool,
		startIndex, stopIndex int,
	) ([]Task, error)
	CreateTask(req CreateTaskReq) error
	UpdateTask(req Task) error
	DeleteTask(id int) error
}

type GetTasksDTO struct {
	Search        string `form:"search"`
	OnlyCompleted bool   `form:"only_completed"`
	Page          int    `form:"page"`
}

type Task struct {
	Id            int     `db:"id"`
	DisplayName   string  `db:"display_name"`
	UrlName       string  `db:"url_name"`
	ContentUrl    string  `db:"content_url"`
	TestcaseCount int     `db:"testcase_count"`
	SolvedNumber  int     `db:"solved_number"`
	Score         float32 `db:"score"` // score of user_id, task_id
}

type CreateTaskReq struct {
	DisplayName   string `db:"display_name"`
	UrlName       string `db:"url_name"`
	ContentUrl    string `db:"content_url"`
	TestcaseCount int    `db:"testcase_count"`
}
