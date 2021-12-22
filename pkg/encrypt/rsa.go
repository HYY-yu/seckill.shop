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
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

// RSA 加解密
// SHA256 or SHA128

type RSAKeyLen int

const (
	// RSAAlgorithmSignSHA256 RSA签名算法 SHA256
	RSAAlgorithmSignSHA256 = crypto.SHA256
	// RSAAlgorithmSignSHA1 RSA签名算法 SHA1
	RSAAlgorithmSignSHA1 = crypto.SHA1

	RSA1024 RSAKeyLen = 1024
	RSA2048 RSAKeyLen = 2048
	RSA4096 RSAKeyLen = 4096
)

// GenerateRSAFile 生成密钥对文件
// pubKeyFilename: 公钥文件名
// priKeyFilename: 私钥文件名
// keyLength: 密钥长度
func GenerateRSAFile(pubKeyFilename, priKeyFilename string, keyLength RSAKeyLen) error {
	if pubKeyFilename == "" {
		pubKeyFilename = "id_rsa.pub"
	}
	if priKeyFilename == "" {
		priKeyFilename = "id_rsa"
	}

	if keyLength == 0 {
		keyLength = RSA1024
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

// GenerateRSAString 生成密钥对字符串
// keyLength 密钥的长度
// 生成 公钥、私钥
func GenerateRSAString(keyLength RSAKeyLen) (string, string, error) {
	if keyLength == 0 || keyLength < 1024 {
		keyLength = 1024
	}

	bufPub := make([]byte, 0, 1024*5)
	pubBuffer := bytes.NewBuffer(bufPub)

	bufPri := make([]byte, 0, 1024*5)
	priBuffer := bytes.NewBuffer(bufPri)

	err := writeRSAKeyPair(pubBuffer, priBuffer, keyLength)
	if err != nil {
		return "", "", nil
	}
	return pubBuffer.String(), priBuffer.String(), nil
}

// writeRSAKeyPair 生成RSA密钥对
func writeRSAKeyPair(publicKeyWriter, privateKeyWriter io.Writer, keyLength RSAKeyLen) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, int(keyLength))
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

// readRSAKeyPairFromFile 从文件读取RSA密钥对
// pubKeyFilename: 公钥文件名称   priKeyFilename: 私钥文件名
func readRSAKeyPairFromFile(pubKeyFilename, priKeyFilename string) ([]byte, []byte, error) {
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
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// NewGoRSA 初始化 GORSA ,读取公钥、私钥，指定Hash算法
func NewGoRSA(publicKey, privateKey []byte) (*GoRSA, error) {
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
			PublicKey:  pub,
			PrivateKey: pri,
		}, nil
	}
	return nil, errors.New("private key not supported")
}

// NewGoRSAFromFile 初始化 GoRSA 从文件中加载秘钥
func NewGoRSAFromFile(pubKeyFilename, priKeyFilename string) (*GoRSA, error) {
	publicKey, privateKey, err := readRSAKeyPairFromFile(pubKeyFilename, priKeyFilename)
	if err != nil {
		return nil, err
	}
	return NewGoRSA(publicKey, privateKey)
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

// PublicEncryptWithEncoding 公钥加密，并转化成 Encoding 编码的字符串
func (r *GoRSA) PublicEncryptWithEncoding(data []byte, e Encoding) (string, error) {
	p, err := r.PublicEncrypt(data)
	if err != nil {
		return "", err
	}
	return e.EncodeToString(p), nil
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

// PrivateDecryptWithEncoding 先用 Encoding 解码 data，再进行解密
func (r *GoRSA) PrivateDecryptWithEncoding(data string, e Encoding) ([]byte, error) {
	encryptData, err := e.DecodeString(data)
	if err != nil {
		return nil, err
	}
	pd, err := r.PrivateDecrypt(encryptData)
	if err != nil {
		return nil, err
	}
	return pd, nil
}

// Sign 利用私钥对数据进行签名
func (r *GoRSA) Sign(signAlgorithm crypto.Hash, encoder Encoding, data string) (string, error) {
	h := signAlgorithm.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, signAlgorithm, hashed)
	if err != nil {
		return "", err
	}
	return encoder.EncodeToString(sign), err
}

// Verify 利用公钥验证数据签名
func (r *GoRSA) Verify(signAlgorithm crypto.Hash, decoder Encoding, data string, sign string) error {
	h := signAlgorithm.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	decodedSign, err := decoder.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(r.PublicKey, signAlgorithm, hashed, decodedSign)
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
