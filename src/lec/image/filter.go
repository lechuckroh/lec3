package image

import "image"

// FilterSource is a source of filter
type FilterSource struct {
	image    image.Image
	filename string
}

// NewFilterSource creates an instance of FilterSource
func NewFilterSource(image image.Image, filename string) *FilterSource {
	return &FilterSource{image, filename}
}

// FilterResult is a result of filter operation
type FilterResult interface {
	Img() image.Image
	Log()
}

// Filter is an interface for filter operation
type Filter interface {
	Run(src *FilterSource) FilterResult
}
