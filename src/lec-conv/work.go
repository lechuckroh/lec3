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

	// List image files
	files, err := lecimg.ListImages(config.src.dir)
	if err != nil {
		log.Println(err)
		return
	}

	// add works
	for _, file := range files {
		workChan <- FilterWork{
			srcDir:   config.src.dir,
			filename: file.Name(),
			destDir:  destDir,
			width:    config.width,
			height:   config.height,
			filters:  filters,
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
	srcDir := config.src.dir
	destFilename := config.FormatDestFilename(srcDir)
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
	srcDir := config.src.dir
	exists, _ := lecimg.Exists(srcDir)
	if !exists {
		log.Printf("Directory not found : %s", srcDir)
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
		metaData := GetMetaData(path.Base(srcDir))
		createPdf(destInfo.dir,
			config.dest.dir,
			destInfo.filename,
			metaData,
			config.quality,
			config.showEdgePoint)
	}
}
