package server

import (
	"context"
	"github.com/feuyeux/hello-grpc-go/common/pb"
	"github.com/google/uuid"
	"io"
	"strconv"
	"strings"
	"time"
)

var helloList = []string{"Hello", "Bonjour", "Hola", "こんにちは", "Ciao", "안녕하세요"}

//implement LandingServiceServer interface
type ProtoServer struct{}

func (s *ProtoServer) Talk(ctx context.Context, request *pb.TalkRequest) (*pb.TalkResponse, error) {
	result := s.buildResult(request.Data)
	return &pb.TalkResponse{
		Status:  200,
		Results: []*pb.TalkResult{result},
	}, nil
}

func (s *ProtoServer) TalkOneAnswerMore(request *pb.TalkRequest, stream pb.LandingService_TalkOneAnswerMoreServer) error {
	datas := strings.Split(request.Data, ",")
	for _, d := range datas {
		result := s.buildResult(d)
		if err := stream.Send(&pb.TalkResponse{
			Status:  200,
			Results: []*pb.TalkResult{result},
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *ProtoServer) TalkMoreAnswerOne(stream pb.LandingService_TalkMoreAnswerOneServer) error {
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
		result := s.buildResult(in.Data)
		rs = append(rs, result)
	}
}

func (s *ProtoServer) TalkBidirectional(stream pb.LandingService_TalkBidirectionalServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		result := s.buildResult(in.Data)
		talkResponse := &pb.TalkResponse{
			Status:  200,
			Results: []*pb.TalkResult{result},
		}
		if err := stream.Send(talkResponse); err != nil {
			return err
		}
	}
}

func (s *ProtoServer) buildResult(id string) *pb.TalkResult {
	index, _ := strconv.Atoi(id)
	kv := make(map[string]string)
	kv["id"] = uuid.New().String()
	kv["idx"] = id
	kv["data"] = helloList[index]
	kv["meta"] = "GOLANG"
	result := new(pb.TalkResult)
	result.Id = time.Now().UnixNano()
	result.Type = pb.ResultType_OK
	result.Kv = kv
	return result
}
