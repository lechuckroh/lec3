package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"sync"

	"lec/lecimg"
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

func collectImages(workChan chan<- IWork,
	finChan chan<- bool,
	config *Config,
	destDir string,
	filters []lecimg.Filter) {
	defer func() {
		finChan <- true
	}()

	srcFilename := config.src.filename
	srcFileInfo, err := os.Stat(srcFilename)
	if err != nil {
		log.Fatal(err)
	}

	addImageWorks := func(dir string, removeSrc bool) {
		// List image files
		files, err := lecimg.ListImages(dir)
		if err != nil {
			log.Fatal(err)
		}

		// add works
		for _, file := range files {
			workChan <- FilterWork{
				srcDir:    dir,
				filename:  file.Name(),
				destDir:   destDir,
				width:     config.width,
				height:    config.height,
				filters:   filters,
				removeSrc: removeSrc,
			}
		}
	}

	if srcFileInfo.IsDir() {
		addImageWorks(srcFilename, false)
	} else {
		ext := lecimg.GetExt(srcFilename)
		if ext == ".zip" {
			os.MkdirAll(destDir, os.ModePerm)
			extractDir, _ := ioutil.TempDir(destDir, "_temp_")
			log.Printf(extractDir)
			if err := Unzip(srcFilename, extractDir); err != nil {
				log.Fatal(err)
			}
			addImageWorks(extractDir, true)
		}
	}
}

func processWorks(worker Worker, wg *sync.WaitGroup) {
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

func getDestDirInfo(config *Config) DestDirInfo {
	srcFilename := config.src.filename
	destFilename := config.FormatDestFilename(srcFilename)
	destFormat := lecimg.GetExt(destFilename)
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

func createPdf(srcDir string,
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

	log.Printf("[WRITE] %s", path.Join(destDir, filename))
	if err := CreateImagePdf(srcDir, destDir, filename, opt); err != nil {
		log.Fatal(err)
	}
	log.Printf("Done.")
}

func createZip(srcDir string, destDir string, filename string) {
	log.Printf("[WRITE] %s", path.Join(destDir, filename))
	if err := CreateImageZip(srcDir, destDir, filename); err != nil {
		log.Fatal(err)
	}
	log.Printf("Done.")
}

func startWorks(config *Config) {
	srcFilename := config.src.filename
	exists, _ := lecimg.Exists(srcFilename)
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
	var filters []lecimg.Filter
	for _, filterOption := range config.filterOptions {
		filters = append(filters, filterOption.filter)
	}

	// Destination information
	destInfo := getDestDirInfo(config)
	if destInfo.temp {
		defer os.RemoveAll(destInfo.dir)
	}

	// start source images collector
	go collectImages(workChan, finChan, config, destInfo.dir, filters)

	// start workers
	for i := 0; i < config.maxProcess; i++ {
		worker := Worker{workChan}
		wg.Add(1)
		go processWorks(worker, &wg)
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
		createZip(destInfo.dir,
			config.dest.dir,
			destInfo.filename)
	case ".pdf":
		metaData := GetMetaData(path.Base(srcFilename))
		createPdf(destInfo.dir,
			config.dest.dir,
			destInfo.filename,
			metaData,
			config.quality,
			config.showEdgePoint)
	}
}
