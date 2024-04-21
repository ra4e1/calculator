package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthServie struct {
	hmacSecret string
}

func NewAuthService() *AuthServie { //создание
	return &AuthServie{
		hmacSecret: "xoBfoI0kxUlmWgylmYCdjL86x4Qe27njO4Vjw1NSfNbeD7rBRNTDo8WOQq1Nokni",
	}
}

func (s *AuthServie) CreateToken(user *DbUser) (string, error) { // Создание нового токена
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"name": user.Name,
		"nbf":  now.Unix(),
		"exp":  now.Add(30 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.hmacSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *AuthServie) ValidateToken(tokenString string) (int64, error) {
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.hmacSecret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := tokenFromString.Claims.(jwt.MapClaims)
	if !ok {
		return 0, err
	}

	return int64(claims["id"].(float64)), nil
}
