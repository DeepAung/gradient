package server_test

import (
	"context"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DeepAung/gradient/grader-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = "localhost:50051"
	client     proto.GraderClient
)

func init() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client = proto.NewGraderClient(conn)
}

func TestGrade(t *testing.T) {
	submitCode(t, "../examples/code.cpp", proto.LanguageType_CPP)
	submitCode(t, "../examples/code.c", proto.LanguageType_C)
	submitCode(t, "../examples/code.go", proto.LanguageType_GO)
	submitCode(t, "../examples/code.py", proto.LanguageType_PYTHON)
}

func submitCode(t *testing.T, codeFilename string, language proto.LanguageType) {
	file, err := os.Open(codeFilename)
	if err != nil {
		t.Fatal("os.Open: ", err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		t.Fatal("io.ReadAll: ", err)
	}
	code := string(b)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.Grade(ctx, &proto.Input{
		Code:        code,
		Language:    language,
		TaskId:      1,
		TimeLimit:   1000,
		MemoryLimit: 1000,
	})
	if err != nil {
		t.Fatal("client.Grade: ", err)
	}

	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal("stream.Recv: ", err)
		}

		log.Println("result")
		log.Println(" - status: ", proto.StatusType_name[int32(result.Status)])
		log.Println(" - time: ", result.Time)
		log.Println(" - memory: ", result.Memory)
	}
}
