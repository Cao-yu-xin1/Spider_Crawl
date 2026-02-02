package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//const secretKey = "1234567890"

//RefreshToken 刷新JWT令牌

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token != "" {
			parseToken, err := ParseToken(token)
			// 解析token有错误（如过期、签名错误）才走刷新逻辑
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 400,
					"msg":  "解析token失败",
				})
				c.Abort()
				return // 解析失败直接终止，避免后续无效操作
			}

			// 解析成功，才检查是否需要刷新
			expVal, expOk := parseToken["exp"].(float64)
			uidVal, uidOk := parseToken["userId"].(string)
			// 确保exp和userId存在且类型正确
			if expOk && uidOk {
				// 检查token是否即将过期（1小时内过期则刷新）
				if int64(expVal) < time.Now().Add(1*time.Hour).Unix() {
					newToken, err := GetJwtToken(int(expVal), uidVal)
					if err == nil && newToken != "" {
						// 将新token放入响应头
						c.Header("token", newToken)
						fmt.Println("刷新后的token：", newToken)
					}
				}
			}
		}
		// 继续执行后续中间件/接口逻辑
		c.Next()
	}
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
