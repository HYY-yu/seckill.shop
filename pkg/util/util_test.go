package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsZero(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"int zero",
			args{
				data: 0,
			},
			true,
		},
		{
			"float zero",
			args{
				data: 0.0,
			},
			true,
		},
		{
			"string zero",
			args{
				data: "",
			},
			true,
		},
		{
			"string not zero",
			args{
				data: "0",
			},
			false,
		},
		{
			"array zero",
			args{
				data: [3]int{},
			},
			true,
		},
		{
			"map zero",
			args{
				data: make(map[string]interface{}),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsZero(tt.args.data); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsZeroForSlice(t *testing.T) {
	var s []int
	st := IsZero(s)
	assert.Equal(t, true, st)

	var s2 = make([]int, 0)
	st2 := IsZero(s2)
	assert.Equal(t, false, st2)
}
