package main

import (
	_ "embed"
	"flag"
	"log"
	"net"
	"os/exec"

	"github.com/DeepAung/gradient/grader-server/graderconfig"
	"github.com/DeepAung/gradient/grader-server/pkg/checker"
	"github.com/DeepAung/gradient/grader-server/pkg/runner"
	"github.com/DeepAung/gradient/grader-server/proto"
	"github.com/DeepAung/gradient/grader-server/server"
	"github.com/DeepAung/gradient/website-server/pkg/storer"
	grpc "google.golang.org/grpc"
)

//go:embed graderconfig.json
var graderConfigFile []byte

var (
	address       = flag.String("address", "localhost:50051", "grader server's address")
	gcpBucketName = flag.String("gcp-bucket-name", "gradient-bucket-dev", "GCP bucket name")
	maxGoroutines = 5
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	exec.Command("isolate --init")
	defer exec.Command("isolate --cleanup")

	cfg := graderconfig.NewConfig(graderConfigFile)
	storer := storer.NewGcpStorer(*gcpBucketName, maxGoroutines)
	runner := runner.NewCodeRunner(cfg)
	checker := checker.NewCodeChecker()
	graderServer := server.NewGraderServer(runner, checker, cfg, storer)

	grpcServer := grpc.NewServer()
	proto.RegisterGraderServer(grpcServer, graderServer)

	log.Printf("grader server running on address %s", *address)
	grpcServer.Serve(lis)
}
