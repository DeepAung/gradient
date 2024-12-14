package submissions

import (
	"context"
	"fmt"
	"io"

	"github.com/DeepAung/gradient/grader-server/proto"
	"github.com/DeepAung/gradient/public-server/modules/types"
)

type submissionSvc struct {
	submissionsRepo types.SubmissionsRepo
	tasksRepo       types.TasksRepo
	graderClient    proto.GraderClient
}

func NewSubmissionSvc(
	submissionsRepo types.SubmissionsRepo,
	tasksRepo types.TasksRepo,
	graderClient proto.GraderClient,
) types.SubmissionsSvc {
	return &submissionSvc{
		submissionsRepo: submissionsRepo,
		tasksRepo:       tasksRepo,
		graderClient:    graderClient,
	}
}

// 1. try create submission in transaction that rollbacks (handle errors)
// 2. grade code (return results by result channel)
// 3. create submission (return by create channel)
func (s *submissionSvc) SubmitCode(
	req types.CreateSubmissionReq,
) (<-chan proto.ResultType, <-chan types.CreateSubmissionRes, error) {
	if err := s.submissionsRepo.CanCreateSubmission(req); err != nil {
		return nil, nil, err
	}

	stream, err := s.graderClient.Grade(context.Background(), &proto.Input{
		Code:     req.Code,
		Language: req.Language,
		TaskId:   uint32(req.TaskId),
	})
	if err != nil {
		fmt.Println("gcp error: ", err.Error())
		return nil, nil, err
	}

	testcaseCount, err := s.tasksRepo.FindOneTaskTestcaseCount(req.TaskId)
	if err != nil {
		return nil, nil, err
	}

	resultCh := make(chan proto.ResultType, testcaseCount)
	createCh := make(chan types.CreateSubmissionRes)
	go func() {
		req.Results = ""
		passCount, totalCount := 0, 0
		for {
			result, err := stream.Recv()

			var resVar proto.ResultType
			if err == io.EOF {
				break
			} else if err != nil {
				resVar = proto.ResultType_COMPILATION_ERROR
				break
			} else {
				resVar = result.Result
			}

			resultCh <- resVar
			char, _ := types.ProtoResultToChar(resVar) // TODO: handle error
			req.Results += char

			totalCount++
			if resVar == proto.ResultType_PASS {
				passCount++
			}
		}
		close(resultCh)

		req.ResultPercent = float32(passCount) / float32(totalCount)

		id, err := s.submissionsRepo.CreateSubmission(req)
		createCh <- types.CreateSubmissionRes{Id: id, Err: err}
		close(createCh)
	}()

	return resultCh, createCh, nil
}

func (s *submissionSvc) GetSubmission(id int) (types.Submission, error) {
	return s.submissionsRepo.FindOneSubmission(id)
}

func (s *submissionSvc) GetSubmissions(req types.GetSubmissionsReq) ([]types.Submission, error) {
	return s.submissionsRepo.FindManySubmissions(req)
}
