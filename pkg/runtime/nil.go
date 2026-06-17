package runtime

import "reflect"

// IsNil checks if v is nil, it uses reflection to check so it will return true
// even for cases when v is a typed interfaces with nil value (which will not be
// equal to nil). For more info see: https://www.youtube.com/watch?v=ynoY2xz-F8s
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	value := reflect.ValueOf(v)
	kind := value.Kind()
	return kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil()
}
