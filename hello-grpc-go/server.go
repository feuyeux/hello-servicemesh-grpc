/*
 * grpc demo
 */

package main

import (
	"context"
	"github.com/feuyeux/hello-grpc-go/common/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	port = ":9996"
)

type protoServer struct{}

func (s *protoServer) Talk(ctx context.Context, in *pb.TalkRequest) (*pb.TalkResponse, error) {
	result := new(pb.TalkResult)
	result.Id = time.Now().UnixNano()
	result.Type = pb.ResultType_OK
	kv := make(map[string]string)
	kv["data"] = in.Data
	kv["meta"] = in.Meta
	kv["timestamp"] = strconv.Itoa(time.Now().Nanosecond())
	result.Kv = kv

	return &pb.TalkResponse{
		Status:  200,
		Results: []*pb.TalkResult{result},
	}, nil
}

func (s *protoServer) TalkOneAnswerMore(in *pb.TalkRequest, stream pb.LandingService_TalkOneAnswerMoreServer) error {
	rs := make([]*pb.TalkResponse, 2)

	kv := make(map[string]string)
	result := new(pb.TalkResult)
	result.Id = time.Now().UnixNano()
	result.Type = pb.ResultType_OK
	kv["data"] = in.Data
	kv["meta"] = in.Meta
	kv["timestamp"] = strconv.Itoa(time.Now().Nanosecond())
	result.Kv = kv
	rs[0] = &pb.TalkResponse{
		Status:  200,
		Results: []*pb.TalkResult{result},
	}

	result2 := new(pb.TalkResult)
	result2.Id = time.Now().UnixNano()
	result2.Type = pb.ResultType_OK
	kv["timestamp"] = strconv.Itoa(time.Now().Nanosecond())
	result2.Kv = kv
	rs[1] = &pb.TalkResponse{
		Status:  200,
		Results: []*pb.TalkResult{result2},
	}

	for _, r := range rs {
		if err := stream.Send(r); err != nil {
			return err
		}
	}
	return nil
}

func (s *protoServer) TalkMoreAnswerOne(stream pb.LandingService_TalkMoreAnswerOneServer) error {
	var rs []*pb.TalkResult
	for {
		in, err := stream.Recv()

		if err == io.EOF {
			talkResponse := &pb.TalkResponse{
				Status:  200,
				Results: rs,
			}
			stream.SendAndClose(talkResponse)
			return nil
		}
		if err != nil {
			return err
		}

		kv := make(map[string]string)
		kv["data"] = in.Data
		kv["meta"] = in.Meta
		kv["timestamp"] = strconv.Itoa(time.Now().Nanosecond())

		result := new(pb.TalkResult)
		result.Id = time.Now().UnixNano()
		result.Type = pb.ResultType_OK
		result.Kv = kv
		rs = append(rs, result)
	}
}

func (s *protoServer) TalkBidirectional(stream pb.LandingService_TalkBidirectionalServer) error {
	var rs []*pb.TalkResponse
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Print(rs)
			return nil
		}
		if err != nil {
			return err
		}

		kv := make(map[string]string)
		kv["data"] = in.Data
		kv["meta"] = in.Meta
		kv["timestamp"] = strconv.Itoa(time.Now().Nanosecond())

		result := new(pb.TalkResult)
		result.Id = time.Now().UnixNano()
		result.Type = pb.ResultType_OK
		result.Kv = kv

		talkResponse := &pb.TalkResponse{
			Status:  200,
			Results: []*pb.TalkResult{result},
		}
		if err := stream.Send(talkResponse); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLandingServiceServer(s, &protoServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
