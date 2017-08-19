package main

import (
	"image"
	"log"
	"path"
	"reflect"
	"strings"

	"lec/lecimg"
)

type FilterWork struct {
	srcDir   string
	filename string
	destDir  string
	width    int
	height   int
	filters  []lecimg.Filter
}

func (w FilterWork) Run() bool {
	log.Printf("[READ] %v\n", w.filename)

	src, err := lecimg.LoadImage(path.Join(w.srcDir, w.filename))
	if err != nil {
		log.Printf("Error : %v : %v\n", w.filename, err)
		return false
	}

	// run filters
	var dest image.Image
	for _, filter := range w.filters {
		result := filter.Run(lecimg.NewFilterSource(src, w.filename))
		result.Log()

		resultImg := result.Img()
		if resultImg == nil {
			log.Printf("Filter result is nil. filter: %v\n", reflect.TypeOf(filter))
			break
		}

		dest = resultImg
		src = dest
	}

	// resize
	dest = lecimg.ResizeImage(dest, w.width, w.height)

	// save dest Image
	filename := strings.ToLower(lecimg.GetBaseWithoutExt(w.filename)) + ".jpg"
	err = lecimg.SaveJpeg(dest, w.destDir, filename, 80)
	if err != nil {
		log.Printf("Error : %v : %v\n", filename, err)
		return false
	}

	return true
}

func (w FilterWork) IsQuit() bool {
	return false
}
