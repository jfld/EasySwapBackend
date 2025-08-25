// Package dao 定义了EasySwap NFT交易所的数据访问层
// 该包负责封装所有的数据库操作，包括 NFT、用户、订单、交易等数据的 CRUD 操作
// 同时管理缓存操作，提高数据访问性能
package dao

import (
	"context"

	"github.com/joinmouse/EasySwapBase/stores/xkv"  // 键值存储操作封装
	"gorm.io/gorm"                                 // GORM ORM 框架
)

// Dao 表示数据访问对象，封装了数据库和缓存操作
// 它是 EasySwap NFT 交易所数据持久化层的核心组件
// 提供统一的数据访问接口，支持事务处理和缓存管理
type Dao struct {
	ctx     context.Context  // 上下文对象，用于传递请求范围内的信息
	DB      *gorm.DB         // GORM 数据库连接，用于执行 SQL 操作
	KvStore *xkv.Store       // 键值存储实例（Redis），用于缓存和会话管理
}

// New 创建一个新的数据访问对象实例
// 该函数初始化 Dao 结构体，将数据库连接和缓存实例传入
//
// 参数:
//   - ctx: 上下文对象，用于传递请求相关信息
//   - db: GORM 数据库连接实例
//   - kvStore: 键值存储实例，用于缓存操作
//
// 返回值:
//   - *Dao: 初始化完成的数据访问对象
func New(ctx context.Context, db *gorm.DB, kvStore *xkv.Store) *Dao {
	return &Dao{
		ctx:     ctx,     // 保存上下文
		DB:      db,      // 保存数据库连接
		KvStore: kvStore, // 保存缓存实例
	}
}
