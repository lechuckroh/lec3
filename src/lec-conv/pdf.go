package main

import (
	"image"
	"image/color"
	"image/draw"
	"path"
	"time"

	"lec/lecimg"

	"github.com/signintech/gopdf"
)

type PdfOption struct {
	Title         string
	Author        string
	Quality       int
	ShowEdgePoint bool
}

func toPdfPoint(pixel int) float64 {
	return float64(pixel) / 128 * 72
}

func CreateImagePdf(srcDir string, destDir string, filename string, opt PdfOption) error {
	files, err := lecimg.ListImages(srcDir)
	if err != nil {
		return err
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{})

	for _, file := range files {
		img, err := lecimg.LoadImage(path.Join(srcDir, file.Name()))
		if err != nil {
			return err
		}

		imageBounds := img.Bounds()
		width := imageBounds.Dx()
		height := imageBounds.Dy()
		rect := gopdf.Rect{
			W: toPdfPoint(width),
			H: toPdfPoint(height),
		}

		// show edge point
		if opt.ShowEdgePoint {
			rgba := image.NewRGBA(image.Rect(0, 0, width, height))
			draw.Draw(rgba, imageBounds, img, imageBounds.Min, draw.Src)

			rgba.Set(0, 0, color.Black)
			rgba.Set(width-1, 0, color.Black)
			rgba.Set(0, height-1, color.Black)
			rgba.Set(width-1, height-1, color.Black)
			img = rgba
		}

		bytes, err := lecimg.ToJpegBytes(img, opt.Quality)
		if err != nil {
			return err
		}

		pdf.AddPageWithOption(gopdf.PageOption{PageSize: rect})
		imgHolder, err := gopdf.ImageHolderByBytes(bytes)
		if err != nil {
			return err
		}
		pdf.ImageByHolder(imgHolder, 0, 0, nil)
	}

	// MetaData
	pdf.SetInfo(gopdf.PdfInfo{
		Title:        opt.Title,
		Author:       opt.Author,
		Creator:      "lec-conv",
		CreationDate: time.Now(),
	})

	pdf.WritePdf(path.Join(destDir, filename))
	return nil
}
