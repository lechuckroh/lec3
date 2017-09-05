package lecimg

import (
	"image"
	"log"

	"github.com/disintegration/gift"
	"github.com/mitchellh/mapstructure"
)

type ResizeOption struct {
	WidthScale  float64
	HeightScale float64
	ScaleCover  bool
}

func NewResizeOption(m map[string]interface{}) (*ResizeOption, error) {
	option := ResizeOption{}

	err := mapstructure.Decode(m, &option)
	if err != nil {
		return nil, err
	}

	return &option, nil
}

type ResizeResult struct {
	image    image.Image
	filename string
	scaled   bool
}

func (r ResizeResult) Img() image.Image {
	return r.image
}

func (r ResizeResult) Log() {
	if !r.scaled {
		log.Printf("Resize skipped : %s\n", r.filename)
	}
}

// ----------------------------------------------------------------------------

type ResizeFilter struct {
	option ResizeOption
}

func NewResizeFilter(option ResizeOption) *ResizeFilter {
	return &ResizeFilter{option: option}
}

func (f ResizeFilter) Run(s *FilterSource) FilterResult {
	if !f.option.ScaleCover && s.index == 0 {
		return ResizeResult{image: s.image, filename: s.filename, scaled: false}
	}

	bbox := s.image.Bounds()
	width := int(f.option.WidthScale * float64(bbox.Dx()))
	height := int(f.option.HeightScale * float64(bbox.Dy()))
	resizedImage := ResizeImage(s.image, width, height, false)
	return ResizeResult{image: resizedImage, filename: s.filename, scaled: true}
}

// ----------------------------------------------------------------------------

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
