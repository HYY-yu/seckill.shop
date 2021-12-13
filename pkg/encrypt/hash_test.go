package encrypt

import "testing"

func TestSalt(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			want: "x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Salt(NewHexEncoding()); got != tt.want {
				t.Errorf("Salt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNonce(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test2",
			args: args{
				n: 6,
			},
			want: "x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Nonce(tt.args.n); got != tt.want {
				t.Errorf("Nonce() = %v, want %v", got, tt.want)
			}
		})
	}
}
