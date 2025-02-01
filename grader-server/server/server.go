package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

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
	cfg *graderconfig.Config
}

func NewGraderServer(cfg *graderconfig.Config) *graderServer {
	return &graderServer{
		cfg: cfg,
	}
}

func (s *graderServer) Grade(
	input *proto.Input,
	stream grpc.ServerStreamingServer[proto.Result],
) error {
	// Pull testcases from taskId
	log.Println("Grader: start Grade() function")
	testcasesDir := fmt.Sprintf("tmp/testcases/%d", input.TaskId)
	testcasePuller := testcasepuller.NewMockTestcasePuller()
	n, err := testcasePuller.Pull(int(input.TaskId), testcasesDir)
	if err != nil {
		log.Println("Grader: err testcase Pull: ", err.Error())
		return err
	}
	log.Println("Grader: testcase pulled")

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
		log.Println("Grader: err os.MkdirAll", err.Error())
		return err
	}
	codeFile, err := os.Create(codeFilename)
	if err != nil {
		log.Println("Grader: err os.Create: ", err.Error())
		return err
	}
	codeFile.Write([]byte(input.Code))
	codeFile.Close()
	log.Println("Grader: create/write file")

	// Init runner & checker
	runner := runner.NewCodeRunner(s.cfg)
	checker := checker.NewCodeChecker()
	log.Println("Grader: init runner and checker")

	// Build code
	ctx := context.Background() // TODO: change to context with timeout
	ok, result := runner.Build(ctx, input.Language, codeFilename)
	if !ok {
		stream.Send(&proto.Result{Status: result})
		return nil
	}
	log.Println("Grader: builded")

	// Run and Check code
	for i := 1; i <= n; i++ {
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
	log.Println("Grader: runned")

	// Delete code file/folder
	if err := os.RemoveAll(submissionDir); err != nil {
		return err
	}
	log.Println("Grader: cleaned")

	return nil
}
