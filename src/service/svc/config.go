package svc

import (
	"github.com/joinmouse/EasySwapBase/evm/erc"     // ERC 标准实现
	"github.com/joinmouse/EasySwapBase/stores/xkv"  // 键值存储操作封装
	"gorm.io/gorm"                                 // GORM ORM 框架

	"github.com/joinmouse/EasySwapBackend/src/dao"  // 数据访问层
)

// CtxConfig 定义了服务上下文的配置参数
// 该结构体用于在创建 ServerCtx 时传递各种依赖组件
type CtxConfig struct {
	db      *gorm.DB       // 数据库连接实例
	dao     *dao.Dao       // 数据访问对象
	KvStore *xkv.Store     // 键值存储实例
	Evm     erc.Erc        // EVM 区块链操作接口
}

// CtxOption 定义了用于配置 ServerCtx 的选项函数类型
// 使用选项模式让 ServerCtx 的创建更加灵活和可扩展
type CtxOption func(conf *CtxConfig)

// NewServerCtx 使用选项模式创建一个新的 ServerCtx 实例
// 该函数接受一系列选项函数，用于配置不同的依赖组件
//
// 参数:
//   - options: 可变参数，一系列配置选项函数
//
// 返回值:
//   - *ServerCtx: 初始化完成的服务上下文实例
func NewServerCtx(options ...CtxOption) *ServerCtx {
	// 创建配置对象
	c := &CtxConfig{}
	
	// 应用所有配置选项
	for _, opt := range options {
		opt(c)
	}
	
	// 根据配置创建 ServerCtx
	return &ServerCtx{
		DB:      c.db,      // 设置数据库连接
		KvStore: c.KvStore, // 设置键值存储
		Dao:     c.dao,     // 设置数据访问层
	}
}

// WithKv 返回一个用于设置键值存储的选项函数
// 该函数用于在创建 ServerCtx 时注入 Redis 等键值存储服务
//
// 参数:
//   - kv: 键值存储实例
//
// 返回值:
//   - CtxOption: 配置选项函数
func WithKv(kv *xkv.Store) CtxOption {
	return func(conf *CtxConfig) {
		conf.KvStore = kv
	}
}

// WithDB 返回一个用于设置数据库连接的选项函数
// 该函数用于在创建 ServerCtx 时注入数据库连接
//
// 参数:
//   - db: GORM 数据库连接实例
//
// 返回值:
//   - CtxOption: 配置选项函数
func WithDB(db *gorm.DB) CtxOption {
	return func(conf *CtxConfig) {
		conf.db = db
	}
}

// WithDao 返回一个用于设置数据访问层的选项函数
// 该函数用于在创建 ServerCtx 时注入数据访问对象
//
// 参数:
//   - dao: 数据访问对象实例
//
// 返回值:
//   - CtxOption: 配置选项函数
func WithDao(dao *dao.Dao) CtxOption {
	return func(conf *CtxConfig) {
		conf.dao = dao
	}
}
