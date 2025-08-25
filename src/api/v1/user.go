// Package v1 定义了 EasySwap NFT 交易所 API v1 版本的处理器函数
// 该包包含了所有 HTTP API 端点的业务逻辑处理
package v1

import (
	"github.com/gin-gonic/gin"                              // Gin Web框架
	"github.com/joinmouse/EasySwapBase/errcode"              // 错误码定义
	"github.com/joinmouse/EasySwapBase/kit/validator"        // 数据验证工具
	"github.com/joinmouse/EasySwapBase/xhttp"                // HTTP 响应封装工具

	"github.com/joinmouse/EasySwapBackend/src/service/svc"   // 服务上下文
	service "github.com/joinmouse/EasySwapBackend/src/service/v1" // 业务逻辑服务层
	"github.com/joinmouse/EasySwapBackend/src/types/v1"      // 数据结构定义
)

// UserLoginHandler 处理用户登录请求的 HTTP 处理器
// 该处理器实现基于区块链签名的身份验证机制，无需传统的用户名密码
// 流程:
// 1. 解析请求体中的登录数据
// 2. 验证请求参数的合法性
// 3. 验证用户的数字签名
// 4. 生成访问令牌并返回给客户端
//
// 参数:
//   - svcCtx: 服务上下文，包含数据库和缓存服务
//
// 返回值:
//   - gin.HandlerFunc: Gin 框架的处理函数
func UserLoginHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析请求体中的 JSON 数据为 LoginReq 结构体
		req := types.LoginReq{}
		if err := c.BindJSON(&req); err != nil {
			// 请求数据解析失败，返回错误响应
			xhttp.Error(c, err)
			return
		}

		// 验证请求参数的完整性和合法性
		// 检查必填字段和数据格式
		if err := validator.Verify(&req); err != nil {
			xhttp.Error(c, errcode.NewCustomErr(err.Error()))
			return
		}

		// 调用业务逻辑层处理登录逻辑
		// 包括签名验证、用户信息查询、令牌生成等
		res, err := service.UserLogin(c.Request.Context(), svcCtx, req)
		if err != nil {
			// 登录失败，返回错误信息
			xhttp.Error(c, errcode.NewCustomErr(err.Error()))
			return
		}

		// 登录成功，返回用户信息和访问令牌
		xhttp.OkJson(c, types.UserLoginResp{
			Result: res,
		})
	}
}

// GetLoginMessageHandler 处理获取登录消息请求的 HTTP 处理器
// 该处理器为指定的用户地址生成一个唯一的消息，用于后续的数字签名验证
// 消息通常包含随机数、时间戳等信息，防止重放攻击
//
// 参数:
//   - svcCtx: 服务上下文
//
// 路由参数:
//   - address: 用户的区块链地址
//
// 返回值:
//   - gin.HandlerFunc: Gin 框架的处理函数
func GetLoginMessageHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 URL 路径参数中获取用户地址
		address := c.Params.ByName("address")
		if address == "" {
			// 地址参数为空，返回错误响应
			xhttp.Error(c, errcode.NewCustomErr("用户地址不能为空"))
			return
		}

		// 调用业务逻辑层生成登录消息
		// 服务层会验证地址格式并生成安全的消息
		res, err := service.GetUserLoginMsg(c.Request.Context(), svcCtx, address)
		if err != nil {
			// 消息生成失败，返回错误信息
			xhttp.Error(c, errcode.NewCustomErr(err.Error()))
			return
		}

		// 成功返回登录消息
		xhttp.OkJson(c, res)
	}
}

// GetSigStatusHandler 处理获取用户签名状态请求的 HTTP 处理器
// 该处理器查询指定用户是否已经完成了数字签名认证
// 可用于客户端轮询用户的认证状态
//
// 参数:
//   - svcCtx: 服务上下文
//
// 路由参数:
//   - address: 用户的区块链地址
//
// 返回值:
//   - gin.HandlerFunc: Gin 框架的处理函数
func GetSigStatusHandler(svcCtx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 URL 路径参数中获取用户地址
		userAddr := c.Params.ByName("address")
		if userAddr == "" {
			// 地址参数为空，返回错误响应
			xhttp.Error(c, errcode.NewCustomErr("用户地址不能为空"))
			return
		}

		// 调用业务逻辑层查询签名状态
		// 服务层会查询数据库或缓存中的用户认证信息
		res, err := service.GetSigStatusMsg(c.Request.Context(), svcCtx, userAddr)
		if err != nil {
			// 查询失败，返回错误信息
			xhttp.Error(c, errcode.NewCustomErr(err.Error()))
			return
		}

		// 成功返回签名状态
		xhttp.OkJson(c, res)
	}
}
