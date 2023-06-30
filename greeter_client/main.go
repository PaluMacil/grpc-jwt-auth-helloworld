package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
	tokenPath   = "token.jwt"
	certFile    = "cert.pem"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		log.Fatalf("could not read %s: %s", tokenPath, err)
	}
	tokenString := string(tokenBytes)

	// Create a credentials object with the JWT
	token := &oauth2.Token{
		AccessToken: tokenString,
	}
	creds := oauth.NewOauthAccess(token)

	// Create a certificate
	cert, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(cert), grpc.WithPerRPCCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
