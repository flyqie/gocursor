//go:build windows && (386 || amd64)
// +build windows
// +build 386 amd64

// Windows平台代码源自
//	https://github.com/lwch/rdesktop
//	https://github.com/TigerVNC/tigervnc
//
// 感谢原作者们创造性的工作!

package platformImpl

import (
	"fmt"
	"github.com/flyqie/gocursor/platformImpl/windef"
	"image"
	"syscall"
	"unsafe"
)

type Impl struct {
	hwnd   uintptr
	hdc    uintptr
	bits   uintptr
	width  uintptr
	height uintptr
	buffer uintptr
}

type ImplConfig struct {
}

type ImplCursorConfig struct {
	// 忽略panic
	IgnorePanic bool
}

func (i *Impl) New(config ImplConfig) error {
	err := i.getHandle()
	if err != nil {
		return err
	}
	bits, _, err := syscall.Syscall(windef.FuncGetDeviceCaps, 2, i.hdc, windef.BITSPIXEL, 0)
	if bits == 0 {
		return fmt.Errorf("get device caps(bits): %v", err)
	}
	if bits != 32 {
		bits = 32
	}
	width, _, err := syscall.Syscall(windef.FuncGetDeviceCaps, 2, i.hdc, windef.HORZRES, 0)
	if width == 0 {
		return fmt.Errorf("get device caps(width): %v", err)
	}
	height, _, err := syscall.Syscall(windef.FuncGetDeviceCaps, 2, i.hdc, windef.VERTRES, 0)
	if height == 0 {
		return fmt.Errorf("get device caps(height): %v", err)
	}
	if i.width != width || i.height != height {
		err := i.resizeBuffer(int(bits), int(width), int(height))
		if err != nil {
			return err
		}
	}
	i.bits = bits
	i.width = width
	i.height = height
	return nil
}

func (i *Impl) getHandle() error {
	hwnd, _, err := syscall.Syscall(windef.FuncGetDesktopWindow, 0, 0, 0, 0)
	if hwnd == 0 {
		return fmt.Errorf("get desktop window: %v", err)
	}
	hdc, _, err := syscall.Syscall(windef.FuncGetDC, 1, i.hwnd, 0, 0)
	if hdc == 0 {
		return fmt.Errorf("get dc: %v", err)
	}
	if i.hdc != 0 {
		syscall.Syscall(windef.FuncReleaseDC, 2, i.hwnd, i.hdc, 0)
	}
	i.hwnd = hwnd
	i.hdc = hdc
	return nil
}

func (i *Impl) resizeBuffer(bits, width, height int) error {
	addr, _, err := syscall.Syscall(windef.FuncGlobalAlloc, 2, windef.GMEMFIXED, uintptr(bits*width*height/8), 0)
	if addr == 0 {
		return fmt.Errorf("global alloc: %v", err)
	}
	if i.buffer != 0 {
		syscall.Syscall(windef.FuncGlobalFree, 1, i.buffer, 0, 0)
	}
	i.buffer = addr
	return nil
}

func (i *Impl) CaptureWithConfig(config ImplCursorConfig) (img *image.RGBA, err error) {
	if config.IgnorePanic {
		defer func(_img *image.RGBA) {
			if _err := recover(); _err != nil {
				_img = nil
			}
		}(img)
	}
	var curInfo windef.CURSORINFO
	curInfo.CbSize = windef.DWORD(unsafe.Sizeof(curInfo))
	ok, _, err := syscall.SyscallN(windef.FuncGetCursorInfo, uintptr(unsafe.Pointer(&curInfo)), 0, 0)
	if ok == 0 {
		return nil, fmt.Errorf("get cursor info: %v", err)
	}
	var info windef.ICONINFO
	ok, _, err = syscall.SyscallN(windef.FuncGetIconInfo, uintptr(curInfo.HCursor), uintptr(unsafe.Pointer(&info)), 0)
	if ok == 0 {
		return nil, fmt.Errorf("get icon info: %v", err)
	}

	var bitmap windef.BITMAP
	ok, _, err = syscall.SyscallN(windef.FuncGetObject, uintptr(info.HbmMask), unsafe.Sizeof(bitmap), uintptr(unsafe.Pointer(&bitmap)))
	if ok == 0 {
		return nil, fmt.Errorf("get object: %v", err)
	}

	defer syscall.SyscallN(windef.FuncDeleteObject, uintptr(info.HbmColor))
	defer syscall.SyscallN(windef.FuncDeleteObject, uintptr(info.HbmMask))

	var copyFromMask bool
	if info.HbmColor == 0 {
		copyFromMask = true
		bitmap.BmHeight /= 2
	}

	hdcMem, _, err := syscall.SyscallN(windef.FuncCreateCompatibleDC, i.hdc)
	if hdcMem == 0 {
		return nil, fmt.Errorf("create compatible dc: %v", err)
	}
	defer syscall.SyscallN(windef.FuncDeleteDC, hdcMem)
	canvas, _, err := syscall.SyscallN(windef.FuncCreateCompatibleBitmap, i.hdc,
		uintptr(bitmap.BmWidth), uintptr(bitmap.BmHeight))
	if canvas == 0 {
		return nil, fmt.Errorf("create compatible bitmap: %v", err)
	}
	defer syscall.SyscallN(windef.FuncDeleteObject, canvas)

	old, _, err := syscall.SyscallN(windef.FuncSelectObject, hdcMem, canvas)
	if old == 0 {
		return nil, fmt.Errorf("select object: %v", err)
	}
	defer syscall.SyscallN(windef.FuncSelectObject, hdcMem, old)

	img = image.NewRGBA(image.Rect(0, 0, int(bitmap.BmWidth), int(bitmap.BmHeight)))
	defer i.copyImageData(i.hdc, canvas, img, copyFromMask, bitmap, info.HbmMask, config.IgnorePanic)

	ok, _, err = syscall.SyscallN(windef.FuncDrawIcon, hdcMem, 0, 0, uintptr(curInfo.HCursor))
	if ok == 0 {
		return nil, fmt.Errorf("draw icon: %v", err)
	}

	return img, nil
}

// BITMAPINFOHEADER https://docs.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-bitmapinfoheader
type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

func (i *Impl) copyImageData(hdc, bitmap uintptr, img *image.RGBA, copyFromMask bool, bitmapInfo windef.BITMAP, hBitmap windef.HBITMAP, ignorePanic bool) {
	if ignorePanic {
		defer func(_img *image.RGBA) {
			if _err := recover(); _err != nil {
				_img = nil
			}
		}(img)
	}
	if copyFromMask {
		_, _, _ = syscall.SyscallN(windef.FuncGetBitmapBits, uintptr(hBitmap), uintptr(bitmapInfo.BmWidthBytes*(bitmapInfo.BmHeight*2)), i.buffer)
		offset := 0
		for dy := 0; dy < img.Rect.Max.Y; dy++ {
			for dx := 0; dx < img.Rect.Max.X; dx++ {
				_byte := dy*int(bitmapInfo.BmWidthBytes) + dx/8
				_bit := 7 - dx%8

				if (*(*uint8)(unsafe.Pointer(i.buffer + uintptr(_byte))) & (1 << _bit)) == 0 {
					// Valid pixel, so make it opaque
					img.Pix[offset+3] = 0xff

					// Black or white?
					if (*(*uint8)(unsafe.Pointer(i.buffer + uintptr(bitmapInfo.BmWidthBytes*bitmapInfo.BmHeight) + uintptr(_byte))) & (1 << _bit)) == 0 {
						img.Pix[offset+0] = 0xff
						img.Pix[offset+1] = 0xff
						img.Pix[offset+2] = 0xff
					} else {
						img.Pix[offset+0] = 0
						img.Pix[offset+1] = 0
						img.Pix[offset+2] = 0
					}
				} else if (*(*uint8)(unsafe.Pointer(i.buffer + uintptr(bitmapInfo.BmWidthBytes*bitmapInfo.BmHeight) + uintptr(_byte))) & (1 << _bit)) != 0 {
					img.Pix[offset+0] = 0
					img.Pix[offset+1] = 0
					img.Pix[offset+2] = 0
					img.Pix[offset+3] = 0xff
				} else {
					// Transparent pixel
					img.Pix[offset+0] = 0
					img.Pix[offset+1] = 0
					img.Pix[offset+2] = 0
					img.Pix[offset+3] = 0
				}
				offset += 4
			}
		}
		// 重设光标到白色, 上面代码拿到的是黑色
		_offset := 0
		for dy := 0; dy < img.Rect.Max.Y; dy++ {
			for dx := 0; dx < img.Rect.Max.X; dx++ {
				if img.Pix[_offset+3] == 0xff && img.Pix[_offset+2] == 0 && img.Pix[_offset+1] == 0 && img.Pix[_offset] == 0 {
					img.Pix[_offset+2] = 0xff
					img.Pix[_offset+1] = 0xff
					img.Pix[_offset] = 0xff
				}
				_offset += 4
			}
		}
	} else {
		var hdr BITMAPINFOHEADER
		hdr.BiSize = uint32(unsafe.Sizeof(hdr))
		hdr.BiPlanes = 1
		hdr.BiBitCount = uint16(i.bits)
		hdr.BiWidth = int32(img.Rect.Max.X)
		hdr.BiHeight = int32(-img.Rect.Max.Y)
		hdr.BiCompression = windef.BIRGB
		hdr.BiSizeImage = 0
		_, _, _ = syscall.SyscallN(windef.FuncGetDIBits, hdc, bitmap, 0, uintptr(img.Rect.Max.Y),
			i.buffer, uintptr(unsafe.Pointer(&hdr)), windef.DIBRGBCOLORS)
		// TODO: support difference of 32 bits
		for j := 0; j < len(img.Pix); j++ {
			img.Pix[j] = *(*uint8)(unsafe.Pointer(i.buffer + uintptr(j)))
		}
		// BGR => RGB
		for j := 0; j < len(img.Pix); j += int(i.bits / 8) {
			img.Pix[j], img.Pix[j+2] = img.Pix[j+2], img.Pix[j]
		}
	}
}

func (i *Impl) Free() error {
	if i.hwnd != 0 && i.hdc != 0 {
		syscall.Syscall(windef.FuncReleaseDC, 2, i.hwnd, i.hdc, 0)
	}
	return nil
}

func (i *Impl) Status() interface{} {
	return nil
}
