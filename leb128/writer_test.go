package leb128

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAppend(t *testing.T) {
	type args struct {
		result []byte
		value  interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Test Append for Int",
			args: args{
				nil,
				123,
			},
			want: []byte{251, 0},
		},
		{
			name: "Test Append for Zero",
			args: args{
				nil,
				0,
			},
			want: []byte{0},
		},
		{
			name: "Test Append for String",
			args: args{
				nil,
				"aaaaaa",
			},
			want: []byte{6, 97, 97, 97, 97, 97, 97},
		},
		{
			name: "Test Append for float64",
			args: args{
				nil,
				float64(34.5),
			},
			want: []byte{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Append(tt.args.result, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	result := make([]byte, 0)
	value := float64(35.5)
	result = Append(result, value)
	fmt.Printf("%v", result)
}
