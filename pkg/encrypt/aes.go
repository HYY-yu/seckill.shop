package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// AES 加密解密
// ( AES-128/192/256 )
// ( CBC - / PKCS5 / 偏移量 block 0 )
// ( ECB - / PKCS5 / )

const (
	AES128 = 16
	AES192 = 24
	AES256 = 32

	ECB = 1
	CBC = 2
)

type GoAES struct {
	key []byte
}

// type :AES-128 AES-192 AES-256
func NewGoAES(key string, aestype int) *GoAES {
	return &GoAES{
		key: paddingKey(key, aestype),
	}
}

func NewGoAESWith(key []byte) *GoAES {
	return &GoAES{
		key: key,
	}
}

// 加密 Base64
func (self *GoAES) EnBase64(src string, model int) (result string, err error) {
	var re []byte
	if model == ECB {
		re, err = self.ECBEncrypt([]byte(src))
	} else if model == CBC {
		re, err = self.CBCEncrypt([]byte(src))
	} else {
		err = errors.New("不支持的 model")
	}

	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(re), nil
}

// 解密 Base64
func (self *GoAES) UnBase64(src string, model int) (string, error) {
	result, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	var origData []byte
	if model == ECB {
		origData, err = self.ECBDecrypt(result)
	} else if model == CBC {
		origData, err = self.CBCDecrypt(result)
	} else {
		err = errors.New("不支持的 model")
	}

	if err != nil {
		return "", err
	}
	return string(origData), nil
}

func (self *GoAES) CBCEncrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(self.key)
	if err != nil {
		return nil, err
	}

	plaintext = pkcs5Padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

func (self *GoAES) CBCDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(self.key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(ciphertext, ciphertext)
	return pkcs5UnPadding(ciphertext), nil
}

func (self *GoAES) CBCDecryptWith(ciphertext []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(self.key)
	if err != nil {
		return nil, err
	}

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(ciphertext, ciphertext)
	return pkcs5UnPadding(ciphertext), nil
}

func (self *GoAES) ECBEncrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(self.key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plaintext = pkcs5Padding(plaintext, blockSize)
	ciphertext := make([]byte, 0)
	temp := make([]byte, aes.BlockSize)

	for i := 0; i < len(plaintext); i += aes.BlockSize {
		block.Encrypt(temp, plaintext[i:i+aes.BlockSize])
		ciphertext = append(ciphertext, temp...)
	}
	return ciphertext, nil
}

func (self *GoAES) ECBDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(self.key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, 0)
	temp := make([]byte, aes.BlockSize)

	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(temp, ciphertext[i:i+aes.BlockSize])
		plaintext = append(plaintext, temp...)
	}
	return pkcs5UnPadding(plaintext), nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	unpadding := int(origData[len(origData)-1])
	return origData[:(len(origData) - unpadding)]
}

// 密码不够blockSize位则填充 0  // 然后返回 blockSize 位密码
func paddingKey(key string, blockSize int) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(key)
	for i := len(key); i < blockSize; i++ {
		buffer.WriteString("0")
	}
	return buffer.Bytes()[:blockSize]
}
