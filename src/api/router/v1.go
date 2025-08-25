package router

import (
	"github.com/gin-gonic/gin"                                 // Gin Web框架

	"github.com/joinmouse/EasySwapBackend/src/api/middleware"   // 中间件包
	v1 "github.com/joinmouse/EasySwapBackend/src/api/v1"        // API v1 版本处理器
	"github.com/joinmouse/EasySwapBackend/src/service/svc"      // 服务上下文
)

// loadV1 加载 API v1 版本的所有路由配置
// 该函数定义了 EasySwap NFT 交易所的所有 API 端点，包括:
// - 用户认证相关 API
// - NFT 集合和物品管理 API
// - 交易活动查询 API
// - 用户投资组合管理 API
// - 订单管理 API
//
// 参数:
//   - r: Gin 路由器实例
//   - svcCtx: 服务上下文，包含数据库、缓存等服务
func loadV1(r *gin.Engine, svcCtx *svc.ServerCtx) {
	// 创建 API v1 版本的路由组
	apiV1 := r.Group("/api/v1")

	// 用户认证相关路由组
	// 处理用户登录、签名验证等功能
	user := apiV1.Group("/user")
	{
		user.GET("/:address/login-message", v1.GetLoginMessageHandler(svcCtx)) // 获取登录签名消息，用于用户签名认证
		user.POST("/login", v1.UserLoginHandler(svcCtx))                       // 用户登录接口，验证签名并返回令牌
		user.GET("/:address/sig-status", v1.GetSigStatusHandler(svcCtx))       // 获取用户签名状态
	}

	// NFT 集合和物品相关路由组
	// 处理 NFT 集合信息、物品详情、交易信息等
	collections := apiV1.Group("/collections")
	{
		// NFT 集合管理 API
		collections.GET("/:address", v1.CollectionDetailHandler(svcCtx))                  // 获取指定 NFT 集合的详细信息
		collections.GET("/:address/bids", v1.CollectionBidsHandler(svcCtx))               // 获取指定集合的所有出价信息
		collections.GET("/:address/:token_id/bids", v1.CollectionItemBidsHandler(svcCtx)) // 获取指定 NFT 物品的出价信息
		collections.GET("/:address/items", v1.CollectionItemsHandler(svcCtx))             // 获取指定集合下的所有 NFT 物品

		// NFT 物品详情 API
		collections.GET("/:address/:token_id", v1.ItemDetailHandler(svcCtx))     // 获取 NFT 物品的详细信息（包括价格、所有者等）
		collections.GET("/:address/:token_id/traits", v1.ItemTraitsHandler(svcCtx)) // 获取 NFT 物品的属性特征信息
		collections.GET("/:address/top-trait", v1.ItemTopTraitPriceHandler(svcCtx)) // 获取集合中最高价的特征信息
		
		// NFT 媒体和元数据 API
		collections.GET("/:address/:token_id/image", 
			middleware.CacheApi(svcCtx.KvStore, 60), // 缓存 60 秒
			v1.GetItemImageHandler(svcCtx))          // 获取 NFT 物品的图片信息
		collections.POST("/:address/:token_id/metadata", v1.ItemMetadataRefreshHandler(svcCtx)) // 刷新 NFT 物品的元数据
		
		// NFT 交易历史和所有权 API
		collections.GET("/:address/history-sales", v1.HistorySalesHandler(svcCtx))       // 获取 NFT 集合的销售历史信息
		collections.GET("/:address/:token_id/owner", v1.ItemOwnerHandler(svcCtx))       // 获取 NFT 物品的当前持有者信息

		// NFT 排行榜 API
		collections.GET("/ranking", 
			middleware.CacheApi(svcCtx.KvStore, 60), // 缓存 60 秒
			v1.TopRankingHandler(svcCtx))            // 获取 NFT 集合排行榜信息
	}

	// 交易活动相关路由组
	// 处理交易历史、交易事件等信息
	activities := apiV1.Group("/activities")
	{
		activities.GET("", v1.ActivityMultiChainHandler(svcCtx)) // 获取多链交易活动信息（买卖、转让等）
	}

	// 用户投资组合相关路由组
	// 处理用户持有的 NFT、挂单、出价等信息
	portfolio := apiV1.Group("/portfolio")
	{
		portfolio.GET("/collections", v1.UserMultiChainCollectionsHandler(svcCtx)) // 获取用户在多链上持有的 NFT 集合信息
		portfolio.GET("/items", v1.UserMultiChainItemsHandler(svcCtx))             // 获取用户在多链上持有的 NFT 物品信息
		portfolio.GET("/listings", v1.UserMultiChainListingsHandler(svcCtx))       // 获取用户在多链上的挂单信息
		portfolio.GET("/bids", v1.UserMultiChainBidsHandler(svcCtx))               // 获取用户在多链上的出价信息
	}

	// 订单管理相关路由组
	// 处理交易订单查询和管理
	orders := apiV1.Group("/bid-orders")
	{
		orders.GET("", v1.OrderInfosHandler(svcCtx)) // 批量查询出价订单信息
	}
}
