// Package main 是EasySwap NFT交易所后端服务的主入口包
// 该服务提供了完整的NFT交易功能，包括用户管理、NFT管理、订单管理、交易管理等核心功能
// 支持多链架构，可以同时处理多个区块链网络上的NFT交易
package main

import (
	"flag"               // 用于解析命令行参数
	_ "net/http/pprof"   // 导入pprof包，用于性能分析和调试

	"github.com/joinmouse/EasySwapBackend/src/api/router"   // 导入路由模块
	"github.com/joinmouse/EasySwapBackend/src/app"          // 导入应用程序核心模块
	"github.com/joinmouse/EasySwapBackend/src/config"       // 导入配置管理模块
	"github.com/joinmouse/EasySwapBackend/src/service/svc"  // 导入服务上下文模块
)

// 常量定义
const (
	// repoRoot 仓库根目录，当前为空字符串表示使用当前目录
	repoRoot = ""
	// defaultConfigPath 默认配置文件路径，指向config目录下的config.toml文件
	defaultConfigPath = "./config/config.toml"
)

// main 是程序的主入口函数
// 负责初始化配置、验证区块链配置、创建服务上下文、初始化路由器并启动应用程序
func main() {
	// 解析命令行参数，获取配置文件路径
	// -conf 参数用于指定配置文件路径，默认使用 defaultConfigPath
	conf := flag.String("conf", defaultConfigPath, "配置文件路径")
	flag.Parse()
	
	// 从指定的配置文件中解析配置信息
	// 配置文件包含数据库连接、API端口、支持的区块链网络等信息
	c, err := config.UnmarshalConfig(*conf)
	if err != nil {
		panic(err)
	}

	// 验证支持的区块链配置
	// 确保每个支持的区块链都有有效的链ID和名称
	for _, chain := range c.ChainSupported {
		if chain.ChainID == 0 || chain.Name == "" {
			panic("无效的区块链配置：链ID不能为0，链名称不能为空")
		}
	}

	// 创建服务上下文，包含数据库连接、Redis连接、区块链服务等
	// 服务上下文是整个应用程序的依赖注入容器
	serverCtx, err := svc.NewServiceContext(c)
	if err != nil {
		panic(err)
	}
	
	// 初始化路由器，设置所有的API端点
	// 路由器配置了中间件、CORS策略和API版本路由
	r := router.NewRouter(serverCtx)
	
	// 创建应用程序平台实例
	// 平台封装了配置、路由器和服务上下文
	app, err := app.NewPlatform(c, r, serverCtx)
	if err != nil {
		panic(err)
	}
	
	// 启动应用程序服务器
	// 开始监听HTTP请求并处理NFT交易相关的API调用
	app.Start()
}
