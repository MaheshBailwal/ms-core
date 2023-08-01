package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTData struct {
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom_claims"`
}

var secretKey = []byte("TokenSecretKey")

func ValidateToken(tokenstr string) (map[string]interface{}, error) {
	claimsMap := make(map[string]interface{})

	token, _ := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return secretKey, nil
	})

	if token == nil {
		fmt.Println("invalid token")
		return claimsMap, errors.New("Token error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("couldn't parse claims")
		return claimsMap, errors.New("Token error")
	}

	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		fmt.Println("token expired")
		return claimsMap, errors.New("Token error")
	}

	for k := range claims {
		claimsMap[k] = claims[k]
	}

	return claimsMap, nil
}
