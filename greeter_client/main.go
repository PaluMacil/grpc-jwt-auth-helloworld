package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

type TokenSource struct {
	oauth2.TokenSource
}

func (ts *TokenSource) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"authorization": "Bearer " + token.AccessToken,
	}, nil
}

func (ts *TokenSource) RequireTransportSecurity() bool {
	return true
}

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
		TokenType:   "Bearer",
	}
	creds := &TokenSource{
		oauth2.StaticTokenSource(token),
	}

	// Create a certificate
	cert, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(cert) {
		log.Fatalf("failed to append certificates")
	}
	config := &tls.Config{
		InsecureSkipVerify: true, // do not use in production!
		RootCAs:            cp,
	}
	connCreds := credentials.NewTLS(config)

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(connCreds), grpc.WithPerRPCCredentials(creds))
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
