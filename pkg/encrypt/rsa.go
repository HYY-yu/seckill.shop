// The MIT License (MIT)
//
// Copyright (c) 2017 雷纳科斯
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// https://github.com/chanyipiaomiao/hltool/blob/master/rsa.go

package encrypt

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

// RSA 加解密
// SHA256 or SHA128
// PKCS8

const (
	// RSAAlgorithmSign RSA签名算法
	RSAAlgorithmSignSHA256 = crypto.SHA256
	RSAAlgorithmSignSHA1   = crypto.SHA1
)

// NewRSAFile 生成密钥对文件
// pubKeyFilename: 公钥文件名 priKeyFilename: 私钥文件名 keyLength: 密钥长度
func NewRSAFile(pubKeyFilename, priKeyFilename string, keyLength int) error {
	if pubKeyFilename == "" {
		pubKeyFilename = "id_rsa.pub"
	}
	if priKeyFilename == "" {
		priKeyFilename = "id_rsa"
	}

	if keyLength == 0 || keyLength < 1024 {
		keyLength = 1024
	}

	// 创建公钥文件
	pubWriter, err := os.Create(pubKeyFilename)
	if err != nil {
		return err
	}
	defer pubWriter.Close()

	// 创建私钥文件
	priWriter, err := os.Create(priKeyFilename)
	if err != nil {
		return err
	}
	defer priWriter.Close()

	// 生成密钥对
	err = writeRSAKeyPair(pubWriter, priWriter, keyLength)
	if err != nil {
		return err
	}
	return nil
}

// NewRSAString 生成密钥对字符串
// keyLength 密钥的长度
func NewRSAString(keyLength int) (string, string, error) {
	if keyLength == 0 || keyLength < 1024 {
		keyLength = 1024
	}

	bufPub := make([]byte, 1024*5)
	pubuffer := bytes.NewBuffer(bufPub)

	bufPri := make([]byte, 1024*5)
	pribuffer := bytes.NewBuffer(bufPri)

	err := writeRSAKeyPair(pubuffer, pribuffer, keyLength)
	if err != nil {
		return "", "", nil
	}
	return pubuffer.String(), pribuffer.String(), nil
}

// writeRSAKeyPair 生成RSA密钥对
func writeRSAKeyPair(publicKeyWriter, privateKeyWriter io.Writer, keyLength int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}

	derStream, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(privateKeyWriter, block)
	if err != nil {
		return err
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)

	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	err = pem.Encode(publicKeyWriter, block)
	if err != nil {
		return err
	}

	return nil
}

// ReadRSAKeyPairFromFile 从文件读取RSA密钥对
// pubKeyFilename: 公钥文件名称   priKeyFilename: 私钥文件名
func ReadRSAKeyPairFromFile(pubKeyFilename, priKeyFilename string) ([]byte, []byte, error) {
	pub, err := ioutil.ReadFile(pubKeyFilename)
	if err != nil {
		return nil, nil, err
	}

	pri, err := ioutil.ReadFile(priKeyFilename)
	if err != nil {
		return nil, nil, err
	}
	return pub, pri, nil
}

// GoRSA RSA加密解密
type GoRSA struct {
	PublicKey        *rsa.PublicKey
	PrivateKey       *rsa.PrivateKey
	RSAAlgorithmSign crypto.Hash
}

func NewGoRSA(publicKey, privateKey []byte, sign crypto.Hash) (*GoRSA, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}

	priInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pri, ok := priInterface.(*rsa.PrivateKey)
	if ok {
		return &GoRSA{
			PublicKey:        pub,
			PrivateKey:       pri,
			RSAAlgorithmSign: sign,
		}, nil
	}
	return nil, errors.New("private key not supported")
}

// NewGoRSAFromFile 初始化 GoRSA对象
func NewGoRSAFromFile(pubKeyFilename, priKeyFilename string, sign crypto.Hash) (*GoRSA, error) {
	publicKey, privateKey, err := ReadRSAKeyPairFromFile(pubKeyFilename, priKeyFilename)
	if err != nil {
		return nil, err
	}
	return NewGoRSA(publicKey, privateKey, sign)
}

// PublicEncrypt 公钥加密
func (r *GoRSA) PublicEncrypt(data []byte) ([]byte, error) {
	partLen := r.PublicKey.N.BitLen()/8 - 11
	chunks := split(data, partLen)
	buffer := bytes.NewBufferString("")

	for _, chunk := range chunks {
		b, err := rsa.EncryptPKCS1v15(rand.Reader, r.PublicKey, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(b)
	}

	return buffer.Bytes(), nil
}

func (r *GoRSA) PublicEncryptBase64(data []byte) (string, error) {
	p, err := r.PublicEncrypt(data)
	if err != nil {
		return "", err
	}
	base64string := base64.RawURLEncoding.EncodeToString(p)
	return base64string, nil
}

// PrivateDecrypt 私钥解密
func (r *GoRSA) PrivateDecrypt(encrypted []byte) ([]byte, error) {
	partLen := r.PublicKey.N.BitLen() / 8
	chunks := split(encrypted, partLen)
	buffer := bytes.NewBufferString("")

	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.PrivateKey, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}
	return buffer.Bytes(), nil
}

// Sign 数据进行签名
func (r *GoRSA) Sign(data string) (string, error) {
	h := r.RSAAlgorithmSign.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, r.RSAAlgorithmSign, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), err
}

// Verify 数据验证签名
func (r *GoRSA) Verify(data string, sign string) error {
	h := r.RSAAlgorithmSign.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	decodedSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(r.PublicKey, r.RSAAlgorithmSign, hashed, decodedSign)
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}
