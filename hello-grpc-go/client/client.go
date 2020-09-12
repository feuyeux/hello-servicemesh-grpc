package main

import (
	"github.com/feuyeux/hello-grpc-go/common/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
)

const (
	address = "localhost:9996"
)

func talk(client pb.LandingServiceClient, request *pb.TalkRequest) {
	r, err := client.Talk(context.Background(), request)
	if err != nil {
		log.Fatalf("fail to talk: %v", err)
	}
	log.Printf("talk response: %d %s", r.Status, r.Results)
}
func talkOneAnswerMore(client pb.LandingServiceClient, request *pb.TalkRequest) {
	stream, err := client.TalkOneAnswerMore(context.Background(), request)
	if err != nil {
		log.Fatalf("%v.TalkOneAnswerMore(_) = _, %v", client, err)
	}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.TalkOneAnswerMore(_) = _, %v", client, err)
		}
		log.Printf("talkOneAnswerMore response: %d %s", r.Status, r.Results)
	}
}
func talkMoreAnswerOne(client pb.LandingServiceClient, requests []*pb.TalkRequest) {
	stream, err := client.TalkMoreAnswerOne(context.Background())
	if err != nil {
		log.Fatalf("%v.TalkMoreAnswerOne(_) = _, %v", client, err)
	}
	for _, request := range requests {
		if err := stream.Send(request); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, request, err)
		}
	}
	r, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.TalkMoreAnswerOne() got error %v, want %v", stream, err, nil)
	}
	log.Printf("talkMoreAnswerOne response: %d %s", r.Status, r.Results)
}

func talkBidirectional(client pb.LandingServiceClient, requests []*pb.TalkRequest) {
	stream, err := client.TalkBidirectional(context.Background())
	if err != nil {
		log.Fatalf("%v.TalkBidirectional(_) = _, %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			r, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("talkBidirectional response: %d %s", r.Status, r.Results)
		}
	}()
	for _, request := range requests {
		if err := stream.Send(request); err != nil {
			log.Fatalf("Failed to send : %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewLandingServiceClient(conn)
	talk(c, &pb.TalkRequest{Data: "query=ai,from=0,size=1000,order=x,sort=y", Meta: "user=eric"})
	talkOneAnswerMore(c, &pb.TalkRequest{Data: "query=ai,from=0,size=1000,order=x,sort=y", Meta: "user=eric"})
	talkBidirectional(c, []*pb.TalkRequest{&pb.TalkRequest{Data: "query=ai", Meta: "user=eric"}, &pb.TalkRequest{Data: "query=dialog", Meta: "user=eric"}})
	talkMoreAnswerOne(c, []*pb.TalkRequest{&pb.TalkRequest{Data: "query=ai", Meta: "user=eric"}, &pb.TalkRequest{Data: "query=dialog", Meta: "user=eric"}})
}
