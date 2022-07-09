package until

import (
	"fmt"
	"reflect"
)

func ToFloat64(number interface{}) (float64, error) {
	v := reflect.ValueOf(number)
	if v.CanFloat() {
		return v.Float(), nil
	}
	return 0, fmt.Errorf("can`t conver %v to Float64", reflect.TypeOf(number))
}

func ToUInt64(number interface{}) (uint64, error) {
	v := reflect.ValueOf(number)
	if v.CanUint() {
		return v.Uint(), nil
	}
	return 0, fmt.Errorf("can`t conver %v to UInt64", reflect.TypeOf(number))
}

func ToInt64(number interface{}) (int64, error) {
	v := reflect.ValueOf(number)
	if v.CanInt() {
		return v.Int(), nil
	}
	return 0, fmt.Errorf("can`t conver %v to Int64", reflect.TypeOf(number))
}

func ToString(str interface{}) (string, error) {
	s, ok := str.(string)
	if ok {
		return s, nil
	}
	return "", fmt.Errorf("can`t conver %v to String", reflect.TypeOf(str))

}
