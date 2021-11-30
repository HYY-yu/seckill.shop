package encrypt

import (
	"fmt"
	"log"
	"testing"
)

func TestGoRC4_DecryptHex(t *testing.T) {
	temp, err := NewGoRC4().DecryptHex("f73989b3fd94e35b563a07", "Gec2")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(temp)
}

func TestGoRC4_DeBase64(t *testing.T) {
	temp, err := NewGoRC4().DeBase32Ma("URDWNFPK9EEJ46JD7YK4V8D9MCNM2BY5NBRS", "4b662a0e53db582574c612")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(temp)
}
