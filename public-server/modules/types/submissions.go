package types

import (
	"time"

	"github.com/DeepAung/gradient/grader-server/proto"
)

type SubmissionsSvc interface {
	SubmitCode(req CreateSubmissionReq) (<-chan proto.StatusType, <-chan CreateSubmissionRes, error)
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
	Id                  int       `db:"id"`
	UserId              int       `db:"user_id"`
	TaskId              int       `db:"task_id"`
	Code                string    `db:"code"`
	LanguageIndex       int       `db:"language_index"`
	Score               float32   `db:"score"`
	MaxTimeMicroSeconds int       `db:"max_time"`
	MaxMemoryKiloBytes  int       `db:"max_memory"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`

	Evaluations  []Evaluation `db:"evaluations"`
	UserUsername string
}

type CreateSubmissionReq struct {
	UserId        int    `validate:"required"`
	TaskId        int    `validate:"required"`
	Code          string `validate:"required"`
	LanguageIndex int

	// assign after complete submiting code
	Score       float32
	Evaluations []CreateEvaluationReq
}

type GetSubmissionsReq struct {
	UserId *int
	TaskId *int
}

type CreateSubmissionRes struct {
	Id  int
	Err error
}

type Evaluation struct {
	Id               int  `db:"id"`
	SubmissionId     int  `db:"submission_id"`
	TimeMicroSeconds int  `db:"time"`
	MemoryKiloBytes  int  `db:"memory"`
	Status           byte `db:"status"`
}

type CreateEvaluationReq struct {
	SubmissionId     int  `db:"submission_id"`
	TimeMicroSeconds int  `db:"time"`
	MemoryKiloBytes  int  `db:"memory"`
	Status           byte `db:"status"`
}
