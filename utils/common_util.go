package utils

func NewBool(a bool) *bool {
	return &a
}

func NewInt(a int) *int {
	return &a
}

func NewInt8(a int8) *int8 {
	return &a
}

func NewInt16(a int16) *int16 {
	return &a
}

func NewInt32(a int32) *int32 {
	return &a
}

func NewInt64(a int64) *int64 {
	return &a
}

func NewFloat32(a float32) *float32 {
	return &a
}

func NewFloat64(a float64) *float64 {
	return &a
}

func NewString(a string) *string {
	return &a
}

func MinInt32(a, b int32) int32 {
	return IfElse(a < b, a, b).(int32)
}

func MaxInt32(a, b int32) int32 {
	return IfElse(a > b, a, b).(int32)
}

func MinInt(a, b int) int {
	return IfElse(a < b, a, b).(int)
}

func MaxInt(a, b int) int {
	return IfElse(a > b, a, b).(int)
}

func IfElse(ok bool, onTrue, onFalse interface{}) interface{} {
	if ok {
		return onTrue
	} else {
		return onFalse
	}
}
