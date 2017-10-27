package lecimg

import (
	"image"

	"image/color"
	"image/draw"

	"github.com/mitchellh/mapstructure"
)

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
type WatermarkOption struct {
	Text     string
	Location string
}

func NewWatermarkOption(m map[string]interface{}) (*WatermarkOption, error) {
	option := WatermarkOption{}

	err := mapstructure.Decode(m, &option)
	if err != nil {
		return nil, err
	}

	return &option, nil
}

type WatermarkResult struct {
	image image.Image
}

func (r WatermarkResult) Img() image.Image {
	return r.image
}

func (r WatermarkResult) Log() {
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
type WatermarkFilter struct {
	option WatermarkOption
}

// Create WatermarkFilter instance
func NewWatermarkFilter(option WatermarkOption) *WatermarkFilter {
	return &WatermarkFilter{option: option}
}

// Implements Filter.Run()
func (f WatermarkFilter) Run(s *FilterSource) FilterResult {
	img := f.run(s.image)
	return WatermarkResult{img}
}

// actual watermark implementation
func (f WatermarkFilter) run(src image.Image) image.Image {
	text := f.option.Text
	textLength := len(text)
	if textLength == 0 {
		return src
	}

	bounds := src.Bounds()
	dest := image.NewRGBA(bounds)

	x, y := 0, 0
	fontWidth := 8
	fontHeight := 16
	baselineMargin := 4
	textWidth := textLength * fontWidth
	switch f.option.Location {
	case "TL":
		x = 0
		y = fontHeight
		break
	case "TC":
		x = (bounds.Dx() - textWidth) / 2
		y = fontHeight
		break
	case "TR":
		x = bounds.Dx() - textWidth
		y = fontHeight
		break
	case "CL":
		x = 0
		y = (bounds.Dy() - fontHeight) / 2
		break
	case "CC":
		x = (bounds.Dx() - textWidth) / 2
		y = (bounds.Dy() - fontHeight) / 2
		break
	case "CR":
		x = bounds.Dx() - textWidth
		y = (bounds.Dy() - fontHeight) / 2
		break
	case "BL":
		x = 0
		y = bounds.Dy() - baselineMargin
	case "BC":
		x = (bounds.Dx() - textWidth) / 2
		y = bounds.Dy() - baselineMargin
	case "BR":
		x = bounds.Dx() - textWidth
		y = bounds.Dy() - baselineMargin
		break
	}
	draw.Draw(dest, dest.Bounds(), src, image.ZP, draw.Src)

	DrawLabelBold8x16(dest, x, y, text, color.Gray{128})
	return dest
}
