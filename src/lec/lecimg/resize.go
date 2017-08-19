package lecimg

import (
	"image"

	"github.com/disintegration/gift"
)

// ResizeImage resizes image to given dimension while preserving aspect ratio.
func ResizeImage(src image.Image, width, height int) image.Image {
	g := gift.New(gift.ResizeToFit(width, height, gift.LanczosResampling))
	dest := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dest, src)
	return dest
}
