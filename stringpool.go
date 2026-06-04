//go:build cgo

package imgui

// #include <stdlib.h>
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type StringPool struct {
	data     []byte
	spill    int
	cleanups []*C.char
	handles  []cgo.Handle
}

func (sp *StringPool) Reset() {
	for _, ptr := range sp.cleanups {
		C.free(unsafe.Pointer(ptr))
	}
	sp.cleanups = sp.cleanups[:0]

	for _, handle := range sp.handles {
		handle.Delete()
	}
	sp.handles = sp.handles[:0]

	if sp.spill > 0 {
		sp.data = make([]byte, 0, cap(sp.data)+sp.spill)
		sp.spill = 0
	} else {
		sp.data = sp.data[:0]
	}
}

func (sp *StringPool) StoreCString(input string) *C.char {
	length := len(input)

	if length+1 > cap(sp.data)-len(sp.data) {
		sp.spill += length + 1
		ptr := C.CString(input)
		sp.cleanups = append(sp.cleanups, ptr)
		return ptr
	}

	at := len(sp.data)
	sp.data = append(sp.data, input...)
	sp.data = append(sp.data, 0)
	s := sp.data[at : at+length+1]

	return (*C.char)(unsafe.Pointer(&s[0]))
}

func (sp *StringPool) StoreCStringArray(input []string) **C.char {
	if len(input) == 0 {
		return nil
	}

	ptr := C.malloc(C.size_t(len(input)) * C.size_t(unsafe.Sizeof((*C.char)(nil))))
	if ptr == nil {
		panic("imgui: failed to allocate C string array")
	}
	sp.cleanups = append(sp.cleanups, (*C.char)(ptr))

	items := unsafe.Slice((**C.char)(ptr), len(input))
	for i, item := range input {
		items[i] = C.CString(item)
		sp.cleanups = append(sp.cleanups, items[i])
	}

	return (**C.char)(ptr)
}

func (sp *StringPool) StoreHandle(value any) cgo.Handle {
	handle := cgo.NewHandle(value)
	sp.handles = append(sp.handles, handle)
	return handle
}
