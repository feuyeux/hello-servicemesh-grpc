/*
 * grpc demo
 */

package main

import (
	"github.com/feuyeux/hello-grpc-go/common/pb"
	"github.com/feuyeux/hello-grpc-go/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":9996"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLandingServiceServer(s, &server.ProtoServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
