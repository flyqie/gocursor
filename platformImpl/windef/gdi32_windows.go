//go:build windows && (386 || amd64)
// +build windows
// +build 386 amd64

package windef

import "syscall"

var (
	libGdi32, _ = syscall.LoadLibrary("Gdi32.dll")
	// FuncGetDeviceCaps https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getdevicecaps
	FuncGetDeviceCaps, _ = syscall.GetProcAddress(syscall.Handle(libGdi32), "GetDeviceCaps")
	// FuncCreateCompatibleDC https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-createcompatibledc
	FuncCreateCompatibleDC, _ = syscall.GetProcAddress(syscall.Handle(libGdi32), "CreateCompatibleDC")
	// FuncCreateCompatibleBitmap https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-createcompatiblebitmap
	FuncCreateCompatibleBitmap, _ = syscall.GetProcAddress(syscall.Handle(libGdi32), "CreateCompatibleBitmap")
	// FuncDeleteDC https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-deletedc
	FuncDeleteDC, _ = syscall.GetProcAddress(syscall.Handle(libGdi32), "DeleteDC")
	// FuncSelectObject https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-selectobject
	FuncSelectObject, _ = syscall.GetProcAddress(syscall.Handle(libGdi32), "SelectObject")
	// FuncDeleteObject https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-deleteobject
	FuncDeleteObject, _ = syscall.GetProcAddress(syscall.Handle(libGdi32), "DeleteObject")
	// FuncGetDIBits https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getdibits
	FuncGetDIBits, _ = syscall.GetProcAddress(syscall.Handle(libGdi32), "GetDIBits")
	// FuncGetObject https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getobject
	FuncGetObject, _ = syscall.GetProcAddress(syscall.Handle(libGdi32), "GetObjectW")
	// FuncGetBitmapBits https://learn.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getbitmapbits
	FuncGetBitmapBits, _ = syscall.GetProcAddress(syscall.Handle(libGdi32), "GetBitmapBits")
)

const (
	// HORZRES https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getdevicecaps
	HORZRES = 8
	// VERTRES https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getdevicecaps
	VERTRES = 10
	// BITSPIXEL https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getdevicecaps
	BITSPIXEL = 12
)

const (
	// BIRGB https://docs.microsoft.com/en-us/previous-versions/dd183376(v=vs.85)
	BIRGB = 0
	// DIBRGBCOLORS https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getdibits
	DIBRGBCOLORS = 0
)

type BITMAP struct {
	BmType       LONG
	BmWidth      LONG
	BmHeight     LONG
	BmWidthBytes LONG
	BmPlanes     WORD
	BmBitsPixel  WORD
	BmPixel      LPVOID
}
