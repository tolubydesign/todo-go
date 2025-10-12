package helper

func ConvertIntToPointerInt(value int) *int {
	return &value
}

// func ConvertPointerIntToInt(pointer *int) int {
// 	unsafePtr := unsafe.Pointer(pointer)
// uintPtrVal := uintptr(unsafePtr)
// 	return uintPtrVal
// }
