package lecimg

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/disintegration/gift"
	"github.com/mitchellh/mapstructure"
)

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
type DeskewOption struct {
	MaxRotation          float32 // max rotation angle (0 <= value <= 360)
	IncrStep             float32 // rotation angle increment step (0 <= value <= 360)
	EmptyLineMaxDotCount int
	EmptyLineMaxDotRate  float32 // max dot count rate (0 <= value < 1.0)
	DebugOutputDir       string
	DebugMode            bool
	Threshold            uint8   // min brightness of space (0~255)
	DetectToleranceRate  float32 // max dot count diff rate (0 <= value < 1.0)
}

func NewDeskewOption(m map[string]interface{}) (*DeskewOption, error) {
	option := DeskewOption{}

	err := mapstructure.Decode(m, &option)
	if err != nil {
		return nil, err
	}

	return &option, nil
}

type DeskewResult struct {
	image        image.Image
	filename     string
	rotatedAngle float32
}

func (r DeskewResult) Img() image.Image {
	return r.image
}

func (r DeskewResult) Log() {
	if r.rotatedAngle != 0 {
		log.Printf("[ROTATE] %v : %.1f", r.filename, r.rotatedAngle)
	}
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type DeskewFilter struct {
	option DeskewOption
}

// Create DeskewFilter instance
func NewDeskewFilter(option DeskewOption) *DeskewFilter {
	return &DeskewFilter{option}
}

// Implements Filter.Run()
func (f DeskewFilter) Run(s *FilterSource) FilterResult {
	resultImage, rotatedAngle := f.run(s.image, s.filename)
	return DeskewResult{resultImage, s.filename, rotatedAngle}
}

// actual deskew implementation
func (f DeskewFilter) run(src image.Image, name string) (image.Image, float32) {
	bounds := src.Bounds()
	var rgba *image.RGBA

	switch src.(type) {
	case *image.RGBA:
		rgba = src.(*image.RGBA)
	default:
		rgba = image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(rgba, bounds, src, bounds.Min, draw.Src)
	}

	if angle := f.detectAngle(rgba, name); angle != 0 {
		return f.rotateImage(rgba, angle), angle
	}
	return src, 0
}

// Rotate image
func (f DeskewFilter) rotateImage(src image.Image, angle float32) image.Image {
	bounds := src.Bounds()
	width, height := CalcRotatedSize(bounds.Dx(), bounds.Dy(), angle)
	dest := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dest, dest.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	rotateFilter := gift.New(gift.Rotate(angle, color.White, gift.CubicInterpolation))
	rotateFilter.Draw(dest, src)
	return dest
}

func (f DeskewFilter) detectAngle(src *image.RGBA, name string) float32 {
	minNonEmptyLineCount := f.calcNonEmptyLineCount(src, 0, name)
	tolerance := int(float32(src.Bounds().Dx()) * f.option.DetectToleranceRate)

	// increase rotation angle by incrStep
	detectedAngle := float32(0)

	positiveDir := true
	negativeDir := true

	incrStep := f.option.IncrStep
	if incrStep > 0 {
		for angle := incrStep; angle <= f.option.MaxRotation; angle += incrStep {
			if positiveDir {
				nonEmptyLineCount := f.calcNonEmptyLineCount(src, angle, name)

				if nonEmptyLineCount < minNonEmptyLineCount {
					minNonEmptyLineCount = nonEmptyLineCount
					detectedAngle = angle
				} else if nonEmptyLineCount >= minNonEmptyLineCount+tolerance {
					positiveDir = false
				}
			}

			if angle > 0 && negativeDir {
				nonEmptyLineCount := f.calcNonEmptyLineCount(src, -angle, name)

				if nonEmptyLineCount < minNonEmptyLineCount {
					minNonEmptyLineCount = nonEmptyLineCount
					detectedAngle = -angle
				} else if nonEmptyLineCount >= minNonEmptyLineCount+tolerance {
					negativeDir = false
				}
			}
		}
	}

	return detectedAngle
}

func (f DeskewFilter) calcNonEmptyLineCount(src *image.RGBA, angle float32, name string) int {
	dy, _ := Sincosf32(angle)
	bounds := src.Bounds()

	thresholdSum := uint32(f.option.Threshold) * 256 * 3
	nonEmptyLineCount := 0
	width, height := bounds.Dx(), bounds.Dy()

	emptyLineMaxDotCount := f.option.EmptyLineMaxDotCount
	if emptyLineMaxDotCount == 0 {
		emptyLineMaxDotCount = int(f.option.EmptyLineMaxDotRate * float32(width))
	}

	for y := 0; y < height; y++ {
		yPos := float32(y)
		dotCount := 0

		for x := 0; x < width; x++ {
			yPosInt := int(yPos)
			if yPosInt < 0 || yPosInt >= height {
				break
			}

			if r, g, b, _ := src.At(x, yPosInt).RGBA(); r+g+b <= thresholdSum {
				dotCount++
			}

			yPos += dy
		}

		if emptyLineMaxDotCount < dotCount {
			nonEmptyLineCount++
		}
	}

	if f.option.DebugMode {
		log.Printf("angle=%v, nonEmptyLineCount=%v\n", angle, nonEmptyLineCount)
	}

	return nonEmptyLineCount
}
