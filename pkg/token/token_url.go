package token

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// UrlSign
// path 请求的路径 (不附带 querystring)
func (t *token) UrlSign(timestamp int64, path string, method string, params url.Values) (tokenString string, err error) {
	// 合法的 Methods
	methods := map[string]bool{
		"get":     true,
		"post":    true,
		"put":     true,
		"path":    true,
		"delete":  true,
		"head":    true,
		"options": true,
	}

	methodName := strings.ToLower(method)
	if !methods[methodName] {
		err = errors.New("method param error")
		return
	}

	// Encode() 方法中自带 sorted by key
	sortParamsEncode := params.Encode()

	// 加密字符串规则 path + method + sortParamsEncode + secret
	encryptStr := fmt.Sprintf("%s%s%s%d%s", path, methodName, sortParamsEncode, timestamp, t.secret)

	// 对加密字符串进行 md5
	s := md5.New()
	s.Write([]byte(encryptStr))
	md5Str := hex.EncodeToString(s.Sum(nil))

	tokenString = md5Str
	return
}
