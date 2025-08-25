package types

import "github.com/shopspring/decimal" // 精度十进制运算库，用于处理价格等金融数据

// ItemInfo 定义了 NFT 物品的基本标识信息
// 用于唯一标识一个 NFT 物品
type ItemInfo struct {
	CollectionAddress string `json:"collection_address"` // NFT 合约地址，标识 NFT 所属的集合
	TokenID           string `json:"token_id"`           // NFT 在合约中的唯一 ID
}

// ItemPriceInfo 定义了 NFT 物品的价格相关信息
// 包含了当前挂单的价格和状态信息
type ItemPriceInfo struct {
	CollectionAddress string          `json:"collection_address"` // NFT 合约地址
	TokenID           string          `json:"token_id"`           // NFT Token ID
	Maker             string          `json:"maker"`              // 挂单制作者的地址
	Price             decimal.Decimal `json:"price"`              // 挂单价格（使用高精度十进制）
	OrderStatus       int             `json:"order_status"`       // 订单状态（0=有效, 1=已取消, 2=已成交）
}

// ItemOwner 定义了 NFT 物品的所有权信息
// 用于记录 NFT 的当前持有者
type ItemOwner struct {
	CollectionAddress string `json:"collection_address"` // NFT 合约地址
	TokenID           string `json:"token_id"`           // NFT Token ID
	Owner             string `json:"owner"`              // NFT 当前持有者的地址
}

// ItemImage 定义了 NFT 物品的图片信息
// 用于存储和返回 NFT 的显示图片
type ItemImage struct {
	CollectionAddress string `json:"collection_address"` // NFT 合约地址
	TokenID           string `json:"token_id"`           // NFT Token ID
	ImageUri          string `json:"image_uri"`          // NFT 图片的 URI 地址（IPFS 或 HTTP）
}

// ItemDetailInfo 定义了 NFT 物品的详细信息
// 包含了 NFT 的完整元数据、价格信息、挂单信息和出价信息
type ItemDetailInfo struct {
	// 基本信息
	ChainID            int    `json:"chain_id"`            // 区块链 ID
	Name               string `json:"name"`                // NFT 名称
	CollectionAddress  string `json:"collection_address"`  // NFT 合约地址
	CollectionName     string `json:"collection_name"`     // NFT 所属集合名称
	CollectionImageURI string `json:"collection_image_uri"` // 集合头像 URI
	TokenID            string `json:"token_id"`            // NFT Token ID
	
	// 媒体信息
	ImageURI  string `json:"image_uri"`  // NFT 图片 URI
	VideoType string `json:"video_type"` // 视频类型（如果有）
	VideoURI  string `json:"video_uri"`  // 视频 URI（如果有）
	
	// 价格信息
	LastSellPrice decimal.Decimal `json:"last_sell_price"` // 最近一次成交价格
	FloorPrice    decimal.Decimal `json:"floor_price"`    // 集合地板价
	
	// 所有权和市场信息
	OwnerAddress  string `json:"owner_address"`  // 当前持有者地址
	MarketplaceID int    `json:"marketplace_id"` // 交易市场 ID

	// 挂单信息（卖单）
	ListOrderID    string          `json:"list_order_id"`    // 挂单订单 ID
	ListTime       int64           `json:"list_time"`        // 挂单时间戳
	ListPrice      decimal.Decimal `json:"list_price"`       // 挂单价格
	ListExpireTime int64           `json:"list_expire_time"` // 挂单过期时间
	ListSalt       int64           `json:"list_salt"`        // 挂单的随机盐值（防重放）
	ListMaker      string          `json:"list_maker"`       // 挂单制作者地址

	// 出价信息（买单）
	BidOrderID    string          `json:"bid_order_id"`    // 出价订单 ID
	BidTime       int64           `json:"bid_time"`        // 出价时间戳
	BidExpireTime int64           `json:"bid_expire_time"` // 出价过期时间
	BidPrice      decimal.Decimal `json:"bid_price"`       // 出价价格
	BidSalt       int64           `json:"bid_salt"`        // 出价的随机盐值
	BidMaker      string          `json:"bid_maker"`       // 出价者地址
	BidType       int64           `json:"bid_type"`        // 出价类型（0=单个 NFT, 1=集合出价）
	BidSize       int64           `json:"bid_size"`        // 出价数量
	BidUnfilled   int64           `json:"bid_unfilled"`    // 未填充的出价数量
}

// ItemDetailInfoResp 定义了 NFT 物品详细信息的 API 响应结构
type ItemDetailInfoResp struct {
	Result interface{} `json:"result"` // 返回结果，通常是 ItemDetailInfo 或错误信息
}

// ListingInfo 定义了 NFT 的挂单信息
// 用于表示在特定市场上的挂单价格
type ListingInfo struct {
	MarketplaceId int32           `json:"marketplace_id"` // 交易市场 ID
	Price         decimal.Decimal `json:"price"`          // 挂单价格
}

// TraitPrice 定义了 NFT 特征的价格信息
// 用于记录和分析具有特定特征的 NFT 的价格趋势
type TraitPrice struct {
	CollectionAddress string          `json:"collection_address"` // NFT 合约地址
	TokenID           string          `json:"token_id"`           // NFT Token ID
	Trait             string          `json:"trait"`              // 特征名称（如 "Background", "Eyes"）
	TraitValue        string          `json:"trait_value"`        // 特征值（如 "Blue", "Rare"）
	Price             decimal.Decimal `json:"price"`              // 具有该特征的 NFT 的价格
}

// ItemTopTraitResp 定义了 NFT 顶级特征信息的 API 响应结构
// 用于返回最有价值或最稀有的 NFT 特征信息
type ItemTopTraitResp struct {
	Result interface{} `json:"result"` // 返回结果，通常是 TraitPrice 数组或错误信息
}
