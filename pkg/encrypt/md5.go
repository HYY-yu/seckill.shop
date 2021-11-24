package encrypt

import (
	cryptoMD5 "crypto/md5"
	"encoding/hex"
)

type md5 struct{}

func New() *md5 {
	return &md5{}
}

func (m *md5) MD5(encryptStr string) string {
	s := cryptoMD5.New()
	s.Write([]byte(encryptStr))
	return hex.EncodeToString(s.Sum(nil))
}
