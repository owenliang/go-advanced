package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type YuerContext struct {
	context.Context
	Gin *gin.Context
}

type ZDMHandleFunc func (c *YuerContext)

func WithYuerContext(zdmHandle ZDMHandleFunc) gin.HandlerFunc {
	return func (c *gin.Context) {
		// 可以在gin.Context中设置key-value
		c.Set("trace", "假设这是一个调用链追踪sdk")

		// 全局超时控制
		timeoutCtx, _ := context.WithTimeout(c, 5 * time.Second)
		// ZDM上下文
		yuerCtx := YuerContext{Context: timeoutCtx, Gin: c}

		// 回调接口
		zdmHandle(&yuerCtx)
	}
}

// 模拟一个MYSQL查询
func dbQuery(ctx context.Context, sql string) {
	// 模拟调用链埋点
	trace := ctx.Value("trace").(string)

	// 模拟长时间逻辑阻塞, 被context的5秒超时中断
	<- ctx.Done()

	fmt.Println(trace)
}

func main() {
	r := gin.New()

	r.GET("/test", WithYuerContext(func(c *YuerContext) {
		// 业务层处理
		dbQuery(c, "select * from xxx")
		// 调用gin应答
		c.Gin.String(200, "请求完成")
	}))

	r.Run()
}
