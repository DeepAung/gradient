package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/grader-server/pkg/checker"
	"github.com/DeepAung/gradient/grader-server/pkg/runner"
	"github.com/DeepAung/gradient/grader-server/pkg/testcasepuller"
	"github.com/DeepAung/gradient/grader-server/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type graderServer struct {
	proto.UnimplementedGraderServer
	cfg            *graderconfig.Config
	testcasePuller testcasepuller.TestcasePuller
}

func NewGraderServer(
	cfg *graderconfig.Config,
	testcasePuller testcasepuller.TestcasePuller,
) *graderServer {
	return &graderServer{
		cfg:            cfg,
		testcasePuller: testcasePuller,
	}
}

func (s *graderServer) Grade(
	input *proto.Input,
	stream grpc.ServerStreamingServer[proto.Result],
) error {
	// Pull testcases from taskId
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testcasesDir := fmt.Sprintf("tmp/testcases/%d", input.TaskId)
	testcaseNumber, err := s.testcasePuller.Pull(ctx, int(input.TaskId), testcasesDir)
	if err != nil {
		return err
	}
	// End

	// Create code file/folder
	submissionId := uuid.NewString()
	submissionDir := fmt.Sprintf("tmp/submissions/%s", submissionId)
	languageInfo, ok := s.cfg.GetLanguageInfoFromProto(input.Language)
	if !ok {
		return errors.New("invalid language")
	}
	codeExt := languageInfo.Extension

	codeFilename := fmt.Sprintf("tmp/submissions/%s/%s", submissionId, "code"+codeExt)
	if err := os.MkdirAll(submissionDir, os.ModePerm); err != nil {
		return err
	}
	codeFile, err := os.Create(codeFilename)
	if err != nil {
		return err
	}
	codeFile.Write([]byte(input.Code))
	codeFile.Close()
	// End

	runner := runner.NewCodeRunner(s.cfg)
	checker := checker.NewCodeChecker()

	// Build code
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(input.TimeLimit))
	defer cancel()

	ok, result := runner.Build(ctx, input.Language, codeFilename)
	if !ok {
		stream.Send(&proto.Result{Status: result})
		return nil
	}
	// End

	// Run and Check code
	for i := 1; i <= testcaseNumber; i++ {
		inputFilename := fmt.Sprintf("%s/%02d.in", testcasesDir, i)

		ok, result := runner.Run(ctx, input.Language, codeFilename, inputFilename)
		if !ok {
			stream.Send(&proto.Result{Status: result})
			continue
		}

		outputFilename := fmt.Sprintf("%s/%02d.out", testcasesDir, i)
		resultFilename := fmt.Sprintf("%s/%02d.result", submissionDir, i)
		ok, err := checker.CheckFile(ctx, outputFilename, resultFilename)
		if err != nil {
			return err
		}
		if ok {
			stream.Send(&proto.Result{Status: proto.StatusType_PASS})
		} else {
			stream.Send(&proto.Result{Status: proto.StatusType_INCORRECT})
		}
	}
	// End

	// Delete code file/folder
	if err := os.RemoveAll(submissionDir); err != nil {
		return err
	}
	// End

	return nil
}
