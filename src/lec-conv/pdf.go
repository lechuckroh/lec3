package main

import (
	"path"

	"lec/lecimg"

	"github.com/signintech/gopdf"
)

func toPdfPoint(pixel int) float64 {
	return float64(pixel) / 128 * 72
}

func CreateImagePdf(srcDir string,
	destDir string,
	filename string,
	width int,
	height int) error {
	files, err := lecimg.ListImages(srcDir)
	if err != nil {
		return err
	}

	pdf := gopdf.GoPdf{}
	rect := gopdf.Rect{
		W: toPdfPoint(width),
		H: toPdfPoint(height),
	}
	pdf.Start(gopdf.Config{PageSize: rect})

	for _, file := range files {
		pdf.AddPage()
		pdf.Image(path.Join(srcDir, file.Name()), 0, 0, nil)
	}

	pdf.WritePdf(path.Join(destDir, filename))
	return nil
}
