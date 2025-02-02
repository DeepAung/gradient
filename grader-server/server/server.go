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
	"github.com/DeepAung/gradient/grader-server/proto"
	"github.com/DeepAung/gradient/website-server/pkg/storer"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type graderServer struct {
	proto.UnimplementedGraderServer
	cfg    *graderconfig.Config
	storer storer.Storer
}

func NewGraderServer(cfg *graderconfig.Config, storer storer.Storer) *graderServer {
	return &graderServer{
		cfg:    cfg,
		storer: storer,
	}
}

func (s *graderServer) Grade(
	input *proto.Input,
	stream grpc.ServerStreamingServer[proto.Result],
) error {
	// Pull testcases from taskId
	remoteDir := fmt.Sprintf("testcases/%d", input.TaskId)
	testcasesDir := fmt.Sprintf("tmp/testcases/%d", input.TaskId)
	cnt, err := s.storer.DownloadFolder(remoteDir, testcasesDir)
	if err != nil {
		return err
	}
	testcaseCount := cnt / 2
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ok, status := runner.Build(ctx, input.Language, codeFilename)
	if !ok {
		stream.Send(&proto.Result{Status: status})
		return nil
	}
	// End

	// Run and Check code
	for i := 1; i <= testcaseCount; i++ {
		inputFilename := fmt.Sprintf("%s/%02d.in", testcasesDir, i)

		ok, result := runner.Run(ctx, input.Language, codeFilename, inputFilename)
		if !ok {
			stream.Send(&result)
			continue
		}

		outputFilename := fmt.Sprintf("%s/%02d.out", testcasesDir, i)
		resultFilename := fmt.Sprintf("%s/%02d.result", submissionDir, i)
		ok, err := checker.CheckFile(ctx, outputFilename, resultFilename)
		if err != nil {
			return err
		}
		if ok {
			stream.Send(&proto.Result{
				Status: proto.StatusType_PASS,
				Time:   result.Time,
				Memory: result.Memory,
			})
		} else {
			stream.Send(&proto.Result{
				Status: proto.StatusType_INCORRECT,
				Time:   result.Time,
				Memory: result.Memory,
			})
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
