//go:build linux
// +build linux

// Linux平台代码基于
//	https://github.com/lwch/rdesktop
//	https://github.com/linuxdeepin/go-x11-client
//	https://github.com/FyshOS/fynedesk
//
// 感谢原作者们创造性的工作!

package platformImpl

import (
	"encoding/binary"
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/xfixes"
	"image"
)

type Impl struct {
	display string
	xConn   *x.Conn
}

type ImplConfig struct {
	// Example :0
	Display string
	// 忽略panic
	IgnorePanic bool
}

type ImplCursorConfig struct {
}

func (i *Impl) New(config ImplConfig) error {
	if config.Display != "" {
		i.display = config.Display
	} else {
		i.display = ""
	}
	var err error
	if i.xConn, err = x.NewConnDisplay(i.display); err != nil {
		return err
	}
	_, _ = xfixes.QueryVersion(i.xConn, xfixes.MajorVersion, xfixes.MinorVersion).Reply(i.xConn)
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
	cookie := xfixes.GetCursorImage(i.xConn)
	rep, err := cookie.Reply(i.xConn)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, int(rep.Width), int(rep.Height)))
	curImgBytes := make([]byte, len(rep.CursorImage)*4)
	for j, v := range rep.CursorImage {
		binary.LittleEndian.PutUint32(curImgBytes[j*4:], v)
	}
	offset := 0
	for dy := 0; dy < int(rep.Height); dy++ {
		for dx := 0; dx < int(rep.Width); dx++ {
			img.Pix[offset+3] = curImgBytes[offset+3] // a
			img.Pix[offset+2] = curImgBytes[offset]   // b
			img.Pix[offset+1] = curImgBytes[offset+1] // g
			img.Pix[offset] = curImgBytes[offset+2]   // r
			offset += 4
		}
	}
	return img, nil
}

func (i *Impl) Free() error {
	i.xConn.Close()
	return nil
}

func (i *Impl) Status() interface{} {
	return map[string]string{
		"display": i.display,
	}
}
