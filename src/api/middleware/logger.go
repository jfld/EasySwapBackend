// Package middleware 定义了EasySwap NFT交易所的HTTP中间件
// 该包包含了日志记录、错误恢复、身份验证等各种中间件功能
package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"                     // Gin Web框架
	"github.com/joinmouse/EasySwapBase/logger/xzap" // 结构化日志库
	"go.uber.org/zap"                              // Uber的高性能日志库
	"go.uber.org/zap/zapcore"                      // Zap日志库核心组件
)

// BodyLogWriter 是一个自定义的响应写入器
// 它封装了 Gin 的原始 ResponseWriter，在写入响应的同时保存响应内容用于日志记录
type BodyLogWriter struct {
	gin.ResponseWriter            // 嵌入 Gin 的原始 ResponseWriter
	body              *bytes.Buffer // 用于存储响应体内容的缓冲区
}

// Write 实现 io.Writer 接口的 Write 方法
// 在写入响应数据的同时，将数据保存到内部缓冲区供日志记录使用
func (w BodyLogWriter) Write(b []byte) (int, error) {
	// 同时写入缓冲区和原始响应写入器
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
// WriteString 实现字符串写入方法
// 在写入响应字符串数据的同时，将数据保存到内部缓冲区供日志记录使用
func (w BodyLogWriter) WriteString(s string) (int, error) {
	// 同时写入缓冲区和原始响应写入器
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// RLog 是一个用于记录 HTTP 请求和响应的中间件函数
// 该中间件会记录请求和响应的详细信息，包括:
// 1. 请求的 URL 路径、查询参数和请求体
// 2. 响应的状态码和响应体
// 3. 请求处理时间
// 4. 客户端 IP、User-Agent 等元数据
// 5. 错误信息（如果有）
//
// 返回值:
//   - gin.HandlerFunc: Gin 中间件函数
func RLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取原始请求路径和查询参数（避免被其他中间件修改）
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 读取并保存请求体内容
		// 使用 TeeReader 在读取的同时保存数据副本
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		requestBody, _ := ioutil.ReadAll(tee)
		// 重新设置请求体，供后续处理器使用
		c.Request.Body = ioutil.NopCloser(&buf)
		
		// 创建自定义的响应写入器，用于捕获响应内容
		bodyLogWriter := &BodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = bodyLogWriter

		// 记录请求开始时间
		start := time.Now()

		// 调用下一个处理器函数
		c.Next()

		// 获取响应体内容
		responseBody := bodyLogWriter.body.Bytes()
		// 获取上下文相关的日志记录器
		logger := xzap.WithContext(c.Request.Context())
		
		if len(c.Errors) > 0 {
			// 如果请求处理过程中出现错误，记录所有错误信息
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			// 计算请求处理的延迟时间（毫秒）
			latency := float64(time.Now().Sub(start).Nanoseconds() / 1000000.0)
			
			// 构建日志字段，记录请求和响应的详细信息
			fields := []zapcore.Field{
				zap.Int("status", c.Writer.Status()),                         // HTTP 状态码
				zap.String("method", c.Request.Method),                       // HTTP 请求方法
				zap.String("function", c.HandlerName()),                     // 处理函数名
				zap.String("path", path),                                    // 请求路径
				zap.String("query", query),                                  // 查询参数
				zap.String("ip", c.ClientIP()),                              // 客户端 IP 地址
				zap.String("user-agent", c.Request.UserAgent()),             // 客户端 User-Agent
				zap.String("token", c.Request.Header.Get("session_id")),     // 会话 ID
				zap.String("content-type", c.Request.Header.Get("Content-Type")), // 请求内容类型
				zap.Float64("latency", latency),                             // 请求处理延迟
				zap.String("request", string(requestBody)),                  // 请求体内容
				zap.String("response", string(responseBody)),                // 响应体内容
			}
			// 记录成功的请求处理日志
			logger.Info("EasySwap API 请求处理完成", fields...)
		}
	}
}
