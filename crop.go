package main

import (
	"image"
	"image/color"
	"image/draw"
)

// GetCropRect gets constraint satisfied cropped Rectagle
func GetCropRect(left, top, right, bottom int, bounds image.Rectangle, maxWidthCropRate, maxHeightCropRate, minRatio, maxRatio float32) image.Rectangle {
	initWidth, initHeight := right-left, bottom-top
	width, height := initWidth, initHeight
	imgWidth, imgHeight := bounds.Dx(), bounds.Dy()

	// maxCropRate
	minWidth := int(float32(imgWidth) * (1 - maxWidthCropRate))
	minHeight := int(float32(imgHeight) * (1 - maxHeightCropRate))
	width, height = Max(width, minWidth), Max(height, minHeight)

	// ratio
	ratio := float32(height) / float32(width)
	if ratio < minRatio {
		height = Max(minHeight, int(float32(width)*minRatio))
	}
	if ratio > maxRatio {
		width = Max(minWidth, int(float32(height)/maxRatio))
	}

	// adjust border
	widthInc, heightInc := width-initWidth, height-initHeight
	widthMargin, heightMargin := width-initWidth, height-initHeight

	if widthInc > 0 {
		widthHalfMargin := int(float32(widthMargin) / 2)
		leftMargin := Min(left, widthHalfMargin)
		rightMargin := Min(imgWidth-right, widthMargin-leftMargin)
		left -= leftMargin
		right += rightMargin

		w := right - left
		dx := widthInc - w + initWidth
		if dx > 0 {
			widthSpace := left + imgWidth - right
			if widthSpace == 0 {
				right += dx
			} else {
				leftRatio := left / widthSpace
				leftSpace := dx * leftRatio
				left -= leftSpace
				right += dx - leftSpace
			}
		}
	}

	if heightInc > 0 {
		heightHalfMargin := int(float32(heightMargin) / 2)
		topMargin := Min(top, heightHalfMargin)
		bottomMargin := Min(imgHeight-bottom, heightMargin-topMargin)
		top -= topMargin
		bottom += bottomMargin

		h := bottom - top
		dy := heightInc - h + initHeight
		if dy > 0 {
			heightSpace := top + imgHeight - bottom
			if heightSpace == 0 {
				bottom += dy
			} else {
				topRatio := top / heightSpace
				topSpace := dy * topRatio
				top -= topSpace
				bottom += dy - topSpace
			}
		}
	}

	return image.Rect(left, top, right, bottom)
}

func CropImage(src image.Image, rect image.Rectangle) image.Image {
	slicedImage := src.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)

	result := CreateImage(rect.Dx(), rect.Dy(), color.White)
	draw.Draw(result,
		image.Rect(-rect.Min.X, -rect.Min.Y, rect.Max.X, rect.Max.Y),
		slicedImage,
		image.ZP,
		draw.Src)

	return result
}
