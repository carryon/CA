package utils

import "reflect"

func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

func DelStringFromSlice(val string, slice *[]string) {
	s := make([]string, len(*slice)-1)
	for n, str := range *slice {
		if str == val {
			copy(s[:], (*slice)[:n])
			copy(s[n:], (*slice)[n+1:])
		}
	}
	*slice = s
}