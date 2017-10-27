package lecimg

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"lec/lecio"

	"golang.org/x/image/font/inconsolata"
)

type Changeable interface {
	Set(x, y int, c color.Color)
}

func SetColorAt(img image.Image, x, y int, color color.Color) {
	cImg := img.(Changeable)
	cImg.Set(x, y, color)
}

// LoadImage loads image from file.
func LoadImage(filename string) (image.Image, error) {
	var decoder func(io.Reader) (image.Image, error)

	ext := lecio.GetExt(filename)
	switch ext {
	case ".jpg", ".jpeg":
		decoder = jpeg.Decode
	case ".gif":
		decoder = gif.Decode
	case ".png":
		decoder = png.Decode
	}

	if decoder == nil {
		return nil, errors.New("Unsupported file format : " + ext)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		file.Close()
	}()

	img, err := decoder(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// SaveJpeg writes image as jpeg file.
func SaveJpeg(img image.Image, dir string, filename string, quality int) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
	}()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
}

func ToJpegBytes(img image.Image, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

// CreateImage creates an image with given size and background color.
func CreateImage(width, height int, bgColor color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, bgColor)
		}
	}
	return img
}

// FillRect draws a filled rectangle.
func FillRect(img *image.RGBA, x1, y1, x2, y2 int, rectColor color.Color) {
	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			img.Set(x, y, rectColor)
		}
	}
}

// DrawLine draw a line.
func DrawLine(img *image.RGBA, x1, y1, x2, y2 int, lineColor color.Color) {
	dx, dy := x2-x1, y2-y1
	if dx <= dy {
		incX := float32(dx) / float32(dy)
		x := float32(x1)
		for y := y1; y < y2; y++ {
			img.Set(int(x), y, lineColor)
			x += incX
		}
	} else {
		incY := float32(dy) / float32(dx)
		y := float32(y1)
		for x := x1; x < x2; x++ {
			img.Set(x, int(y), lineColor)
			y += incY
		}
	}
}

func DrawLabelBold8x16(img *image.RGBA, x, y int, label string, fgColor color.Color) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fgColor),
		Face: inconsolata.Bold8x16,
		Dot:  point,
	}
	d.DrawString(label)
}
