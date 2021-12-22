package token

import (
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
)

var _ Token = (*token)(nil)

// Token 实现 接口令牌 的封装
// JWT 令牌 （代替 Session）
// Sign 签名防篡改
type Token interface {
	// i 为了避免被其他包实现
	i()

	// JwtSign 签名
	JwtSign(userId int64, userName string, expireDuration time.Duration) (tokenString string, err error)

	// JwtParse 解密
	JwtParse(tokenString string) (*claims, error)

	// UrlSign URL 签名
	// 防参数篡改，防重放攻击
	UrlSign(timestamp int64, path string, method string, params url.Values) (tokenString string, err error)
}

type token struct {
	secret string
}

type claims struct {
	UserID   int64
	UserName string
	jwt.StandardClaims
}

// New 根据 secret 生成后续的签名
func New(secret string) Token {
	return &token{
		secret: secret,
	}
}

func (t *token) i() {}
