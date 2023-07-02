package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"os"
	"time"

	"github.com/PaluMacil/grpc-jwt-auth-helloworld/pb"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultName            = "world"
	tokenPath              = "token.jwt"
	certFile               = "cert.pem"
	serverAddr             = "localhost:50051"
	contextTimeoutDuration = time.Second
	authorizationBearer    = "Bearer"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	cachedToken := readCachedToken()
	creds, anonymousConn := createTokenSource(cachedToken)
	defer anonymousConn.Close()
	c, conn := createGreeterClient(creds)
	defer conn.Close()
	sendAndReceiveGreetings(c)
}

func readCachedToken() string {
	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		log.Fatalf("could not read %s: %s", tokenPath, err)
	}
	return string(tokenBytes)
}

func parseJWTToken(tokenString string) jwt.MapClaims {
	parser := &jwt.Parser{}
	claims := jwt.MapClaims{}
	_, _, err := parser.ParseUnverified(tokenString, &claims)
	if err != nil {
		log.Fatalf("failed to parse token: %v", err)
	}
	return claims
}

func createTokenSource(cachedTokenString string) (*TokenSource, *grpc.ClientConn) {
	cachedClaims := parseJWTToken(cachedTokenString)
	anonymousClient, anonymousConn := createAnonymousClient()
	refreshToken, _ := cachedClaims["refresh_token"].(string)
	exp, _ := cachedClaims["exp"].(float64)

	return &TokenSource{
		anonymousClient: anonymousClient,
		accessToken:     cachedTokenString,
		refreshToken:    refreshToken,
		expiry:          time.Unix(int64(exp), 0),
	}, anonymousConn
}

func createAnonymousClient() (pb.TokenClient, *grpc.ClientConn) {
	anonymousConn := createAnonymousConnection()
	return pb.NewTokenClient(anonymousConn), anonymousConn
}

func createAnonymousConnection() *grpc.ClientConn {
	connCreds := createConnectionCredentials()
	anonymousConn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(connCreds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return anonymousConn
}

func createConnectionCredentials() credentials.TransportCredentials {
	cert := readCertificate()
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(cert) {
		log.Fatalf("failed to append certificates")
	}
	config := &tls.Config{
		RootCAs: cp,
	}
	return credentials.NewTLS(config)
}

func readCertificate() []byte {
	cert, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}
	return cert
}

func createGreeterClient(creds *TokenSource) (pb.GreeterClient, *grpc.ClientConn) {
	connCreds := createConnectionCredentials()
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(connCreds), grpc.WithPerRPCCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewGreeterClient(conn), conn
}

func sendAndReceiveGreetings(c pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeoutDuration)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
