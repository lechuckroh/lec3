package main

import (
	"image"
	"image/color"
	"testing"
)

func testChangeLineSpace(t *testing.T, img image.Image, option ChangeLineSpaceOption, expectedHeight int) {
	// Run Filter
	result := NewChangeLineSpaceFilter(option).Run(NewFilterSource(img, "filename", 0))

	// Test result image size
	destBounds := result.Img().Bounds()
	heightMatch := destBounds.Dy() == expectedHeight

	if !heightMatch {
		t.Errorf("height mismatch. exepcted=%v, actual=%v", expectedHeight, destBounds.Dy())
	}
}

func TestChangeLineSpace(t *testing.T) {
	img := CreateImage(100, 300, color.White)
	FillRect(img, 50, 20, 80, 50, color.Black)
	FillRect(img, 30, 180, 70, 200, color.Black)
	FillRect(img, 60, 250, 100, 280, color.Black)

	opt := ChangeLineSpaceOption{
		WidthRatio:         100,
		HeightRatio:        200,
		LineSpaceScale:     0.1,
		MinSpace:           1,
		MaxRemove:          9999,
		Threshold:          180,
		EmptyLineThreshold: 0,
	}
	testChangeLineSpace(t, img, opt, int(opt.HeightRatio))
}

func TestChangeLineSpaceVerticalLine(t *testing.T) {
	img := CreateImage(100, 300, color.White)
	FillRect(img, 50, 20, 80, 50, color.Black)
	FillRect(img, 30, 180, 70, 200, color.Black)
	FillRect(img, 60, 250, 100, 280, color.Black)
	DrawLine(img, 90, 0, 90, 300, color.Black)

	opt := ChangeLineSpaceOption{
		WidthRatio:         100,
		HeightRatio:        200,
		LineSpaceScale:     0.1,
		MinSpace:           1,
		MaxRemove:          9999,
		Threshold:          180,
		EmptyLineThreshold: 0.02,
	}
	testChangeLineSpace(t, img, opt, int(opt.HeightRatio))
}
