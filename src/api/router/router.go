// Package router 定义了EasySwap NFT交易所后端服务的HTTP路由配置
// 该包负责初始化Gin Web框架、配置中间件、CORS策略和所有API路由
package router

import (
	"time"

	"github.com/gin-contrib/cors"                             // Gin CORS 中间件
	"github.com/gin-gonic/gin"                                // Gin Web 框架

	"github.com/joinmouse/EasySwapBackend/src/api/middleware" // 自定义中间件
	"github.com/joinmouse/EasySwapBackend/src/service/svc"    // 服务上下文
)

// NewRouter 创建并配置一个新的 Gin HTTP 路由器
// 该函数负责:
// 1. 初始化 Gin 引擎并设置运行模式
// 2. 配置全局中间件（错误恢复、日志记录、CORS）
// 3. 加载所有API版本的路由配置
//
// 参数:
//   - svcCtx: 服务上下文，包含数据库、缓存等服务实例
//
// 返回值:
//   - *gin.Engine: 配置完成的 Gin 路由器实例
func NewRouter(svcCtx *svc.ServerCtx) *gin.Engine {
	// 强制使用彩色输出，提高日志可读性
	gin.ForceConsoleColor()
	// 设置为发布模式，减少调试信息的输出
	gin.SetMode(gin.ReleaseMode)
	
	// 创建新的 Gin 引擎实例
	r := gin.New()
	
	// 注册全局中间件
	r.Use(middleware.RecoverMiddleware()) // 恢复中间件，捕获panic并返回错误响应
	r.Use(middleware.RLog())              // 日志中间件，记录请求和响应信息

	// 配置 CORS（跨域资源共享）中间件
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // 允许所有来源的跨域请求
		// 允许的 HTTP 方法
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		// 允许的请求头
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"X-CSRF-Token",
			"Authorization",
			"AccessToken",
			"Token",
		},
		// 向客户端暴露的响应头
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"X-GW-Error-Code",
			"X-GW-Error-Message",
		},
		AllowCredentials: true,          // 允许发送身份凭证（如 Cookies）
		MaxAge:           1 * time.Hour, // 预检请求的缓存时间
	}))
	
	// 加载 API v1 版本路由
	loadV1(r, svcCtx)

	return r
}
