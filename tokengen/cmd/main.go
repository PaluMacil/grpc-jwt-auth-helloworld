package main

import (
	"flag"
	"fmt"
	"github.com/PaluMacil/grpc-jwt-auth-helloworld/tokengen"
	"log"
	"os"
	"time"
)

const (
	tokenFile = "token.jwt"
)

var (
	already_expired = flag.Bool("expired", false, "whether to create an expired token")
)

func main() {
	flag.Parse()
	var expires time.Time
	if *already_expired {
		log.Printf("generating new expired token")
		expires = tokengen.OneDayAgo()
	} else {
		log.Printf("generating new valid token")
		expires = tokengen.OneDayFromNow()
	}
	tokenString := tokengen.GenerateToken(expires)

	err := os.WriteFile(tokenFile, []byte(tokenString), 0666)
	if err != nil {
		log.Printf("writing token to %s failed: %v", tokenFile, err)
	} else {
		log.Printf("wrote new %s", tokenFile)
	}
	fmt.Println(tokenString)
}
