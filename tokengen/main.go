package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	refreshToken = "LlXuNq0L2MvV651ybbsVqdlAxY2bmyefY1O0xMirqIw="
	accessToken  = "EgTm8aUjGpW4bp7ChI1f2zm5muoShF+QkNHna3IVEQY="
)

type CustomClaims struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	jwt.RegisteredClaims
}

func main() {
	signingKey, _ := base64.StdEncoding.DecodeString("L7joifscCNr/gr9QEvcD86lp5VO0PPx2IDDRBo5CetA=")

	claims := CustomClaims{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "the_server",
			Subject:   "somebody",
			Audience:  []string{"the_client"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	os.WriteFile("token.jwt", []byte(tokenString), 0666)
	fmt.Println(tokenString)
}
