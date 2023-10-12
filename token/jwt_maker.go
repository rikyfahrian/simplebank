package token

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, errors.New("invalid secret key or to short secretKey")
	}

	return &JWTMaker{secretKey}, nil
}

func (j *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", payload, err
	}
	return token, payload, nil

}

func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	tokenJWT, err := jwt.ParseWithClaims(token, &Payload{}, func(t *jwt.Token) (interface{}, error) {

		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	payload, ok := tokenJWT.Claims.(*Payload)
	if !ok {
		return nil, err
	}

	return payload, nil

}
