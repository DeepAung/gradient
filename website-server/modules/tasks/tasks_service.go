package tasks

import "github.com/DeepAung/gradient/website-server/modules/types"

type tasksSvcImpl struct {
	repo types.TasksRepo
}

func NewTasksSvc(repo types.TasksRepo) types.TasksSvc {
	return &tasksSvcImpl{
		repo: repo,
	}
}

func (s *tasksSvcImpl) GetTask(userId, taskId int) (types.Task, error) {
	return s.repo.FindOneTask(userId, taskId)
}

func (s *tasksSvcImpl) GetTasks(
	userId int,
	search string,
	onlyCompleted bool,
	startIndex, stopIndex int,
) ([]types.Task, error) {
	return s.repo.FindManyTasks(userId, search, onlyCompleted, startIndex, stopIndex)
}

func (s *tasksSvcImpl) CreateTask(req types.CreateUpdateTaskReq) error {
	return s.repo.CreateTask(req)
}

func (s *tasksSvcImpl) UpdateTask(id int, req types.CreateUpdateTaskReq) error {
	return s.repo.UpdateTask(id, req)
}

func (s *tasksSvcImpl) DeleteTask(id int) error {
	return s.repo.DeleteTask(id)
}
