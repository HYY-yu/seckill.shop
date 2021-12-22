package encrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGoAES(t *testing.T) {
	for _, e := range []AESKeyLen{AES128, AES192, AES256} {
		goAes := NewGoAES("low_key", e)
		assert.Len(t, goAes.key, int(e))

		goAes2 := NewGoAES("very loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong key", e)
		assert.Len(t, goAes2.key, int(e))
	}
}

func TestNewGoAESSafety(t *testing.T) {
	for _, e := range []AESKeyLen{AES128, AES192, AES256} {
		goAes, _, err := NewGoAESSafety("low_key", "salt", e)
		assert.NoError(t, err)

		assert.Len(t, goAes.key, int(e))
	}
}

func TestGoAES_Encrypt(t *testing.T) {
	t.Run("without iv base64", func(t *testing.T) {
		goAes := NewGoAES("low_key", AES128)
		ecbEncrypt := `PzeFufy3wliMgicGwGc+H1+49uDFERi/23wtbhpBw7A=`

		for _, e := range []AESModel{ECB, CBC, CTR} {
			goAes.WithModel(e)
			goAes.WithEncoding(NewBase64Encoding())

			encryptData, err := goAes.Encrypt("please encrypt me ~ ")
			assert.NoError(t, err)
			switch e {
			case ECB:
				assert.Equal(t, encryptData, ecbEncrypt)
			case CBC:
				t.Log(encryptData, " AESModel: CBC")
			case CTR:
				t.Log(encryptData, " AESModel: CTR")
			}
		}
	})

	t.Run("with iv base64", func(t *testing.T) {
		goAes := NewGoAES("low_key", AES128)
		iv := []byte(`this is 12321123`) // the len must eq aes.BlockSize
		ecbEncrypt := `PzeFufy3wliMgicGwGc+H1+49uDFERi/23wtbhpBw7A=`
		cbcEncrypt := `sqqvIjJzBjHerY84liwe5pQhg5Pqa90Zoz18/malLx0=`
		ctrEncrypt := `JGQ1/2BrDT8mWEOEGXVOPT9Kibc=`

		for _, e := range []AESModel{ECB, CBC, CTR} {
			goAes.WithModel(e)
			goAes.WithIv(iv)
			goAes.WithEncoding(NewBase64Encoding())

			encryptData, err := goAes.Encrypt("please encrypt me ~ ")
			assert.NoError(t, err)
			switch e {
			case ECB:
				assert.Equal(t, encryptData, ecbEncrypt)
			case CBC:
				assert.Equal(t, encryptData, cbcEncrypt)
			case CTR:
				assert.Equal(t, encryptData, ctrEncrypt)
			}
		}
	})
}

func TestGoAES_Decrypt(t *testing.T) {
	t.Run("without iv base64", func(t *testing.T) {
		goAes := NewGoAES("low_key", AES128)
		plaintext := `please encrypt me ~ `

		for _, e := range []AESModel{ECB, CBC, CTR} {
			goAes.iv = nil
			goAes.WithModel(e)
			goAes.WithEncoding(NewBase64Encoding())
			ciphtData := ""

			switch e {
			case ECB:
				ciphtData = "PzeFufy3wliMgicGwGc+H1+49uDFERi/23wtbhpBw7A="
			case CBC:
				ciphtData = "ywZWaANYdxeIVmiftnb8ajoBK5nbnMZlFzlxLOhnqZyF6sDkBbLWL/DWEknFBfxc"
			case CTR:
				ciphtData = "Wd6RreN2zd5M9QbfxFT/o/YE7O3YMNwiR377xgJsha4wYrx7"
			}
			decryptData, err := goAes.Decrypt(ciphtData)
			assert.NoError(t, err)
			assert.Equal(t, decryptData, plaintext)
			t.Log("Success", decryptData)
		}
	})

	t.Run("with iv base64", func(t *testing.T) {
		goAes := NewGoAES("low_key", AES128)
		plaintext := `please encrypt me ~ `
		iv := []byte(`this is 12321123`) // the len must eq aes.BlockSize

		for _, e := range []AESModel{ECB, CBC, CTR} {
			goAes.WithIv(iv)
			goAes.WithModel(e)
			goAes.WithEncoding(NewBase64Encoding())
			ciphtData := ""

			switch e {
			case ECB:
				ciphtData = "PzeFufy3wliMgicGwGc+H1+49uDFERi/23wtbhpBw7A="
			case CBC:
				ciphtData = "sqqvIjJzBjHerY84liwe5pQhg5Pqa90Zoz18/malLx0="
			case CTR:
				ciphtData = "JGQ1/2BrDT8mWEOEGXVOPT9Kibc="
			}
			decryptData, err := goAes.Decrypt(ciphtData)
			assert.NoError(t, err)
			assert.Equal(t, decryptData, plaintext)
			t.Log("Success", decryptData)
		}
	})
}
