package leb128

import (
	"fmt"
	"reflect"
	"testing"
)

type A struct {
	int
}

func TestDebug(t *testing.T) {
	// value := float64(35.5)
	value := A{543}
	result, err := Compress(value)
	fmt.Printf("%v, %v", result, err)
}

func TestCompress(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Try compress string",
			args: args{
				"aaaaaa",
			},
			want: []byte{6,95,95,95,95,95,95},
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compress() = %v, want %v", got, tt.want)
			}
		})
	}
}
