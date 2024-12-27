package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/grader-server/pkg/checker"
	"github.com/DeepAung/gradient/grader-server/pkg/runner"
	"github.com/DeepAung/gradient/grader-server/pkg/testcasepuller"
	"github.com/DeepAung/gradient/grader-server/proto"
	"github.com/google/uuid"
	grpc "google.golang.org/grpc"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	jsonPath = flag.String("json", ".env.dev.json", "json path")
)

type graderServer struct {
	proto.UnimplementedGraderServer
	cfg *graderconfig.Config
}

func newGraderServer(cfg *graderconfig.Config) *graderServer {
	return &graderServer{
		cfg: cfg,
	}
}

func (s *graderServer) Grade(
	input *proto.Input,
	stream grpc.ServerStreamingServer[proto.Status],
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
	runner, err := runner.NewCodeRunner(input.Language)
	if err != nil {
		return err
	}
	checker := checker.NewCodeChecker()
	log.Println("Grader: init runner and checker")

	// Build code
	ctx := context.Background() // TODO: change to context with timeout
	ok, result := runner.Build(ctx, codeFilename)
	if !ok {
		stream.Send(&proto.Status{Result: result})
		return nil
	}
	log.Println("Grader: builded")

	// Run and Check code
	for i := 1; i <= n; i++ {
		inputFilename := fmt.Sprintf("%s/%02d.in", testcasesDir, i)

		ok, result := runner.Run(ctx, codeFilename, inputFilename)
		if !ok {
			stream.Send(&proto.Status{Result: result})
			continue
		}

		outputFilename := fmt.Sprintf("%s/%02d.out", testcasesDir, i)
		resultFilename := fmt.Sprintf("%s/%02d.result", submissionDir, i)
		ok, err := checker.CheckFile(ctx, outputFilename, resultFilename)
		if err != nil {
			return err
		}
		if ok {
			stream.Send(&proto.Status{Result: proto.StatusType_PASS})
		} else {
			stream.Send(&proto.Status{Result: proto.StatusType_INCORRECT})
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

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cfg := graderconfig.NewConfig(*jsonPath)

	grpcServer := grpc.NewServer()
	graderServer := newGraderServer(cfg)
	proto.RegisterGraderServer(grpcServer, graderServer)

	log.Printf("grader server running on port %d", *port)
	grpcServer.Serve(lis)
}
