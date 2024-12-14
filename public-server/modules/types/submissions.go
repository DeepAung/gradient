package types

import (
	"github.com/DeepAung/gradient/grader-server/proto"
)

type SubmissionsSvc interface {
	SubmitCode(req CreateSubmissionReq) (<-chan proto.ResultType, <-chan CreateSubmissionRes, error)
	GetSubmission(id int) (Submission, error)
	GetSubmissions(req GetSubmissionsReq) ([]Submission, error)
}

type SubmissionsRepo interface {
	CanCreateSubmission(req CreateSubmissionReq) error
	CreateSubmission(req CreateSubmissionReq) (int, error)
	FindOneSubmission(id int) (Submission, error)
	FindManySubmissions(req GetSubmissionsReq) ([]Submission, error)
}

type Submission struct {
	Id            int     `db:"id"`
	UserId        int     `db:"user_id"`
	TaskId        int     `db:"task_id"`
	Code          string  `db:"code"`
	Language      string  `db:"language"`
	Results       string  `db:"results"`
	ResultPercent float32 `db:"result_percent"`
}

type CreateSubmissionReq struct {
	UserId   int    `validate:"required"`
	TaskId   int    `validate:"required"`
	Code     string `validate:"required"`
	Language proto.LanguageType

	// assign after complete submiting code
	Results       string
	ResultPercent float32
}

type GetSubmissionsReq struct {
	UserId *int
	TaskId *int
}

type CreateSubmissionRes struct {
	Id  int
	Err error
}
