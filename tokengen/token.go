package tokengen

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

const (
	refreshToken      = "LlXuNq0L2MvV651ybbsVqdlAxY2bmyefY1O0xMirqIw="
	encodedSigningKey = "L7joifscCNr/gr9QEvcD86lp5VO0PPx2IDDRBo5CetA="
)

type CustomClaims struct {
	RefreshToken string `json:"refresh_token"`
	jwt.RegisteredClaims
}

func OneDayFromNow() time.Time {
	return time.Now().Add(24 * time.Hour)
}

func OneDayAgo() time.Time {
	return time.Now().Add(-24 * time.Hour)
}

func GenerateToken(expires time.Time) string {
	signingKey, _ := base64.StdEncoding.DecodeString(encodedSigningKey)
	user := "somebody"

	claims := CustomClaims{
		RefreshToken: refreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "the_server",
			Subject:   user,
			Audience:  []string{"the_client"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}
	log.Printf("generated new token for %s expiring %v", user, expires)
	return tokenString
}
