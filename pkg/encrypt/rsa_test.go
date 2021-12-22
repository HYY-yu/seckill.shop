package encrypt

import (
	"crypto"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRSAString(t *testing.T) {
	for _, e := range []RSAKeyLen{RSA1024, RSA2048, RSA4096} {
		pubString, priString, err := GenerateRSAString(e)
		assert.NoError(t, err)

		_, err = NewGoRSA([]byte(pubString), []byte(priString))
		assert.NoError(t, err)
	}
}

func TestGenerateRSAFile(t *testing.T) {
	for _, e := range []RSAKeyLen{RSA1024, RSA2048, RSA4096} {
		pubFilePath := fmt.Sprintf("./testdata/test_rsa%d.pub", e)
		priFilePath := fmt.Sprintf("./testdata/rsa%d", e)

		err := GenerateRSAFile(pubFilePath, priFilePath, e)
		assert.NoError(t, err)

		_, err = NewGoRSAFromFile(pubFilePath, priFilePath)
		assert.NoError(t, err)
	}
}

func TestRSAEncrypt(t *testing.T) {
	t.Run("test encrypt and decrypt ", func(t *testing.T) {
		pubFilePath := fmt.Sprintf("./testdata/test_rsa%d.pub", RSA1024)
		priFilePath := fmt.Sprintf("./testdata/rsa%d", RSA1024)

		goRSA, err := NewGoRSAFromFile(pubFilePath, priFilePath)
		assert.NoError(t, err)
		plaintext := `"please encrypt me !"`

		encryptedData, err := goRSA.PublicEncryptWithEncoding([]byte(plaintext), NewBase64Encoding())
		assert.NoError(t, err)

		decryptdData, err := goRSA.PrivateDecryptWithEncoding(encryptedData, NewBase64Encoding())
		assert.NoError(t, err)

		assert.Equal(t, decryptdData, []byte(plaintext))
	})

	t.Run("test encrypt and wrong decrypt ", func(t *testing.T) {
		pubFilePath := fmt.Sprintf("./testdata/test_rsa%d.pub", RSA1024)
		priFilePath := fmt.Sprintf("./testdata/rsa%d", RSA2048)

		goRSA, err := NewGoRSAFromFile(pubFilePath, priFilePath)
		assert.NoError(t, err)
		plaintext := `"please encrypt me !"`

		encryptedData, err := goRSA.PublicEncryptWithEncoding([]byte(plaintext), NewBase64Encoding())
		assert.NoError(t, err)

		decryptdData, err := goRSA.PrivateDecryptWithEncoding(encryptedData, NewBase64Encoding())
		assert.Error(t, err)
		assert.Zero(t, decryptdData)
	})

	t.Run("test different encoding ", func(t *testing.T) {
		pubFilePath := fmt.Sprintf("./testdata/test_rsa%d.pub", RSA1024)
		priFilePath := fmt.Sprintf("./testdata/rsa%d", RSA1024)

		goRSA, err := NewGoRSAFromFile(pubFilePath, priFilePath)
		assert.NoError(t, err)
		plaintext := `"please encrypt me !"`

		encryptedData, err := goRSA.PublicEncryptWithEncoding([]byte(plaintext), NewHexEncoding())
		assert.NoError(t, err)

		decryptdData, err := goRSA.PrivateDecryptWithEncoding(encryptedData, NewBase64Encoding())
		assert.Error(t, err)
		assert.Zero(t, decryptdData)
	})
}

func TestRSASign(t *testing.T) {
	pubFilePath := fmt.Sprintf("./testdata/test_rsa%d.pub", RSA1024)
	priFilePath := fmt.Sprintf("./testdata/rsa%d", RSA1024)
	goRSA, err := NewGoRSAFromFile(pubFilePath, priFilePath)
	assert.NoError(t, err)

	for _, e := range []crypto.Hash{RSAAlgorithmSignSHA1, RSAAlgorithmSignSHA256} {
		signed, err := goRSA.Sign(e, NewHexEncoding(), "sign me thanks")
		assert.NoError(t, err)

		err = goRSA.Verify(e, NewHexEncoding(), "sign me thanks", signed)
		assert.NoError(t, err)
	}
}

var pubString = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDHc+PP8LuTlBL1zCX+lh9kcur
gHHIXFnV/tDK789DaJuhwZvQ1lu5Zdcn+ULbNUKkB6b5tCP0sZxlpoCVKMyKHtde
h/YGXwBD8sMc+XcRs0eh3/tyr4EoBu3bomzHWDGmHjH/F5GotFTrGcB6xQwAROy4
mT5SketlQ3c7tucI+QIDAQAB
-----END PUBLIC KEY-----
`

var priString = `
-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAMMdz48/wu5OUEvX
MJf6WH2Ry6uAcchcWdX+0Mrvz0Nom6HBm9DWW7ll1yf5Qts1QqQHpvm0I/SxnGWm
gJUozIoe116H9gZfAEPywxz5dxGzR6Hf+3KvgSgG7duibMdYMaYeMf8Xkai0VOsZ
wHrFDABE7LiZPlKR62VDdzu25wj5AgMBAAECgYBKcdxYrp5EaHLwjNlIk0ciGfeY
pvhC1yGbqY6mb1soQAhpbkJyKudyVG4EHXGpy6dyiEzoJxg063NdwWp7/sYTHk/N
13UzGTudIKuNacnJk0WKu4owQticC71ZIqUjSZgN0CiEKQ6YfoGOFTzeMqzVYQjI
mPzGdLK74y3YYlmigQJBAObzhzYlWjOypx3klmrPTu2KXPg3ATTEB8kN/isY8bYu
ikVdd2yUd0AvaC7PPwEEjGmsSrEeXw1tsVfZ8VkBaikCQQDYR0+8VzGLdgIFQc/6
+IY5fQlEt/Hc7qsi7JT4o+f+BGJlAT7+OeDMThavKdWq1UvZDyCKdtYRfxQ1jj7D
4yJRAkBrG6InkGcm9sHecTb5Ti+ypqq7Svc6O3fI3L51ylm/PhJOXSyXpLsxf0r3
+pGjrTJZh9gUEJvQpIDM13zA5JERAkBI2zTsED9baIRjuvjR5Xhp00oVARYTw76Y
xDOm0qgq9NUki1fqEhs9F60ikqgspS+oziS7IC8as8FeDS3tlQ0RAkA5OdDvhQRQ
PI75ULyHazTEm4Rak8TKmKl64pmnwcw4GS9fKWs7jRAuem1OtwA8HAqjaDeXC8Cd
6fDfq7z5bZnE
-----END PRIVATE KEY-----
`

func TestRSADecrypt(t *testing.T) {
	goRSA, err := NewGoRSA([]byte(pubString), []byte(priString))
	assert.NoError(t, err)

	outEncryptByPublicKey := `
rdV+gLugw6bOEUAHozggT33lT7wVUzySFs0ZYALtPDuM6gzZcOeOEc2rT/EfOdmx76uvmj8BlgQBFfXl73TiYvlIi5UbgXsWGaSAbtZGCpziePgsMKvQGexLaV5Ey9PXUdC9nO9xi94YLb8tkMxqVtsKNvb6//pS9pJ4WpXrMfk=
`
	plaintext, err := goRSA.PrivateDecryptWithEncoding(outEncryptByPublicKey, NewBase64Encoding())
	assert.NoError(t, err)

	t.Log(string(plaintext))
}

func TestRSAVerify(t *testing.T) {
	goRSA, err := NewGoRSA([]byte(pubString), []byte(priString))
	assert.NoError(t, err)

	outSign := `
d0IoexR4htGwHLPxhnU016fY/reO3okoh8S/j6YNWXXupmVx85ZjXBMZX9Xwhw2xp+2Yw5DUzd8G/LmOH3guvC5o2f9gf+kwMyM6aXbodvvj+6sRVTRrHqmv0QEYzan5cyJbsJrLcDWi4UVmcsdmPyRYw9SKsBx9AwPq7kGOlfA=
`

	err = goRSA.Verify(RSAAlgorithmSignSHA1, NewBase64Encoding(), "sign me thanks", outSign)
	assert.NoError(t, err)
}
