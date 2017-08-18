package main

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"sync"

	limg "lec/image"

	"github.com/signintech/gopdf"
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
	filters []limg.Filter) {
	defer func() {
		finChan <- true
	}()

	// List image files
	files, err := limg.ListImages(config.src.dir)
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

type DestInfo struct {
	dir      string
	filename string
	format   string
	temp     bool
}

func getDestInfo(config *Config) DestInfo {
	srcDir := config.src.dir
	destFilename := config.FormatDestFilename(srcDir)
	destFormat := limg.GetExt(destFilename)
	destDir := config.dest.dir
	isTempDir := destFormat != ""
	if isTempDir {
		os.MkdirAll(destDir, os.ModePerm)
		destDir, _ = ioutil.TempDir(destDir, "_temp_")
	}
	return DestInfo{
		dir:      destDir,
		filename: destFilename,
		format:   destFormat,
		temp:     isTempDir,
	}
}

func createPdf(srcDir string, destDir string, filename string, width int, height int) {
	files, err := limg.ListImages(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	pdf := gopdf.GoPdf{}
	rect := gopdf.Rect{
		W: float64(width),
		H: float64(height),
	}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: rect})

	for _, file := range files {
		pdf.AddPage()
		// TODO: keep aspect ratio
		pdf.Image(path.Join(srcDir, file.Name()), 0, 0, &rect)
	}

	pdf.WritePdf(path.Join(destDir, filename))
}

func createZip(srcDir string, destDir string, filename string) {
	files, err := limg.ListImages(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	newFile, err := os.Create(path.Join(destDir, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	zipWriter := zip.NewWriter(newFile)
	defer zipWriter.Close()

	for _, file := range files {
		data, err := ioutil.ReadFile(path.Join(srcDir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		f, err := zipWriter.Create(file.Name())
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.Write(data)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func startWorks(config *Config) {
	// set maxProcess
	runtime.GOMAXPROCS(config.maxProcess)

	// channels
	workChan := make(chan IWork, 100)
	finChan := make(chan bool)
	wg := sync.WaitGroup{}

	// filters
	var filters []limg.Filter
	for _, filterOption := range config.filterOptions {
		filters = append(filters, filterOption.filter)
	}

	// Destination information
	destInfo := getDestInfo(config)
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

	log.Printf("[WRITE] %s", path.Join(config.dest.dir, destInfo.filename))

	// Create output
	switch destInfo.format {
	case ".cbz", ".zip":
		createZip(destInfo.dir, config.dest.dir, destInfo.filename)
	case ".pdf":
		createPdf(destInfo.dir, config.dest.dir, destInfo.filename, config.width, config.height)
	}
	log.Printf("Done.")
}
