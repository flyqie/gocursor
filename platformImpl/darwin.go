//go:build darwin
// +build darwin

// Darwin平台代码源自
//
//	https://github.com/lwch/rdesktop
//
// 感谢原作者创造性的工作!
package platformImpl

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework AppKit
#include <CoreGraphics/CoreGraphics.h>

void get_cursor_size(int *width, int *height);
void cursor_copy(unsigned char* pixels, int width, int height);
*/
import "C"

import (
	"errors"
	"fmt"
	"image"
	"unsafe"
)

var (
	UnSupportedError = errors.New("unsupported impl")
)

type Impl struct {
	id C.CGDirectDisplayID
}

type ImplConfig struct {
}

type ImplCursorConfig struct {
	// 忽略panic
	IgnorePanic bool
}

func (i *Impl) New(config ImplConfig) error {
	i.id = getDisplayID()
	return nil
}

func (i *Impl) getDisplayID() C.CGDirectDisplayID {
	var id C.CGDirectDisplayID
	if C.CGGetActiveDisplayList(C.uint32_t(1), (*C.CGDirectDisplayID)(unsafe.Pointer(&id)), nil) != C.kCGErrorSuccess {
		return 0
	}
	return id
}

func (i *Impl) CaptureWithConfig(config ImplCursorConfig) (img *image.RGBA, err error) {
	if config.IgnorePanic {
		defer func(_img *image.RGBA) {
			if _err := recover(); _err != nil {
				_img = nil
			}
		}(img)
	}
	var width, height C.int
	C.get_cursor_size(&width, &height)
	if width == 0 || height == 0 {
		return nil, fmt.Errorf("can not get cursor size")
	}
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	C.cursor_copy((*C.uchar)(unsafe.Pointer(&img.Pix[0])), width, height)
	return img, nil
}

func (i *Impl) Free() error {
	C.CGDisplayRelease(i.id)
	return nil
}

func (i *Impl) Status() interface{} {
	return nil
}
