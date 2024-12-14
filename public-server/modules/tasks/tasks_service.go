package tasks

import "github.com/DeepAung/gradient/public-server/modules/types"

type tasksSvc struct {
	repo types.TasksRepo
}

func NewTasksSvc(repo types.TasksRepo) types.TasksSvc {
	return &tasksSvc{
		repo: repo,
	}
}

func (s *tasksSvc) GetTask(userId, taskId int) (types.Task, error) {
	return s.repo.FindOneTask(userId, taskId)
}

func (s *tasksSvc) GetTasks(
	userId int,
	search string,
	onlyCompleted bool,
	startIndex, stopIndex int,
) ([]types.Task, error) {
	return s.repo.FindManyTasks(userId, search, onlyCompleted, startIndex, stopIndex)
}

func (s *tasksSvc) CreateTask(req types.CreateUpdateTaskReq) error {
	return s.repo.CreateTask(req)
}

func (s *tasksSvc) UpdateTask(id int, req types.CreateUpdateTaskReq) error {
	return s.repo.UpdateTask(id, req)
}

func (s *tasksSvc) DeleteTask(id int) error {
	return s.repo.DeleteTask(id)
}
