package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/scrypt"
)

// AES 加密解密封装
// 示例：
//    ciphertext,err := NewGoAES("Some Password",AES128).
//        WithModel(ECB).
//        WithEncoding(NewBase64Encoding()).
//        WithIv(iv).
//        Encrypt(src)
// 默认情况下：
// ciphertext,err := NewGoAES("Some Password",AES128).
//         Encrypt(src)
// 默认为： CBC模式 Base64编码 PKCS5 自动生成iv

// AESModel AES 的不同加密模式
type AESModel int

// AESKeyLen AES 的密钥长度
type AESKeyLen int

const (
	// AES128 16位密钥
	AES128 AESKeyLen = 16
	// AES192 24位密钥
	AES192 AESKeyLen = 24
	// AES256 32位密钥
	AES256 AESKeyLen = 32

	// ECB 加密模式
	// 明文对应唯一密文，相同的明文会被转化成相同的密文
	ECB AESModel = 1
	// CBC 加密模式
	// 明文不加密出唯一密文
	// 只能串行化，速度慢
	// 密文长
	CBC AESModel = 2
	// CTR 加密模式
	// 流加密，速度快，无填充
	// 密文较短
	CTR AESModel = 3
)

// GoAES AES 加密封装
type GoAES struct {
	key     []byte
	model   AESModel
	iv      []byte
	setIv   bool
	encoder Encoding
}

// NewGoAES 新建一个 GoAES 对象
// key 用于后期加、解密的密钥
// aesLen 当key的长度不够\超过，会自动用0填充\根据 aesLen 截断
// 注意：这种密码填充模式是不安全的，若要求安全的密码生成，请使用 NewGoAESSafety
func NewGoAES(key string, aesLen AESKeyLen) *GoAES {
	return &GoAES{
		key: paddingKey(key, aesLen),
	}
}

// NewGoAESSafety 此方法面向安全要求严格的应用
// 利用了 scrypt 包计算密钥，salt 可以使用 encrypt.Salt 生成
func NewGoAESSafety(key string, salt string, aesLen AESKeyLen) (a *GoAES, realKey []byte, err error) {
	realKey, err = scrypt.Key([]byte(key), []byte(salt), 32768, 8, 1, int(aesLen))
	if err != nil {
		return
	}

	return &GoAES{
		key: realKey,
	}, realKey, nil
}

// WithModel 设置 AESModel
// 不设置，默认为ECB
func (a *GoAES) WithModel(m AESModel) *GoAES {
	a.model = m
	return a
}

// WithEncoding
// 不设置，默认为 Base64Encoding
func (a *GoAES) WithEncoding(e Encoding) *GoAES {
	a.encoder = e
	return a
}

// WithIv
// 不设置，默认为空（或者自动生成iv）
func (a *GoAES) WithIv(iv []byte) *GoAES {
	a.iv = iv
	a.setIv = true
	return a
}

func (a *GoAES) initConfig() {
	if a.model == 0 {
		a.model = ECB
	}

	if a.encoder == nil {
		a.encoder = NewBase64Encoding()
	}
}

// Encrypt 加密 plaintext
func (a *GoAES) Encrypt(plaintext string) (string, error) {
	a.initConfig()
	plainBytes := []byte(plaintext)

	var ciphertext []byte
	var err error

	switch a.model {
	case ECB:
		// 此模式不需要 IV
		ciphertext, err = a.ecbEncrypt(plainBytes)
		if err != nil {
			return "", nil
		}
	case CBC:
		a.blockIvEncrypt()
		ciphertext, err = a.cbcEncrypt(plainBytes)
		if err != nil {
			return "", nil
		}
		if !a.setIv {
			// 默认使用 ciphertext 的第一个 block 作为 iv
			ciphertext = append(a.iv, ciphertext...)
		}
	case CTR:
		a.blockIvEncrypt()
		ciphertext, err = a.ctrEncrypt(plainBytes)
		if err != nil {
			return "", nil
		}
		if !a.setIv {
			// 默认使用 ciphertext 的第一个 block 作为 iv
			ciphertext = append(a.iv, ciphertext...)
		}
	}

	return a.encoder.EncodeToString(ciphertext), nil
}

// DecryptBytes 解密 , 和  Decrypt 的区别是，不会用 encoder 处理 ciphertext
// 直接解密 ciphertext
func (a *GoAES) DecryptBytes(ciphertext []byte) (string, error) {
	var plainBytes []byte
	var err error

	switch a.model {
	case ECB:
		// 此模式不需要 IV
		plainBytes, err = a.ecbDecrypt(ciphertext)
		if err != nil {
			return "", nil
		}
	case CBC:
		if !a.setIv {
			// 默认使用 ciphertext 的第一个 block 作为 iv
			ciphertext, a.iv = a.blockIvDecrypt(ciphertext)
		}
		plainBytes, err = a.cbcDecrypt(ciphertext)
		if err != nil {
			return "", nil
		}
	case CTR:
		if !a.setIv {
			// 默认使用 ciphertext 的第一个 block 作为 iv
			ciphertext, a.iv = a.blockIvDecrypt(ciphertext)
		}
		plainBytes, err = a.ctrDecrypt(ciphertext)
		if err != nil {
			return "", nil
		}
	}

	return string(plainBytes), nil
}

// Decrypt 解密
func (a *GoAES) Decrypt(ciphertext string) (string, error) {
	a.initConfig()
	cipherBytes, err := a.encoder.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	return a.DecryptBytes(cipherBytes)
}

func (a *GoAES) blockIvEncrypt() {
	if len(a.iv) == 0 {
		iv := make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			panic(err)
		}
		a.iv = iv
	}
}

func (a *GoAES) blockIvDecrypt(ciphertext []byte) (ct []byte, iv []byte) {
	if len(a.iv) == 0 {
		iv = ciphertext[:aes.BlockSize]
		ct = ciphertext[aes.BlockSize:]
		return
	}
	return ciphertext, a.iv
}

func (a *GoAES) ctrEncrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plaintext))

	stream := cipher.NewCTR(block, a.iv)
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}

func (a *GoAES) ctrDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, a.iv)
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

func (a *GoAES) cbcEncrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	plaintext = pkcs5Padding(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))

	cipher.NewCBCEncrypter(block, a.iv).CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func (a *GoAES) cbcDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	cipher.NewCBCDecrypter(block, a.iv).CryptBlocks(ciphertext, ciphertext)
	return pkcs5UnPadding(ciphertext), nil
}

func (a *GoAES) ecbEncrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
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

func (a *GoAES) ecbDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
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
	pt := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, pt...)
}

func pkcs5UnPadding(origData []byte) []byte {
	unpadding := int(origData[len(origData)-1])
	return origData[:(len(origData) - unpadding)]
}

// paddingKey 密码不够blockSize位则填充 0
// 然后返回 blockSize 位密码
func paddingKey(key string, blockSize AESKeyLen) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(key)
	for i := len(key); i < int(blockSize); i++ {
		buffer.WriteString("0")
	}
	return buffer.Bytes()[:blockSize]
}
