package gocursor

import (
	"image/png"
	"os"
	"testing"
)

func TestCaptureCursor(t *testing.T) {
	i, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer i.Free()
	img, err := i.Capture()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Cursor Image W=%d H=%d", img.Bounds().Dx(), img.Bounds().Dy())
	f, err := os.Create("cursor.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCaptureCursorHelper(t *testing.T) {
	img, err := CaptureCursor()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Cursor Image W=%d H=%d", img.Bounds().Dx(), img.Bounds().Dy())
	f, err := os.Create("fast_cursor.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		t.Fatal(err)
	}
}
