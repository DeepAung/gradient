package testcasepuller

import "errors"

type TestcasePuller interface {
	Pull(taskId int, directory string) (int, error)
}

type mockTestcasePuller struct{}

func NewMockTestcasePuller() TestcasePuller {
	return &mockTestcasePuller{}
}

func (m *mockTestcasePuller) Pull(taskId int, directory string) (int, error) {
	if taskId == 1 {
		return 3, nil
	}
	return 0, errors.New("Pull: task not found")
}
