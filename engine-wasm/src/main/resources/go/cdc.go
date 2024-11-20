package main

import (
    "unsafe"
)

// alloc/free implementation from:
// https://github.com/tinygo-org/tinygo/blob/2a76ceb7dd5ea5a834ec470b724882564d9681b3/src/runtime/arch_tinygowasm_malloc.go#L7
var allocs = make(map[uintptr][]byte)

//export malloc
func libc_malloc(size uintptr) unsafe.Pointer {
	if size == 0 {
		return nil
	}
	buf := make([]byte, size)
	ptr := unsafe.Pointer(&buf[0])
	allocs[uintptr(ptr)] = buf
	return ptr
}

//export free
func libc_free(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}
	if _, ok := allocs[uintptr(ptr)]; ok {
		delete(allocs, uintptr(ptr))
	} else {
		panic("free: invalid pointer")
	}
}

func readCString(offset uint32) string {
    length := 0
    for {
        s := *(*int32)(unsafe.Pointer(uintptr(offset) + uintptr(length)))
        if (byte(s) == 0) {
            break
        }
        length++
    }

    buffer := make([]byte, length)
    for i := 0; i < int(length); i++ {
        s := *(*int32)(unsafe.Pointer(uintptr(offset) + uintptr(i)))
        buffer[i] = byte(s)
    }
    return string(buffer)
}

// inspired by:
// https://github.com/tinygo-org/tinygo/blob/2a76ceb7dd5ea5a834ec470b724882564d9681b3/src/runtime/string.go#L278
func writeCString(offset uintptr, str string) {
    stringData := []byte(str)
    for i := 0; i < len(stringData); i++ {
        *(*byte)(unsafe.Pointer(uintptr(offset) + uintptr(i))) = stringData[i]
    }
    *(*byte)(unsafe.Pointer(uintptr(offset) + uintptr(len(stringData)))) = 0 // trailing 0 byte
}

//go:wasm-module env
//export struct_get_string
func envStructGetString(structPtr, fieldNamePtr uint32) uint32

func structGetString(structPtr uint32, fieldName string) string {
    return "ciao"
}

//export process
func process(structPtr uint32) uint32 {
    envStructGetString(structPtr, 2)

    str := readCString(structPtr)

    var result = "giusto"
    if (str != "ciao") {
        result = "sbagliato"
    }

    var resultLen = len(result) + 1
    var resultPtr = libc_malloc(uintptr(resultLen))
    writeCString(uintptr(resultPtr), result)
    // return uint32(uintptr(resultPtr))
    return uint32(uintptr(resultPtr))
}

func main() {}
