package encrypt

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

// Encoding 封装go底层的 encoding 包
// 目前支持：
// 		base64
// 		base32
// 		hex(base16)
type Encoding interface {
	EncodeToString(src []byte) string
	DecodeString(s string) ([]byte, error)
}

// Base64Encoding Base64 编码简化
type Base64Encoding struct {
	internalEncoding *base64.Encoding
}

// NewBase64Encoding 创建一个 Base64Encoding
// e 可以指定 Encoding
func NewBase64Encoding(e ...*base64.Encoding) *Base64Encoding {
	if len(e) == 0 {
		e = []*base64.Encoding{base64.StdEncoding}
	}

	return &Base64Encoding{
		internalEncoding: e[0],
	}
}

// EncodeToString 编码字节
func (b *Base64Encoding) EncodeToString(src []byte) string {
	return b.internalEncoding.EncodeToString(src)
}

// DecodeString 解码字节，请勿忽略 error
func (b *Base64Encoding) DecodeString(s string) ([]byte, error) {
	return b.internalEncoding.DecodeString(s)
}

// Base32Encoding Base32 编码简化
type Base32Encoding struct {
	internalEncoding *base32.Encoding
}

// NewBase32Encoding 创建一个 Base64Encoding
// e 可以指定 Encoding
func NewBase32Encoding(e ...*base32.Encoding) *Base32Encoding {
	if len(e) == 0 {
		e = []*base32.Encoding{base32.StdEncoding}
	}

	return &Base32Encoding{
		internalEncoding: e[0],
	}
}

// NewBase32Human 生成的编码中没有 O o 0 I 1 T 这种易被混淆的字符
// 对人类阅读十分友好
func NewBase32Human() *Base32Encoding {
	e := base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ23456789").WithPadding(base32.NoPadding)
	return &Base32Encoding{
		internalEncoding: e,
	}
}

// EncodeToString 编码字节
func (b *Base32Encoding) EncodeToString(src []byte) string {
	return b.internalEncoding.EncodeToString(src)
}

// DecodeString 解码，请勿忽略 error
func (b *Base32Encoding) DecodeString(s string) ([]byte, error) {
	return b.internalEncoding.DecodeString(s)
}

// HexEncoding Hex 编码简化
type HexEncoding struct {
}

// NewHexEncoding 创建一个 HexEncoding
func NewHexEncoding() *HexEncoding {
	return &HexEncoding{}
}

// EncodeToString 编码字节
func (h *HexEncoding) EncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}

// DecodeString 解码，请勿忽略 error
func (h *HexEncoding) DecodeString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
