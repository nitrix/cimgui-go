//go:build cgo

package imgui

/*
#define CIMGUI_DEFINE_ENUMS_AND_STRUCTS 1
#include "dist/cimgui/cimgui.h"

extern int imguiInputTextCallbackBridge(ImGuiInputTextCallbackData* data);
extern void imguiSizeCallbackBridge(ImGuiSizeCallbackData* data);

static inline ImGuiInputTextCallback imguiGoInputTextCallback(void) {
	return imguiInputTextCallbackBridge;
}

static inline ImGuiSizeCallback imguiGoSizeCallback(void) {
	return imguiSizeCallbackBridge;
}
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
)

type InputTextCallbackFunc func(data *InputTextCallbackData) int

type SizeCallbackFunc func(data *SizeCallbackData)

func handleUserData(handle cgo.Handle) unsafe.Pointer {
	return unsafe.Pointer(uintptr(handle))
}

func InputTextWithCallback(label string, buf []byte, flags InputTextFlags, callback InputTextCallbackFunc) bool {
	a0 := stringPool.StoreCString(label)
	a1 := (*C.char)(nil)
	if len(buf) > 0 {
		a1 = (*C.char)(unsafe.Pointer(&buf[0]))
	}
	a2 := (C.size_t)(len(buf))
	a3 := (C.ImGuiInputTextFlags)(flags)
	a4 := C.ImGuiInputTextCallback(nil)
	a5 := unsafe.Pointer(nil)
	if callback != nil {
		handle := stringPool.StoreHandle(callback)
		a4 = C.imguiGoInputTextCallback()
		a5 = handleUserData(handle)
	}

	call := C.igInputText(a0, a1, a2, a3, a4, a5)
	return bool(call)
}

func InputTextMultilineWithCallback(label string, buf []byte, size mgl32.Vec2, flags InputTextFlags, callback InputTextCallbackFunc) bool {
	a0 := stringPool.StoreCString(label)
	a1 := (*C.char)(nil)
	if len(buf) > 0 {
		a1 = (*C.char)(unsafe.Pointer(&buf[0]))
	}
	a2 := (C.size_t)(len(buf))
	a3 := mglVec2ToImVec2(size)
	a4 := (C.ImGuiInputTextFlags)(flags)
	a5 := C.ImGuiInputTextCallback(nil)
	a6 := unsafe.Pointer(nil)
	if callback != nil {
		handle := stringPool.StoreHandle(callback)
		a5 = C.imguiGoInputTextCallback()
		a6 = handleUserData(handle)
	}

	call := C.igInputTextMultiline(a0, a1, a2, a3, a4, a5, a6)
	return bool(call)
}

func InputTextWithHintAndCallback(label string, hint string, buf []byte, flags InputTextFlags, callback InputTextCallbackFunc) bool {
	a0 := stringPool.StoreCString(label)
	a1 := stringPool.StoreCString(hint)
	a2 := (*C.char)(nil)
	if len(buf) > 0 {
		a2 = (*C.char)(unsafe.Pointer(&buf[0]))
	}
	a3 := (C.size_t)(len(buf))
	a4 := (C.ImGuiInputTextFlags)(flags)
	a5 := C.ImGuiInputTextCallback(nil)
	a6 := unsafe.Pointer(nil)
	if callback != nil {
		handle := stringPool.StoreHandle(callback)
		a5 = C.imguiGoInputTextCallback()
		a6 = handleUserData(handle)
	}

	call := C.igInputTextWithHint(a0, a1, a2, a3, a4, a5, a6)
	return bool(call)
}

func SetNextWindowSizeConstraintsWithCallback(sizeMin mgl32.Vec2, sizeMax mgl32.Vec2, callback SizeCallbackFunc) {
	a0 := mglVec2ToImVec2(sizeMin)
	a1 := mglVec2ToImVec2(sizeMax)
	a2 := C.ImGuiSizeCallback(nil)
	a3 := unsafe.Pointer(nil)
	if callback != nil {
		handle := stringPool.StoreHandle(callback)
		a2 = C.imguiGoSizeCallback()
		a3 = handleUserData(handle)
	}

	C.igSetNextWindowSizeConstraints(a0, a1, a2, a3)
}
