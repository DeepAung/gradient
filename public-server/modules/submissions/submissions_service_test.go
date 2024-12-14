package submissions

import (
	"fmt"
	"log"
	"testing"

	"github.com/DeepAung/gradient/grader-server/proto"
	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/database"
	"github.com/DeepAung/gradient/public-server/modules/tasks"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/modules/users"
	"github.com/DeepAung/gradient/public-server/pkg/asserts"
	"github.com/DeepAung/gradient/public-server/pkg/graderclient"
	"github.com/jmoiron/sqlx"
)

var (
	migrateSourceName = "../../migrations/migrate.sql"
	seedSourceName    = "../../migrations/seed.sql"
	cfg               *config.Config
	db                *sqlx.DB
	client            proto.GraderClient
	submissionsRepo   types.SubmissionsRepo
	svc               types.SubmissionsSvc

	createReq = types.CreateSubmissionReq{
		UserId:   1,
		TaskId:   1,
		Code:     "for _ in range(len(int(input()))): print(int(input()) + int(input()))",
		Language: proto.LanguageType_PYTHON,
	}

	listSubmissions = []types.Submission{
		{
			Id:            9,
			UserId:        2,
			TaskId:        1,
			Code:          "print(123456)",
			Language:      "python",
			Results:       "----------",
			ResultPercent: 0,
		},
		{
			Id:            8,
			UserId:        1,
			TaskId:        3,
			Code:          "print(123456)",
			Language:      "python",
			Results:       "-",
			ResultPercent: 0,
		},
		{
			Id:            7,
			UserId:        1,
			TaskId:        2,
			Code:          "print(123456)",
			Language:      "python",
			Results:       "----------",
			ResultPercent: 0,
		},
		{
			Id:            6,
			UserId:        1,
			TaskId:        2,
			Code:          "for _ in range(len(int(input()))): print(int(input()) + int(input()))",
			Language:      "python",
			Results:       "PPPPPPPPPP",
			ResultPercent: 100,
		},
		{
			Id:            5,
			UserId:        1,
			TaskId:        1,
			Code:          "println(123456)",
			Language:      "go",
			Results:       "----------",
			ResultPercent: 0,
		},
		{
			Id:            4,
			UserId:        1,
			TaskId:        1,
			Code:          "print(123456)",
			Language:      "python",
			Results:       "----------",
			ResultPercent: 0,
		},
		{
			Id:            3,
			UserId:        1,
			TaskId:        1,
			Code:          "for _ in range(len(int(input()))): print(int(input()) + int(input()))",
			Language:      "python",
			Results:       "PPPPPPPPPP",
			ResultPercent: 100,
		},
		{
			Id:            2,
			UserId:        1,
			TaskId:        1,
			Code:          "for _ in range(len(int(input()))): print(int(input()) + int(input()))",
			Language:      "python",
			Results:       "PPPPPPPPPP",
			ResultPercent: 100,
		},
		{
			Id:            1,
			UserId:        1,
			TaskId:        1,
			Code:          "for _ in range(len(int(input()))): print(int(input()) + int(input()))",
			Language:      "python",
			Results:       "PPPPPPPPPP",
			ResultPercent: 100,
		},
	}
)

func init() {
	cfg = config.NewConfig("../../.env.dev")
	db = database.InitDB(cfg.App.DbUrl)
	database.RunSQL(db, migrateSourceName)
	database.RunSQL(db, seedSourceName)
	submissionsRepo = NewSubmissionRepo(db, cfg.App.Timeout)
	tasksRepo := tasks.NewTasksRepo(db, cfg.App.Timeout)
	client = graderclient.NewGraderClientMock(10, 0)
	svc = NewSubmissionSvc(submissionsRepo, tasksRepo, client)
}

func TestGetSubmission(t *testing.T) {
	t.Run("id not found", func(t *testing.T) {
		submission, err := svc.GetSubmission(1000)
		asserts.EqualError(t, err, ErrSubmissionNotFound)
		asserts.Equal(t, "submission", submission, types.Submission{})
	})

	t.Run("normal get submission", func(t *testing.T) {
		expectedSubmission := listSubmissions[len(listSubmissions)-1]
		submission, err := svc.GetSubmission(expectedSubmission.Id)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "submission", submission, expectedSubmission)
	})
}

func TestGetSubmissions(t *testing.T) {
	t.Run("normal get submissions", func(t *testing.T) {
		submissions, err := svc.GetSubmissions(types.GetSubmissionsReq{UserId: nil, TaskId: nil})
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "submissions", submissions, listSubmissions)
	})

	t.Run("with user id = 1 (DeepAung)", func(t *testing.T) {
		userId := 1
		submissions, err := svc.GetSubmissions(
			types.GetSubmissionsReq{UserId: &userId, TaskId: nil},
		)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "submissions", submissions, listSubmissions[1:])
	})

	t.Run("with user id = 2 (admin)", func(t *testing.T) {
		userId := 2
		submissions, err := svc.GetSubmissions(
			types.GetSubmissionsReq{UserId: &userId, TaskId: nil},
		)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "submissions", submissions, listSubmissions[0:1])
	})

	t.Run("with task id = 1", func(t *testing.T) {
		taskId := 1
		submissions, err := svc.GetSubmissions(
			types.GetSubmissionsReq{UserId: nil, TaskId: &taskId},
		)
		expectedSubmissions := make([]types.Submission, 6)
		expectedSubmissions[0] = listSubmissions[0]
		copy(expectedSubmissions[1:], listSubmissions[4:])
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "submissions", submissions, expectedSubmissions)
	})

	t.Run("with task id = 2", func(t *testing.T) {
		taskId := 2
		submissions, err := svc.GetSubmissions(
			types.GetSubmissionsReq{UserId: nil, TaskId: &taskId},
		)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "submissions", submissions, listSubmissions[2:4])
	})

	t.Run("with task id = 3", func(t *testing.T) {
		taskId := 3
		submissions, err := svc.GetSubmissions(
			types.GetSubmissionsReq{UserId: nil, TaskId: &taskId},
		)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "submissions", submissions, listSubmissions[1:2])
	})

	t.Run("non-exist task id", func(t *testing.T) {
		taskId := 1000
		submissions, err := svc.GetSubmissions(
			types.GetSubmissionsReq{UserId: nil, TaskId: &taskId},
		)
		asserts.EqualError(t, err, nil)
		if len(submissions) != 0 {
			log.Fatalf("invalid submissions len, expect=0, got=%d", len(submissions))
		}
	})
}

func TestSubmitCodeMockGrader(t *testing.T) {
	t.Run("user id not found", func(t *testing.T) {
		req := createReq
		req.UserId = 1000
		_, _, err := svc.SubmitCode(req)
		asserts.EqualError(t, err, users.ErrUserNotFound)
	})

	t.Run("task is not found", func(t *testing.T) {
		req := createReq
		req.TaskId = 1000
		_, _, err := svc.SubmitCode(req)
		asserts.EqualError(t, err, tasks.ErrTaskNotFound)
	})

	t.Run("invalid language", func(t *testing.T) {
		req := createReq
		req.Language = 1000
		_, _, err := svc.SubmitCode(req)
		asserts.EqualError(t, err, ErrInvalidLanguage)
	})

	t.Run("normal submit code", func(t *testing.T) {
		req := createReq
		resultCh, createCh, err := svc.SubmitCode(req)

		resultLen := 0
		for result := range resultCh {
			fmt.Println("result: ", result)
			resultLen++
		}
		asserts.NotEqual(t, "result length", resultLen, 0)

		createRes := <-createCh
		id, err := createRes.Id, createRes.Err
		asserts.EqualError(t, err, nil)

		submission, err := svc.GetSubmission(id)
		asserts.EqualError(t, err, nil)

		asserts.Equal(t, "submission id", submission.Id, id)
		asserts.Equal(t, "submission user id", submission.UserId, req.UserId)
		asserts.Equal(t, "submission task id", submission.TaskId, req.TaskId)
		asserts.Equal(t, "submission code", submission.Code, req.Code)

		language, ok := types.ProtoLanguageToString(req.Language)
		asserts.Equal(t, "ok", ok, true)
		asserts.Equal(t, "submission language", submission.Language, language)
		asserts.Equal(t, "submission result length", len(submission.Results), resultLen)
	})
}

// // TODO: when grader server finished, add these tests
// func TestSubmitCode(t *testing.T) {
// 	t.Run("compilation error", func(t *testing.T) {})
// 	t.Run("time limit exceeded", func(t *testing.T) {})
// 	t.Run("memory limit exceeded", func(t *testing.T) {})
// 	t.Run("partial invalid", func(t *testing.T) {})
// 	t.Run("all pass", func(t *testing.T) {})
// }
//
