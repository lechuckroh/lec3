package image

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/gift"
)

// CalcRotatedSize returns image width/height after rotating given angle.
func CalcRotatedSize(w, h int, angle float32) (int, int) {
	if w <= 0 || h <= 0 {
		return 0, 0
	}

	xoff := float32(w)/2 - 0.5
	yoff := float32(h)/2 - 0.5

	asin, acos := Sincosf32(angle)
	x1, y1 := RotatePoint(0-xoff, 0-yoff, asin, acos)
	x2, y2 := RotatePoint(float32(w-1)-xoff, 0-yoff, asin, acos)
	x3, y3 := RotatePoint(float32(w-1)-xoff, float32(h-1)-yoff, asin, acos)
	x4, y4 := RotatePoint(0-xoff, float32(h-1)-yoff, asin, acos)

	minx := Minf32(x1, Minf32(x2, Minf32(x3, x4)))
	maxx := Maxf32(x1, Maxf32(x2, Maxf32(x3, x4)))
	miny := Minf32(y1, Minf32(y2, Minf32(y3, y4)))
	maxy := Maxf32(y1, Maxf32(y2, Maxf32(y3, y4)))

	neww := maxx - minx + 1
	if neww-Floorf32(neww) > 0.01 {
		neww += 2
	}
	newh := maxy - miny + 1
	if newh-Floorf32(newh) > 0.01 {
		newh += 2
	}
	return int(neww), int(newh)
}

// RotatePoint rotates a given point
func RotatePoint(x, y, asin, acos float32) (float32, float32) {
	newx := x*acos - y*asin
	newy := x*asin + y*acos
	return newx, newy
}

// RotateImage rotates the image by given angle.
// empty area after rotation is filled with bgColor
func RotateImage(
	src image.Image,
	angle float32,
	bgColor color.Color) image.Image {
	bounds := src.Bounds()
	width, height := CalcRotatedSize(bounds.Dx(), bounds.Dy(), angle)
	dest := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dest, dest.Bounds(),
		&image.Uniform{color.White},
		image.ZP,
		draw.Src)
	rotateFilter := gift.Rotate(angle, bgColor, gift.CubicInterpolation)
	gift.New(rotateFilter).Draw(dest, src)
	return dest
}
