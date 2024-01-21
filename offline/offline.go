package main

// #include <stdlib.h>
import "C"

//export Run
func Run(*C.char) C.int {
	println("Please enable network during cargo build")
	return 0
}

func main() {}
