package encrypt

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// MD5 计算md5的封装
// 32位小写
func MD5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

// SHA1 封装 sha1 算法
func SHA1(s string) string {
	return SHA1WithEncoding(s, NewHexEncoding())
}

// SHA1WithEncoding SHA1 算法，可选数据编码格式
func SHA1WithEncoding(s string, e Encoding) string {
	m := sha1.New()
	m.Write([]byte(s))
	return e.EncodeToString(m.Sum(nil))
}

// SHA1FileHash 计算文件的 sha1 值
func SHA1FileHash(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SHA256 封装 sha256 算法
func SHA256(s string) string {
	return SHA256WithEncoding(s, NewHexEncoding())
}

// SHA256WithEncoding SHA256 算法，可选数据编码格式
func SHA256WithEncoding(s string, e Encoding) string {
	m := sha256.New()
	m.Write([]byte(s))
	return e.EncodeToString(m.Sum(nil))
}

// Salt 利用安全随机数生成器生成盐值
func Salt(e Encoding) string {
	nonce := make([]byte, 4)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	//             8字节   4字节
	// base64.Raw  11字    6字
	// base64.Std  12字    8字
	// base32.Raw  13字    7字
	// base32.Std  16字    8字
	// hex         16字    8字
	return e.EncodeToString(nonce)
}

// Nonce 生成一个n位 ，只包含小写字母+数字的随机字符串
func Nonce(n int) string {
	nonce := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	return NewHexEncoding().EncodeToString(nonce)[:n]
}
