package main

// #include <stdlib.h>
import "C"

import (
	"github.com/webteleport/ufo/apps/teleport"
	"log"
	"unsafe"
)

//go:generate env CGO_ENABLED=1 go build -v -o libteleport.so -buildmode=c-shared .

//export Run2
func Run2(cstr1, cstr2 *C.char) {
	str1 := C.GoString(cstr1)
	str2 := C.GoString(cstr2)
	println(str1)
	println(str2)
	teleport.Run([]string{str1, str2})
	println("2")
}

func isNULL(p *C.char) bool {
	if uint(uintptr(unsafe.Pointer(p))) == 0x00 {
		return true
	}
	return byte((uint8)(*p)) == 0x00
}

func offsetof(n int, base uintptr) uintptr {
	return base * uintptr(n)
}

//export Run
func Run(cstrs *C.char) {
	strs := []string{}
	for !isNULL(cstrs) {
		str := C.GoString(cstrs)
		strs = append(strs, str)
		// println(str)
		cstrs = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cstrs)) + offsetof((len(str)+1), unsafe.Sizeof(*cstrs))))
	}
	if err := teleport.Run(strs); err != nil {
		log.Println(err)
	}
}

func main() {
}
