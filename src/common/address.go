// Package common 定义了 EasySwap NFT 交易所的通用工具函数
// 该包包含了地址处理、数据转换等通用功能
package common

import (
	"github.com/ethereum/go-ethereum/common"              // 以太坊通用工具库
	"github.com/joinmouse/EasySwapBase/evm/eip"            // EIP 标准实现，包含 EIP-55 校验和地址
	"github.com/pkg/errors"                              // 错误处理库

	"github.com/joinmouse/EasySwapBackend/src/common/utils" // 内部工具函数
)

// UnifyAddress 统一化区块链地址格式
// 该函数将输入的地址转换为标准的 EIP-55 校验和地址格式
// 确保所有地址在系统中都使用统一的格式，避免因大小写不同导致的问题
//
// 参数:
//   - address: 原始地址字符串，可能包含大小写不一致的问题
//
// 返回值:
//   - string: 标准化后的 EIP-55 校验和地址
//   - error: 地址验证或转换过程中的错误
func UnifyAddress(address string) (string, error) {
	// 验证地址的基本格式
	// 地址必须大于 2 个字符（包含 0x 前缀）且符合十六进制地址格式
	if len(address) <= 2 || !common.IsHexAddress(address) {
		return "", errors.New("用户地址格式不合法")
	}

	// 使用 EIP-55 标准转换为校验和地址
	// EIP-55 通过大小写混合的方式提供地址校验功能
	addr, err := eip.ToCheckSumAddress(address)
	if err != nil {
		return "", errors.Wrap(err, "无效的地址格式")
	}

	// 再次验证转换后的地址是否有效
	// 这是一个额外的安全检查，确保地址的一致性
	if addr != utils.ToValidateAddress(addr) {
		return "", errors.Wrap(err, "地址统一化失败")
	}

	return addr, nil
}
