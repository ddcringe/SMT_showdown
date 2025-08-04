package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
type TokenUtil struct {
	secret []byte
	ttl    time.Duration
}

// NewTokenUtil создает новый экземпляр JWT утилиты
func NewTokenUtil(secret string, ttl time.Duration) *TokenUtil {
	return &TokenUtil{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

// GenerateToken создает новый JWT токен
func (t *TokenUtil) GenerateToken(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "smt-showdown-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.secret)
}

// ParseToken валидирует и парсит JWT токен
func (t *TokenUtil) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return t.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
