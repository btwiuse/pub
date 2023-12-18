package main

// #include <stdlib.h>
import "C"

import (
	"github.com/webteleport/ufo/apps/teleport"
	"unsafe"
	"log"
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

func convertToByte(c C.char) byte {
    return byte((uint8)(c))
}

func offsetof(n int, base uintptr) uintptr {
	return base * uintptr(n)
}

//export Run
func Run(cstrs *C.char) {
	strs := []string{}
	for {
		leading := *cstrs
		if convertToByte(leading) == 0x00 {
			break
		}
		str := C.GoString(cstrs)
		strs = append(strs, str)
		// println(str)
		cstrs = (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cstrs)) + offsetof((len(str) + 1), unsafe.Sizeof(*cstrs))))
	}
	if err := teleport.Run(strs); err != nil {
		log.Println(err)
	}
}

func main() {
}
