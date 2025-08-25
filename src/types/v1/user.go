// Package types 定义了 EasySwap NFT 交易所 API v1 版本的数据结构
// 该包包含了所有 API 请求和响应的数据结构定义，用于 JSON 序列化和反序列化
package types

// LoginReq 定义了用户登录请求的数据结构
// 使用区块链签名进行身份验证，无需传统的用户名密码
type LoginReq struct {
	ChainID   int    `json:"chain_id"`  // 区块链 ID，用于标识用户所在的区块链网络
	Message   string `json:"message"`   // 待签名的消息内容，由服务器生成
	Signature string `json:"signature"` // 用户对消息的数字签名
	Address   string `json:"address"`   // 用户的区块链地址（钉包地址）
}

// UserLoginInfo 定义了用户登录成功后的信息
// 包含了访问令牌和用户权限状态
type UserLoginInfo struct {
	Token     string `json:"token"`      // JWT 令牌或会话 ID，用于后续 API 访问身份验证
	IsAllowed bool   `json:"is_allowed"` // 用户是否被允许访问系统（用于权限控制）
}

// UserLoginResp 定义了用户登录请求的响应数据结构
// 使用通用的 interface{} 类型以支持不同类型的响应数据
type UserLoginResp struct {
	Result interface{} `json:"result"` // 登录结果，可能是 UserLoginInfo 或错误信息
}

// UserLoginMsgResp 定义了获取登录消息的响应数据结构
// 用于返回用户需要签名的消息内容
type UserLoginMsgResp struct {
	Address string `json:"address"` // 用户地址，用于确认身份
	Message string `json:"message"` // 需要签名的消息内容，通常包含随机数和时间戳
}

// UserSignStatusResp 定义了用户签名状态的响应数据结构
// 用于查询用户是否已经完成签名操作
type UserSignStatusResp struct {
	IsSigned bool `json:"is_signed"` // 用户是否已经完成签名认证
}
