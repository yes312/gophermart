package jwtpackage

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	TokenExp time.Duration
	Secret   string
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

var ErrInvalidToken = errors.New("invalid token")

func NewToken(tokenExp time.Duration, secret string) *Token {
	return &Token{
		TokenExp: tokenExp,
		Secret:   secret,
	}
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func (t *Token) BuildJWTString(UserID string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.TokenExp)),
		},
		// собственное утверждение
		UserID: UserID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(t.Secret))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func (tok *Token) GetUserId(tokenString string) (string, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(tok.Secret), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	return claims.UserID, err
}
