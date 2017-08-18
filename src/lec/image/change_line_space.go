package image

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"sort"

	"github.com/mitchellh/mapstructure"
)

// ChangeLineSpaceOption contains options for changeLineSpace filter
type ChangeLineSpaceOption struct {
	WidthRatio         float64
	HeightRatio        float64
	LineSpaceScale     float64
	MinSpace           int
	MaxRemove          int
	Threshold          uint32
	EmptyLineThreshold float64
	DebugMode          bool
}

func NewChangeLineSpaceOption(m map[string]interface{}) (*ChangeLineSpaceOption, error) {
	option := ChangeLineSpaceOption{}

	err := mapstructure.Decode(m, &option)
	if err != nil {
		return nil, err
	}

	return &option, nil
}

type ChangeLineSpaceResult struct {
	image image.Image
	rect  image.Rectangle
}

func (r ChangeLineSpaceResult) Img() image.Image {
	return r.image
}

func (r ChangeLineSpaceResult) Log() {
}

// ----------------------------------------------------------------------------

type lineRange struct {
	start        int
	end          int
	height       int
	targetHeight int
	emptyLine    bool
}

func (r *lineRange) calc(scale float64, minHeight, maxRemove int) {
	r.height = r.end - r.start + 1
	if !r.emptyLine || r.height <= minHeight {
		r.targetHeight = r.height
	} else {
		if maxRemove > 0 {
			r.targetHeight = Max(minHeight, int(float64(r.height)*scale+0.5))
			if removed := r.height - r.targetHeight; removed > maxRemove {
				r.targetHeight = r.height - maxRemove
			}
		}
	}
}

func (r lineRange) getReducedHeight() int {
	return r.height - r.targetHeight
}

func (r lineRange) String() string {
	return fmt.Sprintf("(%4d-%4d) h:%d, targetH: %d, empty: %v",
		r.start, r.end, r.height, r.targetHeight, r.emptyLine)
}

type lineRanges []lineRange

func (r lineRanges) Len() int {
	return len(r)
}

// Less sorts in descending order by targetHeight
func (r lineRanges) Less(i, j int) bool {
	return r[i].targetHeight > r[j].targetHeight
}
func (r lineRanges) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// ----------------------------------------------------------------------------

type ChangeLineSpaceFilter struct {
	option ChangeLineSpaceOption
}

func NewChangeLineSpaceFilter(option ChangeLineSpaceOption) *ChangeLineSpaceFilter {
	return &ChangeLineSpaceFilter{option: option}
}

func getBrightness(r, g, b uint32) uint32 {
	return (r + g + b) / 3
}

// getLineRanges returns list of text and empty lines
func (f ChangeLineSpaceFilter) getLineRanges(src image.Image) lineRanges {
	bounds := src.Bounds()
	srcWidth, srcHeight := bounds.Dx(), bounds.Dy()
	threshold16 := f.option.Threshold * 256

	var ranges lineRanges
	var r lineRange

	maxDotCount := int(f.option.EmptyLineThreshold)
	if f.option.EmptyLineThreshold < 1 {
		maxDotCount = int(float64(srcWidth) * f.option.EmptyLineThreshold)
	}
	for y := 0; y < srcHeight; y++ {
		emptyLine := true
		dotCount := 0
		for x := 0; x < srcWidth; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			brightness := getBrightness(r, g, b)
			if brightness < threshold16 {
				dotCount++
				if dotCount >= maxDotCount {
					emptyLine = false
					break
				}
			}
		}

		if emptyLine {
			if y == 0 {
				r = lineRange{start: y, end: y, emptyLine: true}
			} else {
				if r.emptyLine {
					r.end = y
				} else {
					ranges = append(ranges, r)
					r = lineRange{start: y, end: y, emptyLine: true}
				}
			}
		} else {
			if y == 0 {
				r = lineRange{start: y, end: y, emptyLine: false}
			} else {
				if r.emptyLine {
					ranges = append(ranges, r)
					r = lineRange{start: y, end: y, emptyLine: false}
				} else {
					r.end = y
				}
			}
		}

	}
	ranges = append(ranges, r)
	return ranges
}

func (f ChangeLineSpaceFilter) processLineRanges(ranges lineRanges, width int) int {
	targetHeight := 0
	rangeCount := len(ranges)
	emptyRangeCount := 0
	for i := 0; i < rangeCount; i++ {
		r := &ranges[i]
		r.calc(f.option.LineSpaceScale, f.option.MinSpace, f.option.MaxRemove)
		targetHeight += r.targetHeight
		if r.emptyLine {
			emptyRangeCount++
		}
	}

	minTargetHeight := int(f.option.HeightRatio * float64(width) / f.option.WidthRatio)
	if f.option.DebugMode {
		log.Printf("min targetHeight: %v", minTargetHeight)
	}

	loop := 0
	maxLoopCount := 5
	for targetHeight < minTargetHeight && loop < maxLoopCount {
		totalReducedHeight := 0
		for i := 0; i < rangeCount; i++ {
			r := ranges[i]
			if r.emptyLine {
				totalReducedHeight += r.getReducedHeight()
			}
		}

		totalInc := 0
		if totalReducedHeight > 0 {
			totalEmptyLinesToInc := minTargetHeight - targetHeight
			for i := 0; i < len(ranges); i++ {
				r := &ranges[i]
				if r.emptyLine {
					reducedHeight := r.getReducedHeight()
					heightToInc := Min(reducedHeight,
						reducedHeight*totalEmptyLinesToInc/totalReducedHeight)
					inc := int(float32(heightToInc) + 0.5)
					r.targetHeight = r.targetHeight + inc
					totalInc += inc
				}
			}
		}

		targetHeight = 0
		for i := 0; i < rangeCount; i++ {
			r := &ranges[i]
			targetHeight += r.targetHeight
		}

		if f.option.DebugMode {
			log.Printf("[%d] totalInc=%d, targetHeight=%d", loop, totalInc, targetHeight)
		}

		loop++

		if totalInc <= emptyRangeCount {
			if remainInc := minTargetHeight - targetHeight; remainInc != 0 {
				sortedRanges := make(lineRanges, rangeCount)
				copy(sortedRanges, ranges)
				sort.Sort(sortedRanges)
				for i := 0; i < rangeCount && remainInc != 0; i++ {
					if r := sortedRanges[i]; r.emptyLine {
						if remainInc > 0 {
							r.targetHeight++
							targetHeight++
							remainInc--
						} else {
							r.targetHeight--
							targetHeight--
							remainInc++
						}
					}
				}
			}
			break
		}
	}
	if f.option.DebugMode {
		log.Printf("TargetHeight: %v", targetHeight)
	}

	return targetHeight
}

func (f ChangeLineSpaceFilter) Run(s *FilterSource) FilterResult {
	img, rect := f.run(s.image)
	return &ChangeLineSpaceResult{img, rect}
}

func (f ChangeLineSpaceFilter) run(src image.Image) (image.Image, image.Rectangle) {
	ranges := f.getLineRanges(src)
	rangeCount := len(ranges)

	if rangeCount <= 1 {
		return src, src.Bounds()
	} else {
		width := src.Bounds().Dx()
		targetHeight := f.processLineRanges(ranges, width)

		bounds := image.Rect(0, 0, width, targetHeight)
		dest := CreateImage(width, targetHeight, color.White)
		destY := 0
		for i := 0; i < rangeCount; i++ {
			r := &ranges[i]
			rangeHeight := r.height
			rangeTargetHeight := r.targetHeight

			if rangeHeight > 0 && rangeTargetHeight > 0 {
				srcRect := image.Rect(0, r.start, width, r.start+rangeTargetHeight)
				subImage := src.(interface {
					SubImage(r image.Rectangle) image.Image
				}).SubImage(srcRect)

				destRect := image.Rect(0, destY, width, targetHeight)
				draw.Draw(dest, destRect, subImage, image.ZP, draw.Src)

				destY -= rangeHeight - rangeTargetHeight
			}
		}
		return dest, bounds
	}
}
