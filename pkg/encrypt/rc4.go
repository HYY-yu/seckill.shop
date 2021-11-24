package encrypt

import (
	"crypto/rc4"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

// rc4 加密算法
// 1. 速度快，密文长度等于正文
// 2. 对密钥长度无要求

type GoRC4 struct {
}

func NewGoRC4() *GoRC4 {
	return &GoRC4{}
}

// Encrypt encrypts the first block in src into dst.
// Dst and src may point at the same memory.
func (self *GoRC4) Encrypt(src, key []byte) ([]byte, error) {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(src))
	cipher.XORKeyStream(dst, src)
	return dst, nil
}

// Decrypt decrypts the first block in src into dst.
// Dst and src may point at the same memory.
func (self *GoRC4) Decrypt(src, key []byte) ([]byte, error) {
	cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(src))
	cipher.XORKeyStream(dst, src)
	return dst, nil
}

// EncryptHex encrypts the first block in src into dst.
func (self *GoRC4) EncryptHex(src, key string) (string, error) {
	ciphertext, err := self.Encrypt([]byte(src), []byte(key))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ciphertext), nil
}

// DecryptHex decrypts the first block in src into dst.
func (self *GoRC4) DecryptHex(src, key string) (string, error) {
	plaintext, err := hex.DecodeString(src)
	if err != nil {
		return "", err
	}
	dst, err := self.Decrypt(plaintext, []byte(key))
	return string(dst), nil
}

func (self *GoRC4) EnBase64URL(src, key string) (string, error) {
	ciphertext, err := self.Encrypt([]byte(src), []byte(key))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func (self *GoRC4) DeBase64URL(src, key string) (string, error) {
	plaintext, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	dst, err := self.Decrypt(plaintext, []byte(key))
	return string(dst), nil
}

// 生成的密文中没有 O o 0 I 1 T 这种易被混淆的字符
func (self *GoRC4) EnBase32Ma(src, key string) (string, error) {
	ciphertext, err := self.Encrypt([]byte(src), []byte(key))
	if err != nil {
		return "", err
	}
	return base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ23456789").WithPadding(base32.NoPadding).EncodeToString(ciphertext), nil
}

func (self *GoRC4) DeBase32Ma(src, key string) (string, error) {
	plaintext, err := base32.NewEncoding("ABCDEFGHJKLMNPQRSTUVWXYZ23456789").WithPadding(base32.NoPadding).DecodeString(src)
	if err != nil {
		return "", err
	}
	dst, err := self.Decrypt(plaintext, []byte(key))
	return string(dst), nil
}
