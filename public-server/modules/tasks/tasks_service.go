package tasks

import "github.com/DeepAung/gradient/public-server/modules/types"

type TaskSvc struct {
	repo types.TaskRepo
}

func NewTaskSvc(repo types.TaskRepo) types.TaskSvc {
	return &TaskSvc{
		repo: repo,
	}
}

func (s *TaskSvc) GetTask(id int) (types.Task, error) {
	return s.repo.FindOneTask(id)
}

func (s *TaskSvc) GetTasks(
	search string,
	onlyCompleted bool,
	startIndex, stopIndex int,
) ([]types.Task, error) {
	return s.repo.FindManyTasks(search, onlyCompleted, startIndex, stopIndex)
}

func (s *TaskSvc) CreateTask(req types.CreateTaskReq) error {
	return s.repo.CreateTask(req)
}

func (s *TaskSvc) UpdateTask(req types.Task) error {
	return s.repo.UpdateTask(req)
}

func (s *TaskSvc) DeleteTask(id int) error {
	return s.repo.DeleteTask(id)
}
