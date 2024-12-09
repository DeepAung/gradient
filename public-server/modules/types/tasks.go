package types

type TaskSvc interface {
	GetTask(id int) (Task, error)
	GetTasks(search string, onlyCompleted bool, startIndex, stopIndex int) ([]Task, error)
	CreateTask(req CreateTaskReq) error
	UpdateTask(req Task) error
	DeleteTask(id int) error
}

type TaskRepo interface {
	FindOneTask(id int) (Task, error)
	// [startIndex, stopIndex)
	FindManyTasks(search string, onlyCompleted bool, startIndex, stopIndex int) ([]Task, error)
	CreateTask(req CreateTaskReq) error
	UpdateTask(req Task) error
	DeleteTask(id int) error
}

type Task struct {
	Id            int    `db:"id"`
	DisplayName   string `db:"display_name"`
	UrlName       string `db:"url_name"`
	ContentUrl    string `db:"content_url"`
	TestcaseCount int    `db:"testcase_count"`
}

type CreateTaskReq struct {
	DisplayName   string `db:"display_name"`
	UrlName       string `db:"url_name"`
	ContentUrl    string `db:"content_url"`
	TestcaseCount int    `db:"testcase_count"`
}
