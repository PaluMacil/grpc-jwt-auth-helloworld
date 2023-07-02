package main

import (
	"context"
	"github.com/PaluMacil/grpc-jwt-auth-helloworld/pb"
	"github.com/PaluMacil/grpc-jwt-auth-helloworld/tokengen"
)

type tokenServer struct {
	pb.UnimplementedTokenServer
}

func (s *tokenServer) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshReply, error) {
	newToken := tokengen.GenerateToken(tokengen.OneDayFromNow())
	return &pb.RefreshReply{AccessToken: newToken}, nil
}

type greeterServer struct {
	pb.UnimplementedGreeterServer
}
