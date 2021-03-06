package lecimg

import (
	"image"
	"image/color"
	"testing"
)

func testDeskew(t *testing.T,
	img image.Image,
	option DeskewOption,
	rotatedAngleMin, rotatedAngleMax float32) {
	// Run Filter
	result := NewDeskewFilter(option).Run(NewFilterSource(img, "filename", 0))
	rotatedAngle := result.(DeskewResult).rotatedAngle

	// Test result image size
	if !InRangef32(rotatedAngle, rotatedAngleMin, rotatedAngleMax) {
		t.Errorf("angle mismatch. exepcted=(%v ~ %v), actual=%v",
			rotatedAngleMin, rotatedAngleMax, rotatedAngle)
	}
}

func TestDeskewCCW(t *testing.T) {
	img := CreateImage(400, 700, color.White)
	FillRect(img, 50, 50, 350, 650, color.Black)
	rotatedImg := RotateImage(img, -1.4, color.White)

	// Run Filter
	option := DeskewOption{
		MaxRotation:         2,
		IncrStep:            0.2,
		Threshold:           220,
		EmptyLineMaxDotRate: 0.01,
		DetectToleranceRate: 0.005,
	}
	testDeskew(t, rotatedImg, option, 1.2, 1.6)
}

func TestDeskewCW(t *testing.T) {
	img := CreateImage(400, 700, color.White)
	FillRect(img, 50, 50, 550, 650, color.Black)
	rotatedImg := RotateImage(img, 1.4, color.White)

	// Run Filter
	option := DeskewOption{
		MaxRotation:          2,
		IncrStep:             0.2,
		Threshold:            220,
		EmptyLineMaxDotCount: 0,
	}
	testDeskew(t, rotatedImg, option, -1.6, -1.2)
}
