package main

import (
	"image"
	"log"
	"path"
	"reflect"
	"runtime"
	"sync"

	limg "lec/image"
)

// Work contains an information of image file to process
type Work struct {
	dir      string
	filename string
	quit     bool
}

// Worker is channel of Work
type Worker struct {
	workChan <-chan Work
}

func collectImages(workChan chan<- Work, finChan chan<- bool, srcDir string) {
	defer func() {
		finChan <- true
	}()

	// List image files
	files, err := limg.ListImages(srcDir)
	if err != nil {
		log.Println(err)
		return
	}

	// add works
	for _, file := range files {
		workChan <- Work{srcDir, file.Name(), false}
	}
}

func work(worker Worker, filters []limg.Filter, config *Config, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for {
		work := <-worker.workChan
		if work.quit {
			break
		}

		log.Printf("[READ] %v\n", work.filename)

		src, err := limg.LoadImage(path.Join(work.dir, work.filename))
		if err != nil {
			log.Printf("Error : %v : %v\n", work.filename, err)
			continue
		}

		// run filters
		var dest image.Image
		for _, filter := range filters {
			result := filter.Run(limg.NewFilterSource(src, work.filename))
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
		dest = limg.ResizeImage(dest, config.width, config.height)

		// save dest Img
		// ---------------
		filename := config.FormatDestFilename(work.filename)
		err = limg.SaveJpeg(dest, config.dest.dir, filename, 80)
		if err != nil {
			log.Printf("Error : %v : %v\n", filename, err)
			continue
		}
	}
}

func startWorks(config *Config) {
	// set maxProcess
	runtime.GOMAXPROCS(config.maxProcess)

	// Create channels
	workChan := make(chan Work, 100)
	finChan := make(chan bool)

	// WaitGroup
	wg := sync.WaitGroup{}

	// start collector
	go collectImages(workChan, finChan, config.src.dir)

	var filters []limg.Filter
	for _, filterOption := range config.filterOptions {
		filters = append(filters, filterOption.filter)
	}

	// start workers
	for i := 0; i < config.maxProcess; i++ {
		worker := Worker{workChan}
		wg.Add(1)
		go work(worker, filters, config, &wg)
	}

	// wait for collector finish
	<-finChan

	// finish workers
	for i := 0; i < config.maxProcess; i++ {
		workChan <- Work{"", "", true}
	}

	wg.Wait()
}
