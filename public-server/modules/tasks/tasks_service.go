package tasks

import "github.com/DeepAung/gradient/public-server/modules/types"

type TasksSvc struct {
	repo types.TasksRepo
}

func NewTasksSvc(repo types.TasksRepo) types.TasksSvc {
	return &TasksSvc{
		repo: repo,
	}
}

func (s *TasksSvc) GetTask(id int) (types.Task, error) {
	return s.repo.FindOneTask(id)
}

func (s *TasksSvc) GetTasks(
	userId int,
	search string,
	onlyCompleted bool,
	startIndex, stopIndex int,
) ([]types.Task, error) {
	return s.repo.FindManyTasks(userId, search, onlyCompleted, startIndex, stopIndex)
}

func (s *TasksSvc) CreateTask(req types.CreateTaskReq) error {
	return s.repo.CreateTask(req)
}

func (s *TasksSvc) UpdateTask(req types.Task) error {
	return s.repo.UpdateTask(req)
}

func (s *TasksSvc) DeleteTask(id int) error {
	return s.repo.DeleteTask(id)
}
