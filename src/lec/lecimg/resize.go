package lecimg

import (
	"image"

	"github.com/disintegration/gift"
)

// ResizeImage resizes image to given dimension while preserving aspect ratio.
func ResizeImage(src image.Image, width, height int, keepAspectRatio bool) image.Image {
	var g *gift.GIFT
	if keepAspectRatio {
		g = gift.New(gift.ResizeToFit(width, height, gift.LanczosResampling))
	} else {
		g = gift.New(gift.Resize(width, height, gift.LanczosResampling))
	}
	dest := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dest, src)
	return dest
}
