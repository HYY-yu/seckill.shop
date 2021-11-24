package encrypt

import "testing"

func TestEncrypt(t *testing.T) {
	t.Log(New().MD5("123456"))
}

func BenchmarkEncrypt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New().MD5("123456")
	}
}
