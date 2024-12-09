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
	UserId   int
	TaskId   int
	Code     string
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

func ProtoResultToChar(res proto.ResultType) string {
	switch res {
	case proto.ResultType_COMPILATION_ERROR:
		return "C"
	case proto.ResultType_PASS:
		return "P"
	case proto.ResultType_INCORRECT:
		return "-"
	case proto.ResultType_RUNTIME_ERROR:
		return "X"
	case proto.ResultType_TIME_LIMIT_EXCEEDED:
		return "T"
	case proto.ResultType_MEMORY_LIMIT_EXCEEDED:
		return "M"
	default:
		return "C"
	}
}

func ProtoLanguageToString(language proto.LanguageType) string {
	switch language {
	case proto.LanguageType_CPP:
		return "cpp"
	case proto.LanguageType_C:
		return "c"
	case proto.LanguageType_GO:
		return "go"
	case proto.LanguageType_PYTHON:
		return "python"
	default:
		return ""
	}
}
