package gocursor

import (
	"errors"
	"github.com/flyqie/gocursor/platformImpl"
	"github.com/hashicorp/go-multierror"
	"image"
)

// Errors
var (
	NewPlatformImplError = errors.New("platform impl new() error")
	InstanceUsableError  = errors.New("instance not usable")
	CaptureError         = errors.New("cursor capture error")
	FreeError            = errors.New("instance has been freed")
)

type Cursor struct {
	usable   bool
	platImpl *platformImpl.Impl
}

// NewWithConfig 使用自定义配置创建实例
func NewWithConfig(config platformImpl.ImplConfig) (*Cursor, error) {
	cur := &Cursor{}
	cur.platImpl = &platformImpl.Impl{}
	if err := cur.platImpl.New(config); err != nil {
		return nil, multierror.Append(nil, NewPlatformImplError, err)
	}
	cur.usable = true
	return cur, nil
}

// New 使用默认配置创建实例
func New() (*Cursor, error) {
	return NewWithConfig(platformImpl.ImplConfig{})
}

// CaptureCursor 快速使用默认配置捕获鼠标指针图像
func CaptureCursor() (img *image.RGBA, err error) {
	cur, err := New()
	if err != nil {
		return
	}
	img, err = cur.Capture()
	if err != nil {
		return
	}
	err = cur.Free()
	return
}

// CaptureWithConfig 使用自定义配置捕捉鼠标指针图像
func (c *Cursor) CaptureWithConfig(config platformImpl.ImplCursorConfig) (*image.RGBA, error) {
	if !c.usable {
		return nil, multierror.Append(nil, InstanceUsableError)
	}
	curImg, err := c.platImpl.CaptureWithConfig(config)
	if err != nil {
		return nil, multierror.Append(nil, CaptureError, err)
	}
	return curImg, nil
}

// Capture 使用默认配置捕捉鼠标指针图像
func (c *Cursor) Capture() (*image.RGBA, error) {
	return c.CaptureWithConfig(platformImpl.ImplCursorConfig{})
}

// Free 释放实例
func (c *Cursor) Free() error {
	if !c.usable {
		// 实例已被释放
		return nil
	}
	err := c.platImpl.Free()
	if err != nil {
		return multierror.Append(nil, FreeError, err)
	}
	c.usable = false
	c.platImpl = nil
	return nil
}
