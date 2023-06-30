package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"os"
	"time"

	pb "github.com/PaluMacil/grpc-jwt-auth-helloworld/helloworld"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	anonymousClient pb.TokenClient
	refreshToken    string
}

func (ts *TokenSource) Token() (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := ts.anonymousClient.Refresh(ctx, &pb.RefreshRequest{RefreshToken: ts.refreshToken})
	if err != nil {
		return nil, err
	}
	tokenString := r.GetAccessToken()
	// parse refresh_token from new JWT
	parser := &jwt.Parser{}
	claims := jwt.MapClaims{}
	_, _, err = parser.ParseUnverified(tokenString, &claims)
	if err != nil {
		return nil, err
	}
	ts.refreshToken = claims["refresh_token"].(string)
	return &oauth2.Token{
		AccessToken: tokenString,
		TokenType:   "Bearer",
	}, nil
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

	// parse refresh_token from JWT
	parser := &jwt.Parser{}
	claims := jwt.MapClaims{}
	_, _, err = parser.ParseUnverified(tokenString, &claims)
	if err != nil {
		log.Fatalf("failed to parse token: %v", err)
	}
	refreshToken := claims["refresh_token"].(string)

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

	// Set up a connection to the anonymous server without authentication.
	anonymousConn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(connCreds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer anonymousConn.Close()
	anonymousClient := pb.NewTokenClient(anonymousConn)

	// Create a TokenSource that uses RefreshToken method to get new tokens when needed.
	creds := &TokenSource{
		anonymousClient: anonymousClient,
		refreshToken:    refreshToken,
	}

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
