package encrypt

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func Test_paddingKey(t *testing.T) {
	x := paddingKey("123456789abcdefghi", AES128)
	fmt.Println(string(x))
	x = paddingKey("12345678", AES192)
	fmt.Println(string(x))
	x = paddingKey("12345678", AES256)
	fmt.Println(string(x))
}

func TestGoAES_UnBase64(t *testing.T) {
	key, err := base64.StdEncoding.DecodeString("qZe60QZFxuirub2ey4+7+Q==")
	if err != nil {
		fmt.Println(err)
	}

	goaes := NewGoAES(string(key), AES128)
	r, err := goaes.UnBase64("8O2O2SsToX9q8anImAxh7Q==",
		ECB)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r)
}

func TestGoAES_EnBase64(t *testing.T) {
	key, err := base64.StdEncoding.DecodeString("qZe60QZFxuirub2ey4+7+Q==")
	if err != nil {
		fmt.Println(err)
	}

	goaes := NewGoAES(string(key), AES192)

	r, err := goaes.EnBase64("3sCmQNfp4yu", ECB)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
}
