package encrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSalt(t *testing.T) {
	for _, e := range []Encoding{
		NewHexEncoding(),
		NewBase64Encoding(),
		NewBase32Encoding(),
		NewBase32Human(),
	} {
		s := Salt(e)
		t.Log(s)
		assert.NotZero(t, s)
	}
}

func TestNonce(t *testing.T) {
	for _, e := range []int{5, 50, 5000} {
		n1 := Nonce(e)
		assert.Len(t, n1, e)
	}
}
