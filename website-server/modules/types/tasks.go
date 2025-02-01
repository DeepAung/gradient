package types

type TasksSvc interface {
	GetTask(userId, taskId int) (Task, error)
	GetTasks(
		userId int,
		search string,
		onlyCompleted bool,
		startIndex, stopIndex int,
	) ([]Task, error)
	CreateTask(req CreateUpdateTaskReq) error
	UpdateTask(id int, req CreateUpdateTaskReq) error
	DeleteTask(id int) error
}

type TasksRepo interface {
	FindOneTask(userId, taskId int) (Task, error)
	FindOneTaskTestcaseCount(taskId int) (int, error) // TODO: write test
	// [startIndex, stopIndex)
	FindManyTasks(
		userId int,
		search string,
		onlyCompleted bool,
		startIndex, stopIndex int,
	) ([]Task, error)
	CreateTask(req CreateUpdateTaskReq) error
	UpdateTask(id int, req CreateUpdateTaskReq) error
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

type CreateUpdateTaskReq struct {
	DisplayName   string `db:"display_name"`
	UrlName       string `db:"url_name"`
	ContentUrl    string `db:"content_url"`
	TestcaseCount int    `db:"testcase_count"`
}
