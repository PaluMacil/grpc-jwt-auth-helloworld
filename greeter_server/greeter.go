package main

import (
	"context"
	"fmt"
	"github.com/PaluMacil/grpc-jwt-auth-helloworld/pb"
	"google.golang.org/grpc/metadata"
	"log"
)

func (s *greeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	if !valid(md) {
		return nil, fmt.Errorf("unauthorized")
	}

	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
