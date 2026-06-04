//go:build cgo

package imgui

// #define CIMGUI_DEFINE_ENUMS_AND_STRUCTS 1
// #include "dist/cimgui/cimgui.h"
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

//export imguiInputTextCallbackBridge
func imguiInputTextCallbackBridge(data *C.ImGuiInputTextCallbackData) C.int {
	if data == nil || data.UserData == nil {
		return 0
	}

	handle := cgo.Handle(uintptr(data.UserData))
	callback, ok := handle.Value().(InputTextCallbackFunc)
	if !ok || callback == nil {
		return 0
	}

	return C.int(callback((*InputTextCallbackData)(unsafe.Pointer(data))))
}

//export imguiSizeCallbackBridge
func imguiSizeCallbackBridge(data *C.ImGuiSizeCallbackData) {
	if data == nil || data.UserData == nil {
		return
	}

	handle := cgo.Handle(uintptr(data.UserData))
	callback, ok := handle.Value().(SizeCallbackFunc)
	if ok && callback != nil {
		callback((*SizeCallbackData)(unsafe.Pointer(data)))
	}
}
