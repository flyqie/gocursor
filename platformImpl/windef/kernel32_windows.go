//go:build windows && (386 || amd64)
// +build windows
// +build 386 amd64

package windef

import "syscall"

var (
	libKernel32, _ = syscall.LoadLibrary("kernel32.dll")
	// FuncGlobalAlloc https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globalalloc
	FuncGlobalAlloc, _ = syscall.GetProcAddress(syscall.Handle(libKernel32), "GlobalAlloc")
	// FuncGlobalFree https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globalfree
	FuncGlobalFree, _ = syscall.GetProcAddress(syscall.Handle(libKernel32), "GlobalFree")
)

const (
	// GMEMFIXED https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globalalloc
	GMEMFIXED = 0x0000
)
