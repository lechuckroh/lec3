package main

import (
	"image"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type FilterWork struct {
	srcDir    string
	filename  string
	index     int
	destDir   string
	width     int
	height    int
	filters   []Filter
	removeSrc bool
}

func (w FilterWork) Run() bool {
	log.Printf("[READ] %v\n", w.filename)

	srcPath := filepath.Join(w.srcDir, w.filename)
	src, err := LoadImage(srcPath)
	if err != nil {
		log.Printf("Error : %v : %v\n", w.filename, err)
		return false
	}

	// run filters
	var dest image.Image
	for _, filter := range w.filters {
		result := filter.Run(NewFilterSource(src, w.filename, w.index))
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
	dest = ResizeImage(dest, w.width, w.height, true)

	// save dest Image
	filename := strings.ToLower(GetBaseWithoutExt(w.filename)) + ".jpg"
	err = SaveJpeg(dest, w.destDir, filename, 80)
	if err != nil {
		log.Printf("Error : %v : %v\n", filename, err)
		return false
	}

	if w.removeSrc {
		os.Remove(srcPath)
	}

	return true
}

func (w FilterWork) IsQuit() bool {
	return false
}
