package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type IWork interface {
	Run() bool
	IsQuit() bool
}

// Work contains an information of image file to process
type Work struct {
	dir      string
	filename string
	destDir  string
	quit     bool
}

// Worker is channel of IWork
type Worker struct {
	workChan <-chan IWork
}

type Convert struct {
}

func (c *Convert) collectImages(workChan chan<- IWork,
	finChan chan<- bool,
	config *ConfigConv,
	destDir string,
	filters []Filter) {
	defer func() {
		finChan <- true
	}()

	srcFilename := config.src.filename
	srcFileInfo, err := os.Stat(srcFilename)
	if err != nil {
		log.Fatal(err)
	}

	addWork := func(dir string, filename string, index int, removeSrc bool) {
		workChan <- FilterWork{
			srcDir:    dir,
			filename:  filename,
			index:     index,
			destDir:   destDir,
			width:     config.width,
			height:    config.height,
			filters:   filters,
			removeSrc: removeSrc,
		}
	}

	// Add works for the images in the 'dir' directory.
	addImageWorks := func(dir string, removeSrc bool) {
		// List image files
		files, err := ListImages(dir)
		if err != nil {
			log.Fatal(err)
		}

		// add works
		for i, file := range files {
			addWork(dir, file.Name(), i, removeSrc)
		}
	}

	if srcFileInfo.IsDir() {
		addImageWorks(srcFilename, false)
	} else {
		ext := GetExt(srcFilename)
		if ext == ".zip" || ext == ".cbz" {
			os.MkdirAll(destDir, os.ModePerm)
			extractDir, _ := ioutil.TempDir(destDir, "_temp_")

			callback := func(dir, filename string, index int) {
				addWork(extractDir, filename, index, true)
			}

			if err := Unzip(srcFilename, extractDir, callback); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (c *Convert) processWorks(worker Worker, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for {
		work := <-worker.workChan
		if work.IsQuit() {
			break
		}

		work.Run()
	}
}

type DestDirInfo struct {
	dir      string
	filename string
	format   string
	temp     bool
}

func (c *Convert) getDestDirInfo(config *ConfigConv) DestDirInfo {
	srcFilename := config.src.filename
	destFilename := config.FormatDestFilename(srcFilename)
	destFormat := GetExt(destFilename)
	destDir := config.dest.dir
	isTempDir := destFormat != ""
	if isTempDir {
		os.MkdirAll(destDir, os.ModePerm)
		destDir, _ = ioutil.TempDir(destDir, "_temp_")
	}
	return DestDirInfo{
		dir:      destDir,
		filename: destFilename,
		format:   destFormat,
		temp:     isTempDir,
	}
}

func (c *Convert) createPdf(srcDir string,
	destDir string,
	filename string,
	metaData MetaData,
	quality int,
	showEdgePoint bool) {

	opt := PdfOption{
		Title:         metaData.Title,
		Author:        metaData.Author,
		Quality:       quality,
		ShowEdgePoint: showEdgePoint,
	}

	log.Printf("[WRITE] %s", filepath.Join(destDir, filename))
	if err := CreateImagePdf(srcDir, destDir, filename, opt); err != nil {
		log.Fatal(err)
	}
	log.Printf("Done.")
}

func (c *Convert) createZip(srcDir string, destDir string, filename string) {
	log.Printf("[WRITE] %s", filepath.Join(destDir, filename))
	if err := CreateImageZip(srcDir, destDir, filename); err != nil {
		log.Fatal(err)
	}
	log.Printf("Done.")
}

func (c *Convert) run(config *ConfigConv) {
	srcFilename := config.src.filename
	exists, _ := Exists(srcFilename)
	if !exists {
		log.Printf("File not found : %s", srcFilename)
		return
	}

	// set maxProcess
	runtime.GOMAXPROCS(config.maxProcess)

	// channels
	workChan := make(chan IWork, 100)
	finChan := make(chan bool)
	wg := sync.WaitGroup{}

	// filters
	var filters []Filter
	for _, filterOption := range config.filterOptions {
		filters = append(filters, filterOption.filter)
	}

	// Destination information
	destInfo := c.getDestDirInfo(config)
	if destInfo.temp {
		defer os.RemoveAll(destInfo.dir)
	}

	// start source images collector
	go c.collectImages(workChan, finChan, config, destInfo.dir, filters)

	// start workers
	for i := 0; i < config.maxProcess; i++ {
		worker := Worker{workChan}
		wg.Add(1)
		go c.processWorks(worker, &wg)
	}

	// wait for collector finish
	<-finChan

	// finish workers
	for i := 0; i < config.maxProcess; i++ {
		workChan <- QuitWork{}
	}

	wg.Wait()

	// Create output
	switch destInfo.format {
	case ".cbz", ".zip":
		c.createZip(destInfo.dir,
			config.dest.dir,
			destInfo.filename)
	case ".pdf":
		metaData := GetMetaData(filepath.Base(srcFilename))
		c.createPdf(destInfo.dir,
			config.dest.dir,
			destInfo.filename,
			metaData,
			config.quality,
			config.showEdgePoint)
	}
}
