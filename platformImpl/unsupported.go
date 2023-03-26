//go:build !(windows && (386 || amd64)) && !linux && !darwin
// +build !windows !386,!amd64
// +build !linux
// +build !darwin

package platformImpl

import (
	"errors"
	"image"
)

var (
	UnSupportedError = errors.New("unsupported impl")
)

type Impl struct {
}

type ImplConfig struct {
}

type ImplCursorConfig struct {
}

func (i *Impl) New(config ImplConfig) error {
	return UnSupportedError
}

func (i *Impl) CaptureWithConfig(config ImplCursorConfig) (img *image.RGBA, err error) {
	return nil, UnSupportedError
}

func (i *Impl) Free() error {
	return nil
}

func (i *Impl) Status() interface{} {
	return nil
}
