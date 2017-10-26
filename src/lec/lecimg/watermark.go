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
	FontSize int
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
	if len(text) == 0 {
		return src
	}

	bounds := src.Bounds()
	dest := image.NewRGBA(bounds)

	x, y := 0, 0
	switch f.option.Location {
	case "TL":
		break
	case "TC":
		break
	case "TR":
		break
	case "CL":
		break
	case "CC":
		break
	case "CR":
		break
	case "BL":
		x = 0
		y = bounds.Dy()
	case "BC":
		// TODO: align center
		x = bounds.Dx() / 2
		y = bounds.Dy()
	case "BR":
		break
	}
	draw.Draw(dest, dest.Bounds(), src, image.ZP, draw.Src)

	// TODO: use font size option
	DrawLabel(dest, x, y, text, color.Black)
	return dest
}
