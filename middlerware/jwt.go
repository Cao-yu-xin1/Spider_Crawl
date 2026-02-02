package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
// @exp: 过期时间，单位秒

const SecretKey = "1234567890"
const Seconds = 86400 * 7

func GetJwtToken(exp int, payload string) (string, error) {
	iat := time.Now().Unix()
	//seconds := 86400 * 7
	claims := make(jwt.MapClaims)
	claims["exp"] = int(iat) + Seconds
	claims["iat"] = iat
	claims["payload"] = payload
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(SecretKey))
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(SecretKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims, err
	} else {
		fmt.Println(err)
	}
	return nil, err
}
