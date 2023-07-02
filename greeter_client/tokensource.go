package main

import (
	"context"
	"github.com/PaluMacil/grpc-jwt-auth-helloworld/pb"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"time"
)

type TokenSource struct {
	anonymousClient pb.TokenClient
	accessToken     string
	refreshToken    string
	expiry          time.Time
}

func (ts *TokenSource) Token() (*oauth2.Token, error) {
	// If the token is not expired, return the existing token.
	if time.Now().Before(ts.expiry) {
		return &oauth2.Token{
			AccessToken: ts.accessToken,
			TokenType:   authorizationBearer,
		}, nil
	}

	// If the token is expired, fetch a new one.
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeoutDuration)
	defer cancel()
	r, err := ts.anonymousClient.Refresh(ctx, &pb.RefreshRequest{RefreshToken: ts.refreshToken})
	if err != nil {
		return nil, err
	}
	tokenString := r.GetAccessToken()

	// parse refresh_token and expiry from new JWT
	parser := &jwt.Parser{}
	claims := jwt.MapClaims{}
	_, _, err = parser.ParseUnverified(tokenString, &claims)
	if err != nil {
		return nil, err
	}
	ts.refreshToken, _ = claims["refresh_token"].(string)
	exp, _ := claims["exp"].(float64) // JWT spec specifies exp is NumericDate (float64)
	ts.expiry = time.Unix(int64(exp), 0)
	return &oauth2.Token{
		AccessToken: tokenString,
		TokenType:   authorizationBearer,
	}, nil
}

func (ts *TokenSource) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"authorization": authorizationBearer + " " + token.AccessToken,
	}, nil
}

func (ts *TokenSource) RequireTransportSecurity() bool {
	return true
}
