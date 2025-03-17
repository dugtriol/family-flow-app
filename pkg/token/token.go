package token

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExp       = time.Hour * 3
	TokenSecretKey = "supersecretkey"
)

type Token interface {
	Create(userID string) (string, error)
	Check(tokenString string) bool
}

func Create(userID string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)), Subject: userID},
	)
	signedString, err := token.SignedString([]byte(TokenSecretKey))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func Check(tokenString string) (*jwt.Token, error) {
	data := &jwt.RegisteredClaims{}
	var token *jwt.Token
	var err error

	if token, err = jwt.ParseWithClaims(
		tokenString, data,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(TokenSecretKey), nil
		},
	); err != nil {
		log.Printf(fmt.Sprintf("Token - Check - jwt.ParseWithClaims: %v", token))
		return nil, err
	}

	return token, nil
}
