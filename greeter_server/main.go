// Package main implements a server for Greeter service.
package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
)

const (
	port                = 50051
	expectedAccessToken = "EgTm8aUjGpW4bp7ChI1f2zm5muoShF+QkNHna3IVEQY="
	certFile            = "cert.pem"
	keyFile             = "key.pem"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
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

func valid(md metadata.MD) bool {
	tokenString := md.Get("authorization")[0]
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	signingKey, _ := base64.StdEncoding.DecodeString("L7joifscCNr/gr9QEvcD86lp5VO0PPx2IDDRBo5CetA=")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		log.Printf("Failed to parse token: %v", err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accessToken, _ := claims["access_token"].(string)
		return accessToken == expectedAccessToken
	}

	return false
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Load the certificates from disk
	certificate, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("could not load server key pair: %s", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer(grpc.Creds(certificate))
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
