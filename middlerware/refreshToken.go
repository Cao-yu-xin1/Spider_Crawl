package middleware

import (
	jwt2 "Spider_Crawl/pkg/jwt"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

//const secretKey = "1234567890"

//RefreshToken 刷新JWT令牌

func RefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwt2.SecretKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//fmt.Println(claims["foo"], claims["nbf"])
		//payload, err := ParseToken(tokenString)
		// 提取原payload
		// 计算令牌剩余过期时间
		exp, ok := claims["exp"].(float64)
		if !ok {
			exp = float64(time.Now().Unix() + 86400*7) // 默认7天
		}
		remainingTime := time.Until(time.Unix(int64(exp), 0))
		if remainingTime <= 0 {
			return "", errors.New("token已过期")
		}
		/*// 从白名单中删除令牌
		wKey := "whitelist" + tokenString
		if err := config.Rdb.Del(config.Ctx, wKey).Err(); err != nil {
			return "", errors.New("从白名单中删除令牌失败")
		}*/
		payload, ok := claims["payload"].(string)
		if !ok {
			return "", errors.New("invalid payload format")
		}

		// 生成新令牌
		return jwt2.GetJwtToken(payload)
		//return claims["payload"].(string), err
	} else {
		fmt.Println(err)
	}
	return "", errors.New("invalid token structure")
}

/*func ReFreShJwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			parseToken, err := ParseToken(token)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 400,
					"msg":  "解析token失败",
				})
				c.Abort()
				exp := parseToken["exp"]
				uid := parseToken["userId"]
				if exp != nil {
					exp = parseToken["exp"].(float64)
				}
				if uid != nil {
					uid = parseToken["userId"].(string)
				}
				newToken := ""
				if int64(exp.(float64)) < time.Now().Add(time.Hour*time.Duration(1)).Unix() {
					newToken, _ = TokenHandler(uid.(string), int64(int64(exp.(float64))))
					c.Header("token", newToken)
					fmt.Println(newToken)
				}
			}
		}
	}
}*/
