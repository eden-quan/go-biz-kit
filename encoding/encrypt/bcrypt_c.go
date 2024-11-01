package encrypt

//
///*
// #cgo CFLAGS: -I/Users/sinsay/project/work/lainuo/project/eden/go-biz-kit/encoding/encrypt
// #include "bcrypt_all.h"
//*/
//import "C"
//import (
//	"unsafe"
//)
//
//// BcryptGenerateHash ...
//func BcryptGenerateHash(psw string) string {
//	cInput := C.CString(psw)
//	defer C.free(unsafe.Pointer(cInput))
//	output := make([]byte, 1024)
//	cOutput := (*C.char)(unsafe.Pointer(&output[0]))
//	C.GenerateHash(cInput, C.uint(len(psw)), C.uint(4), cOutput)
//	return C.GoString(cOutput)
//}
//
//// BcryptValidatePassword ...
//func BcryptValidatePassword(psw string, hashed string) bool {
//	cPsw := C.CString(psw)
//	cHashed := C.CString(hashed)
//	cInt := C.ValidatePassword(cPsw, C.uint(len(psw)), cHashed)
//	return cInt == 0
//}
