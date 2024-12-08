package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/DeepAung/gradient/grader-server/pkg/checker"
	"github.com/DeepAung/gradient/grader-server/pkg/runner"
	"github.com/DeepAung/gradient/grader-server/pkg/testcasepuller"
	"github.com/DeepAung/gradient/grader-server/proto"
	"github.com/google/uuid"
	grpc "google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "The server port")

type graderServer struct {
	proto.UnimplementedGraderServer
}

func newGraderServer() *graderServer {
	return &graderServer{}
}

func (s *graderServer) Grade(
	input *proto.Input,
	stream grpc.ServerStreamingServer[proto.Result],
) error {
	// Pull testcases from taskId
	testcasesDir := fmt.Sprintf("tmp/testcases/%d", input.TaskId)
	testcasePuller := testcasepuller.NewMockTestcasePuller()
	n, err := testcasePuller.Pull(int(input.TaskId), testcasesDir)
	if err != nil {
		return err
	}

	// Create code file/folder
	submissionId := uuid.NewString()
	submissionDir := fmt.Sprintf("tmp/submissions/%s", submissionId)
	codeFilename := fmt.Sprintf("tmp/submissions/%s/%s", submissionId, input.CodeFilename)
	if err := os.MkdirAll(submissionDir, os.ModePerm); err != nil {
		return err
	}
	codeFile, err := os.Create(codeFilename)
	if err != nil {
		return err
	}
	codeFile.Write([]byte(input.Code))
	codeFile.Close()

	// Init runner & checker
	runner, err := runner.NewCodeRunner(input.Language)
	if err != nil {
		return err
	}
	checker := checker.NewCodeChecker()

	// Build code
	ctx := context.Background() // TODO: change to context with timeout
	ok, result := runner.Build(ctx, codeFilename)
	if !ok {
		stream.Send(&proto.Result{Result: result})
		return nil
	}

	// Run and Check code
	for i := 1; i <= n; i++ {
		inputFilename := fmt.Sprintf("%s/%02d.in", testcasesDir, i)

		ok, result := runner.Run(ctx, codeFilename, inputFilename)
		if !ok {
			stream.Send(&proto.Result{Result: result})
			continue
		}

		outputFilename := fmt.Sprintf("%s/%02d.out", testcasesDir, i)
		resultFilename := fmt.Sprintf("%s/%02d.result", submissionDir, i)
		ok, err := checker.Check(ctx, outputFilename, resultFilename)
		if err != nil {
			return err
		}
		if ok {
			stream.Send(&proto.Result{Result: proto.ResultType_PASS})
		} else {
			stream.Send(&proto.Result{Result: proto.ResultType_INCORRECT})
		}
	}

	// Delete code file/folder
	if err := os.RemoveAll(submissionDir); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	graderServer := newGraderServer()
	proto.RegisterGraderServer(grpcServer, graderServer)

	log.Printf("grader server running on port %d", *port)
	grpcServer.Serve(lis)
}
