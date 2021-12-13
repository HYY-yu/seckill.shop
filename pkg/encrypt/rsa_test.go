package encrypt

import (
	"fmt"
	"testing"
)

func TestNewRSAString(t *testing.T) {
	pub, pri, _ := NewRSAString(1024)
	fmt.Println("pub", pub)
	fmt.Println("pri", pri)
}

var pubFile = "key/rsa_2048.pub"
var priFile = "key/rsa_2048"

var test = `<?xml version="1.0" encoding="UTF-8" ?>
<mo version="1.0.0">
  <head>
    <businessType>PERSONAL_GOODS_DECLAR</businessType>
  </head>
  <body>
    <jkfSign>
      <companyCode>WC14060601</companyCode>
      <businessNo>LP057</businessNo>
    </jkfSign>
    <jkfGoodsDeclar>
      <personalGoodsFormNo>CFAK00022900100057</personalGoodsFormNo>
      <approveResult>51</approveResult>
      <approveComment>人工放行</approveComment>
      <processTime>2014-09-29 13:31:20</processTime>
    </jkfGoodsDeclar>
  </body>
</mo>`

func TestNewRSAFile(t *testing.T) {
	NewRSAFile(pubFile, priFile, 2048)
}

func TestReadRSAKeyPairFromFile(t *testing.T) {
	gorsa, err := NewGoRSAFromFile(pubFile, priFile, RSAAlgorithmSignSHA1)
	if err != nil {
		fmt.Println(err)
	}

	s, err := gorsa.Sign(test)
	fmt.Println(s)

	e := gorsa.Verify(test, "XHin4uUFqrKDEhKBD/hQisXLFFSxM6EZCvCPqnWCQJq3uEp3ayxmFuUgVE0Xoh4AIWjIIsOWdnaToL1bXvAFKwjCtXnkaRwUpvWrk+Q0eqwsoAdywsVQDEceG5stas1CkPtrznAIW2eBGXCWspOj+aumEAcPyYDxLhDN646Krzw=")
	fmt.Println(e)
}
