package middleware

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @description: 协商缓存中间件
// @param: c *gin.Context
// NegotiationCacheMiddleware 协商缓存中间件(你的缓存代码)
// 1. 计算当前内容的 ETag（根据内容版本）
// 2. 检查客户端携带的 If-None-Match（核心校验）
// 3. 缓存命中则返回 304（核心逻辑）
// 4. 设置响应头（让浏览器缓存 ETag）

func NegotiationCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentVersion := "v1.0.0"
		etag := fmt.Sprintf(`W/"%x"`, md5.Sum([]byte(contentVersion)))

		// 2. 检查客户端携带的 ETag（核心校验）
		clientETag := c.GetHeader("If-None-Match")

		// 3. 缓存命中则返回 304（核心逻辑）
		if clientETag == etag {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}

		// 4. 设置响应头（让浏览器缓存 ETag）
		c.Header("ETag", etag)
		c.Header("Cache-Control", "no-cache") // 可选保留，强制走协商缓存

		c.Next()
	}
}

/*
marshal, err := json.Marshal(goods)
	if err != nil {
		begin.Rollback()
		return nil, errors.New("商品序列化失败")
	}
	key := fmt.Sprintf("goods:%s", in.Name)
	nx := config.Rdb.SetNX(config.Ctx, key, marshal, time.Minute*30).Val()
	if !nx {
		begin.Rollback()
		return nil, errors.New("商品已存在,请勿重复上架")
	}
	begin.Commit()
------------------------------
key := fmt.Sprintf("goods:%s", in.Name)
	err := config.Rdb.Exists(config.Ctx, key).Val()
	if err == 0 {
		return nil, errors.New("商品不存在")
	} else {
		result, _ := config.Rdb.Get(config.Ctx, key).Result()
		json.Unmarshal([]byte(result), &goods)
	}
*/

/*
// NegotiationCacheMiddleware 协商缓存中间件
func NegotiationCacheMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 生成 ETag (基于内容特征或版本号)
        contentVersion := "v1.0.0" // 或基于数据哈希
        etag := fmt.Sprintf(`W/"%x"`, md5.Sum([]byte(contentVersion)))

        // 2. 设置 Last-Modified
        lastModified := time.Now().Add(-1 * time.Hour).Format(http.TimeFormat)

        // 3. 检查客户端请求头（协商过程）
        clientETag := c.GetHeader("If-None-Match")
        clientModified := c.GetHeader("If-Modified-Since")

        // 4. 验证是否匹配
        if clientETag == etag {
            // ETag 匹配，返回 304
            c.AbortWithStatus(http.StatusNotModified)
            return
        }

        if clientModified == lastModified {
            // 时间匹配，返回 304
            c.AbortWithStatus(http.StatusNotModified)
            return
        }

        // 5. 设置响应头（只有不匹配时才设置）
        c.Header("ETag", etag)
        c.Header("Last-Modified", lastModified)

        // 建议同时禁用强缓存，强制走协商缓存
        c.Header("Cache-Control", "no-cache") // 或 "max-age=0, must-revalidate"

        c.Next()
    }
}
*/
