package tool

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var hmacSampleSecret = []byte("selbasehmacSampleSecret")

// NewJWT 创建一个新的jwt
func NewJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": strconv.FormatInt(time.Now().Unix(), 10),
	})
	return token.SignedString(hmacSampleSecret)
}

// JWTVal jwt验证
func JWTVal(tokenString string) bool {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		v, ok := claims["nbf"].(string)
		if !ok {
			return false
		}
		clainmTimeStr, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return false
		}
		clainmTime := time.Unix(clainmTimeStr, 0)
		if time.Now().Sub(clainmTime) > time.Minute*3 {
			return false
		}
		return true
	}
	return false
}
