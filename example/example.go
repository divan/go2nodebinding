package main

import "C"

//export MyFunc
func MyFunc(arg1, arg2 *C.char, arg3 C.int) *C.char {
	// do nothing
}
