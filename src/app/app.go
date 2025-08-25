// Package app 定义了EasySwap NFT交易所的应用程序核心结构
// 该包负责管理应用程序的生命周期，包括启动、配置管理和服务器运行
package app

import (
	"context"

	"github.com/gin-gonic/gin"                              // Gin Web框架，用于构建REST API
	"github.com/joinmouse/EasySwapBase/logger/xzap"         // 日志库，基于zap的结构化日志
	"go.uber.org/zap"                                       // Uber的高性能日志库

	"github.com/joinmouse/EasySwapBackend/src/config"       // 配置管理模块
	"github.com/joinmouse/EasySwapBackend/src/service/svc"  // 服务上下文模块
)

// Platform 表示EasySwap NFT交易所的主应用程序平台
// 它封装了应用程序运行所需的所有组件，包括配置、HTTP路由器和服务上下文
type Platform struct {
	config    *config.Config    // 应用程序配置，包含数据库、API、区块链等配置信息
	router    *gin.Engine       // Gin HTTP路由器，处理所有的API请求
	serverCtx *svc.ServerCtx    // 服务上下文，包含数据库连接、缓存、区块链服务等
}

// NewPlatform 创建一个新的应用程序平台实例
// 参数:
//   - config: 应用程序配置，包含所有必要的配置信息
//   - router: 已初始化的Gin路由器，包含所有API端点
//   - serverCtx: 服务上下文，包含数据库、缓存等服务
//
// 返回值:
//   - *Platform: 初始化完成的平台实例
//   - error: 初始化过程中的错误（当前始终返回 nil）
func NewPlatform(config *config.Config, router *gin.Engine, serverCtx *svc.ServerCtx) (*Platform, error) {
	return &Platform{
		config:    config,     // 保存应用程序配置
		router:    router,     // 保存HTTP路由器
		serverCtx: serverCtx,  // 保存服务上下文
	}, nil
}

// Start 启动应用程序平台
// 该方法会记录启动信息并开始HTTP服务器的监听
// 服务器将在配置指定的端口上接收和处理HTTP请求
// 此方法会阻塞运行，直到服务器关闭或发生错误
func (p *Platform) Start() {
	// 记录服务器启动日志，包含监听端口信息
	xzap.WithContext(context.Background()).Info(
		"EasySwap NFT交易所后端服务器已启动", 
		zap.String("port", p.config.Api.Port),  // 记录监听端口
	)
	
	// 启动Gin HTTP服务器
	// 在指定端口上开始监听并处理HTTP请求
	if err := p.router.Run(p.config.Api.Port); err != nil {
		// 如果服务器启动失败，直接崩溃程序
		panic(err)
	}
}
