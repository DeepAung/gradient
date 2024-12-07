package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/DeepAung/gradient/grader-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String(
		"addr",
		"localhost:50051",
		"The server address in the format of host:port",
	)
	serverHostOverride = flag.String(
		"server_host_override",
		"x.test.example.com",
		"The server name used to verify the hostname returned by the TLS handshake",
	)
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewGraderClient(conn)

	submitCode(client, "example.cpp", proto.LanguageType_CPP)
	// submitCode(client, "example.c", proto.LanguageType_C)
	// submitCode(client, "example.go", proto.LanguageType_GO)
	// submitCode(client, "example.py", proto.LanguageType_PYTHON)
}

func submitCode(client proto.GraderClient, codeFilename string, language proto.LanguageType) {
	file, err := os.Open("examples/" + codeFilename)
	if err != nil {
		log.Fatal("os.Open: ", err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("io.ReadAll: ", err)
	}
	code := string(b)

	stream, err := client.Grade(context.Background(), &proto.Input{
		Code:         code,
		CodeFilename: codeFilename,
		Language:     language,
		TaskId:       1,
	})
	if err != nil {
		log.Fatal("client.Grade: ", err)
	}

	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("stream.Recv: ", err)
		}

		log.Println("result: ", proto.ResultType_name[int32(result.Result)])
	}
}
