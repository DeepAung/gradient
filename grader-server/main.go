package main

// import (
// 	_ "embed"
// 	"flag"
// )
//
// var (
// 	address       = flag.String("address", "localhost:50051", "grader server's address")
// 	gcpBucketName = flag.String("gcp-bucket-name", "gradient-bucket-dev", "GCP bucket name")
// 	maxGoroutines = 5
// )
//
// func main() {
// 	flag.Parse()
// 	lis, err := net.Listen("tcp", *address)
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
//
// 	storer := storer.NewGcpStorer(*gcpBucketName, maxGoroutines)
// 	graderServer := server.NewGraderServer(cfg, storer)
//
// 	grpcServer := grpc.NewServer()
// 	proto.RegisterGraderServer(grpcServer, graderServer)
//
// 	log.Printf("grader server running on address %s", *address)
// 	grpcServer.Serve(lis)
// }
