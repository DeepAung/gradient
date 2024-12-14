package graderclient

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/DeepAung/gradient/grader-server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type graderClient struct{}

func NewGraderClient(
	address string,
	opts ...grpc.DialOption,
) (proto.GraderClient, io.Closer, error) {
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to dial: %v", err)
	}
	return proto.NewGraderClient(conn), conn, nil
}

// ----------------------------  Mock ---------------------------- //

type graderClientMock struct {
	testcaseCount int
}

func NewGraderClientMock(testcaseCount int) proto.GraderClient {
	return &graderClientMock{
		testcaseCount: testcaseCount,
	}
}

func (g *graderClientMock) SetTestcaseCount(testcaseCount int) {
	g.testcaseCount = testcaseCount
}

func (g *graderClientMock) Grade(
	ctx context.Context,
	in *proto.Input,
	opts ...grpc.CallOption,
) (grpc.ServerStreamingClient[proto.Result], error) {
	return &serverStreamingClientMock{testcaseCount: g.testcaseCount}, nil
}

type serverStreamingClientMock struct {
	testcaseCount int
}

func (s *serverStreamingClientMock) Recv() (*proto.Result, error) {
	if s.testcaseCount <= 0 {
		return &proto.Result{}, io.EOF
	}

	s.testcaseCount--
	time.Sleep(1000 * time.Millisecond)
	return &proto.Result{Result: proto.ResultType(rand.Intn(6))}, nil
}
func (s *serverStreamingClientMock) Header() (metadata.MD, error) { return nil, nil }
func (s *serverStreamingClientMock) Trailer() metadata.MD         { return nil }
func (s *serverStreamingClientMock) CloseSend() error             { return nil }
func (s *serverStreamingClientMock) Context() context.Context     { return nil }
func (s *serverStreamingClientMock) SendMsg(m any) error          { return nil }
func (s *serverStreamingClientMock) RecvMsg(m any) error          { return nil }
