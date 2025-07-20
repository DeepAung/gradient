package submissions

import (
	"errors"

	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/grader-server/proto"
	"github.com/DeepAung/gradient/website-server/modules/types"
)

type submissionSvcImpl struct {
	submissionsRepo types.SubmissionsRepo
	tasksRepo       types.TasksRepo
	graderClient    proto.GraderClient
	graderCfg       *graderconfig.Config
}

func NewSubmissionSvc(
	submissionsRepo types.SubmissionsRepo,
	tasksRepo types.TasksRepo,
	graderClient proto.GraderClient,
	graderCfg *graderconfig.Config,
) types.SubmissionsSvc {
	return &submissionSvcImpl{
		submissionsRepo: submissionsRepo,
		tasksRepo:       tasksRepo,
		graderClient:    graderClient,
		graderCfg:       graderCfg,
	}
}

// 1. try create submission in transaction that rollbacks (handle errors)
// 2. grade code (return results by result channel)
// 3. create submission (return by create channel)
func (s *submissionSvcImpl) SubmitCode(
	req types.CreateSubmissionReq,
) (<-chan proto.StatusType, <-chan types.CreateSubmissionRes, error) {
	return nil, nil, errors.New("not implemented yet")
	// if err := s.submissionsRepo.CanCreateSubmission(req); err != nil {
	// 	return nil, nil, err
	// }
	//
	// languageInfo, ok := s.graderCfg.GetLanguageInfoFromProtoIndex(req.LanguageIndex)
	// if !ok {
	// 	return nil, nil, ErrInvalidLanguage
	// }
	//
	// stream, err := s.graderClient.Grade(context.Background(), &proto.Input{
	// 	Code:     req.Code,
	// 	Language: languageInfo.Proto,
	// 	TaskId:   uint32(req.TaskId),
	// })
	// if err != nil {
	// 	fmt.Println("gcp error: ", err.Error())
	// 	return nil, nil, err
	// }
	//
	// testcaseCount, err := s.tasksRepo.FindOneTaskTestcaseCount(req.TaskId)
	// if err != nil {
	// 	return nil, nil, err
	// }
	//
	// resultCh := make(chan proto.StatusType, testcaseCount)
	// createCh := make(chan types.CreateSubmissionRes)
	// go func() {
	// 	req.Evaluations = make([]types.CreateEvaluationReq, 0)
	// 	passCount, totalCount := 0, 0
	// 	for {
	// 		result, err := stream.Recv()
	//
	// 		var resVar proto.StatusType
	// 		if err == io.EOF {
	// 			break
	// 		} else if err != nil {
	// 			resVar = proto.StatusType_COMPILATION_ERROR
	// 			break
	// 		} else {
	// 			resVar = result.Result
	// 		}
	//
	// 		resultCh <- resVar
	// 		resultInfo, _ := s.graderCfg.GetResultInfoFromProto(resVar) // TODO: handle error
	// 		req.Evaluations = append(req.Evaluations, types.CreateEvaluationReq{
	//
	// 		})
	// 		req.Results += resultInfo.Char
	//
	// 		totalCount++
	// 		if resVar == proto.StatusType_PASS {
	// 			passCount++
	// 		}
	// 	}
	// 	close(resultCh)
	//
	// 	req.ResultPercent = float32(passCount) / float32(totalCount)
	//
	// 	id, err := s.submissionsRepo.CreateSubmission(req)
	// 	createCh <- types.CreateSubmissionRes{Id: id, Err: err}
	// 	close(createCh)
	// }()
	//
	// return resultCh, createCh, nil
}

func (s *submissionSvcImpl) GetSubmission(id int) (types.Submission, error) {
	return s.submissionsRepo.FindOneSubmission(id)
}

func (s *submissionSvcImpl) GetSubmissions(
	req types.GetSubmissionsReq,
) ([]types.Submission, error) {
	return s.submissionsRepo.FindManySubmissions(req)
}
