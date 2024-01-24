package main

// #include <stdlib.h>
import "C"

import (
	"log"
	"unsafe"

	"github.com/btwiuse/pub"
)

//export Run2
func Run2(cstr1, cstr2 *C.char) {
	str1 := C.GoString(cstr1)
	str2 := C.GoString(cstr2)
	println(str1)
	println(str2)
	pub.Run([]string{str1, str2})
	println("2")
}

// Helper function to check if C-style string is NULL
func isNULL(p *C.char) bool {
	return uint(uintptr(unsafe.Pointer(p))) == 0x00 || byte((uint8)(*p)) == 0x00
}

// Helper function to advance C-style string array pointer
func advancePointer(cArray *C.char, n int) *C.char {
	return (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cArray)) + unsafe.Sizeof(*cArray)*(uintptr(n))))
}

// Helper function to convert C-style string array to Go string slice
func cArrayToGoSlice(cArray *C.char) []string {
	strs := []string{}
	for !isNULL(cArray) {
		str := C.GoString(cArray)
		cArray = advancePointer(cArray, len(str)+1)
		strs = append(strs, str)
	}
	return strs
}

//export Run
func Run(cstrs *C.char) C.int {
	strs := cArrayToGoSlice(cstrs)
	if err := pub.Run(strs); err != nil {
		log.Println(err)
		return -1
	}
	return 0
}

func main() {}
