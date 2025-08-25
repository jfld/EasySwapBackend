// Package config 定义了EasySwap NFT交易所后端服务的所有配置结构
// 该包负责从配置文件中读取并解析各种配置项，包括数据库、API、缓存、区块链等配置
package config

import (
	"strings"

	"github.com/joinmouse/EasySwapBase/evm/erc"        // ERC标准实现，用于处理NFT相关操作
	logging "github.com/joinmouse/EasySwapBase/logger" // 日志配置结构
	"github.com/joinmouse/EasySwapBase/stores/gdb"     // 数据库配置结构
	"github.com/spf13/viper"                          // 配置文件解析库
)

// Config 是应用程序的主配置结构体
// 包含了运行 EasySwap NFT 交易所后端服务所需的所有配置信息
type Config struct {
	Api            `toml:"api" json:"api"`                                                               // API 服务器配置，包括端口和请求限制
	ProjectCfg     *ProjectCfg     `toml:"project_cfg" mapstructure:"project_cfg" json:"project_cfg"`         // 项目基本信息配置
	Log            logging.LogConf `toml:"log" json:"log"`                                                   // 日志系统配置
	DB             gdb.Config      `toml:"db" json:"db"`                                                     // 数据库连接配置
	Kv             *KvConf         `toml:"kv" json:"kv"`                                                     // 键值存储（Redis）配置
	Evm            *erc.NftErc     `toml:"evm" json:"evm"`                                                   // EVM 区块链相关配置
	MetadataParse  *MetadataParse  `toml:"metadata_parse" mapstructure:"metadata_parse" json:"metadata_parse"` // NFT 元数据解析配置
	ChainSupported []*ChainSupported `toml:"chain_supported" mapstructure:"chain_supported" json:"chain_supported"` // 支持的区块链列表配置
}

// ProjectCfg 定义了项目的基本信息配置
type ProjectCfg struct {
	Name string `toml:"name" mapstructure:"name" json:"name"` // 项目名称，用于标识应用程序
}

// Api 定义了 HTTP API 服务器的配置参数
type Api struct {
	Port   string `toml:"port" json:"port"`     // HTTP 服务器监听端口，格式为 ":8080"
	MaxNum int64  `toml:"max_num" json:"max_num"` // 最大并发请求数量限制
}

// KvConf 定义了键值存储（主要是 Redis）的配置
type KvConf struct {
	Redis []*Redis `toml:"redis" mapstructure:"redis" json:"redis"` // Redis 服务器配置列表，支持多实例配置
}

// Redis 定义了单个 Redis 实例的连接配置
type Redis struct {
	MasterName string `toml:"master_name" mapstructure:"master_name" json:"master_name"` // Redis 主节点名称（用于 Sentinel 模式）
	Host       string `toml:"host" json:"host"`                                         // Redis 服务器地址和端口，格式为 "host:port"
	Type       string `toml:"type" json:"type"`                                         // Redis 连接类型（如 "node", "cluster", "sentinel"）
	Pass       string `toml:"pass" json:"pass"`                                         // Redis 连接密码
}

// MetadataParse 定义了 NFT 元数据解析的配置参数
// 用于从不同来源的 NFT 元数据中提取标准化信息
type MetadataParse struct {
	NameTags       []string `toml:"name_tags" mapstructure:"name_tags" json:"name_tags"`             // NFT 名称字段的可能标签名列表
	ImageTags      []string `toml:"image_tags" mapstructure:"image_tags" json:"image_tags"`          // NFT 图片 URL 字段的可能标签名列表
	AttributesTags []string `toml:"attributes_tags" mapstructure:"attributes_tags" json:"attributes_tags"` // NFT 属性字段的可能标签名列表
	TraitNameTags  []string `toml:"trait_name_tags" mapstructure:"trait_name_tags" json:"trait_name_tags"`   // NFT 特征名称字段的可能标签名列表
	TraitValueTags []string `toml:"trait_value_tags" mapstructure:"trait_value_tags" json:"trait_value_tags"` // NFT 特征值字段的可能标签名列表
}

// ChainSupported 定义了系统支持的区块链网络配置
// EasySwap 支持多链架构，可以同时处理多个区块链上的 NFT 交易
type ChainSupported struct {
	Name     string `toml:"name" mapstructure:"name" json:"name"`         // 区块链名称（如 "Ethereum", "Polygon", "BSC"）
	ChainID  int    `toml:"chain_id" mapstructure:"chain_id" json:"chain_id"` // 区块链 ID（如 Ethereum 主网是 1）
	Endpoint string `toml:"endpoint" mapstructure:"endpoint" json:"endpoint"` // 区块链 RPC 连接端点 URL
}

// UnmarshalConfig 从指定的配置文件中解析配置信息
// 该函数使用 Viper 库来读取 TOML 格式的配置文件，并支持环境变量覆盖
//
// 参数:
//   - configFilePath: 配置文件的完整路径
//
// 返回值:
//   - *Config: 解析完成的配置对象
//   - error: 解析过程中的错误，如文件不存在、格式错误等
func UnmarshalConfig(configFilePath string) (*Config, error) {
	// 设置配置文件路径
	viper.SetConfigFile(configFilePath)
	// 设置配置文件类型为 TOML
	viper.SetConfigType("toml")
	// 启用自动环境变量读取功能
	viper.AutomaticEnv()
	// 设置环境变量前缀为 "CNFT"
	viper.SetEnvPrefix("CNFT")
	// 创建环境变量名称的替换器，将 "." 替换为 "_"
	// 例如 "db.host" 对应的环境变量是 "CNFT_DB_HOST"
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	
	// 创建默认配置对象
	config, err := DefaultConfig()
	if err != nil {
		return nil, err
	}

	// 将读取的配置数据解析到配置对象中
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}
	
	return config, nil
}

// DefaultConfig 创建一个默认的配置对象
// 返回一个空的 Config 结构体，所有字段都使用默认值
//
// 返回值:
//   - *Config: 默认配置对象
//   - error: 创建过程中的错误（当前始终返回 nil）
func DefaultConfig() (*Config, error) {
	return &Config{}, nil
}
