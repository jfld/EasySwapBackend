// Package svc 定义了EasySwap NFT交易所后端服务的服务上下文
// 该包负责管理应用程序的所有依赖服务，包括数据库、缓存、区块链连接等
package svc

import (
	"context"

	"github.com/joinmouse/EasySwapBase/chain/nftchainservice" // NFT 区块链服务，用于与区块链交互
	"github.com/joinmouse/EasySwapBase/logger/xzap"         // 结构化日志库
	"github.com/joinmouse/EasySwapBase/stores/gdb"          // 数据库操作封装
	"github.com/joinmouse/EasySwapBase/stores/xkv"          // 键值存储操作封装
	"github.com/pkg/errors"                                // 错误处理库
	"github.com/zeromicro/go-zero/core/stores/cache"        // go-zero 缓存组件
	"github.com/zeromicro/go-zero/core/stores/kv"           // go-zero 键值存储组件
	"github.com/zeromicro/go-zero/core/stores/redis"        // go-zero Redis 组件
	"gorm.io/gorm"                                         // GORM ORM 框架

	"github.com/joinmouse/EasySwapBackend/src/config"       // 配置管理模块
	"github.com/joinmouse/EasySwapBackend/src/dao"          // 数据访问层
)

// ServerCtx 表示服务器的上下文信息
// 它包含了运行 EasySwap NFT 交易所后端服务所需的所有依赖组件
// 该结构体通过依赖注入的方式统一管理各种服务
type ServerCtx struct {
	C        *config.Config                        // 应用程序配置
	DB       *gorm.DB                              // 数据库连接实例，用于数据持久化
	Dao      *dao.Dao                              // 数据访问对象，封装了所有数据库操作
	KvStore  *xkv.Store                            // 键值存储实例，主要用于缓存和会话管理
	RankKey  string                                // 排行榜缓存的键名前缀
	NodeSrvs map[int64]*nftchainservice.Service    // 区块链服务实例映射，键为链ID，值为对应的区块链服务
}

// NewServiceContext 创建一个新的服务上下文实例
// 该函数根据提供的配置初始化所有必要的服务组件，包括:
// - 日志系统
// - Redis 缓存服务
// - 数据库连接
// - 区块链服务（支持多链）
// - 数据访问层
//
// 参数:
//   - c: 应用程序配置对象
//
// 返回值:
//   - *ServerCtx: 初始化完成的服务上下文
//   - error: 初始化过程中的错误
func NewServiceContext(c *config.Config) (*ServerCtx, error) {
	var err error

	// 初始化日志系统
	// 根据配置设置日志级别、输出格式等
	_, err = xzap.SetUp(c.Log)
	if err != nil {
		return nil, err
	}

	// 构建 Redis 配置
	// 将配置文件中的 Redis 配置转换为 go-zero 所需的格式
	var kvConf kv.KvConf
	for _, con := range c.Kv.Redis {
		kvConf = append(kvConf, cache.NodeConf{
			RedisConf: redis.RedisConf{
				Host: con.Host,  // Redis 服务器地址
				Type: con.Type,  // Redis 连接类型
				Pass: con.Pass,  // Redis 连接密码
			},
			Weight: 1,           // 节点权重，用于负载均衡
		})
	}

	// 初始化 Redis 存储
	store := xkv.NewStore(kvConf)
	
	// 初始化数据库连接
	db, err := gdb.NewDB(&c.DB)
	if err != nil {
		return nil, err
	}

	// 初始化区块链服务
	// 为每个支持的区块链创建对应的服务实例
	nodeSrvs := make(map[int64]*nftchainservice.Service)
	for _, supported := range c.ChainSupported {
		// 为每个区块链创建 NFT 链上服务
		nodeSrvs[int64(supported.ChainID)], err = nftchainservice.New(
			context.Background(),
			supported.Endpoint,           // 区块链 RPC 端点
			supported.Name,               // 区块链名称
			supported.ChainID,            // 区块链 ID
			c.MetadataParse.NameTags,     // NFT 名称字段标签
			c.MetadataParse.ImageTags,    // NFT 图片字段标签
			c.MetadataParse.AttributesTags, // NFT 属性字段标签
			c.MetadataParse.TraitNameTags,  // NFT 特征名称字段标签
			c.MetadataParse.TraitValueTags, // NFT 特征值字段标签
		)

		if err != nil {
			return nil, errors.Wrap(err, "初始化区块链同步服务失败")
		}
	}

	// 初始化数据访问层
	dao := dao.New(context.Background(), db, store)
	
	// 使用选项模式创建服务上下文
	serverCtx := NewServerCtx(
		WithDB(db),     // 注入数据库连接
		WithKv(store),  // 注入键值存储
		WithDao(dao),   // 注入数据访问层
	)
	
	// 设置其他属性
	serverCtx.C = c               // 保存配置引用
	serverCtx.NodeSrvs = nodeSrvs // 保存区块链服务映射

	return serverCtx, nil
}
