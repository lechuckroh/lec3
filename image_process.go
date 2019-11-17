package main

import (
	"image"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"
)

// IPWork represents a job to do
type IPWork struct {
	dir      string
	filename string
	quit     bool
}

// IPWorker is a worker to process images.
type IPWorker struct {
	workChan <-chan IPWork
}

type ImageProcess struct {
}

func (ip *ImageProcess) collectImages(workChan chan<- IPWork, finChan chan<- bool, srcDir string, watch bool, watchDelay int) {
	defer func() {
		finChan <- true
	}()

	lastCheckTime := time.Unix(0, 0)
	var files []os.FileInfo
	var err error

	for {
		// List modified image files
		files, lastCheckTime, err = ListModifiedImages(srcDir, watchDelay, lastCheckTime)
		if err != nil {
			log.Println(err)
			break
		}

		// add works
		for _, file := range files {
			workChan <- IPWork{srcDir, file.Name(), false}
		}

		if watch {
			// sleep for a while
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			break
		}
	}
}

func (ip *ImageProcess) work(worker IPWorker, filters []Filter, destDir string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for {
		work := <-worker.workChan
		if work.quit {
			break
		}

		log.Printf("[R] %v\n", work.filename)

		src, err := LoadImage(filepath.Join(work.dir, work.filename))
		if err != nil {
			log.Printf("Error : %v : %v\n", work.filename, err)
			continue
		}

		// run filters
		var dest image.Image
		for _, filter := range filters {
			result := filter.Run(NewFilterSource(src, work.filename, -1))
			result.Log()

			resultImg := result.Img()
			if resultImg == nil {
				log.Printf("Filter result is nil. filter: %v\n", reflect.TypeOf(filter))
				break
			}

			dest = resultImg
			src = dest
		}

		// save dest Img
		err = SaveJpeg(dest, destDir, work.filename, 80)
		if err != nil {
			log.Printf("Error : %v : %v\n", work.filename, err)
			continue
		}
	}
}

func (ip *ImageProcess) run(config *ConfigIP) {
	// set maxProcess
	runtime.GOMAXPROCS(config.maxProcess)

	// Create channels
	workChan := make(chan IPWork, 100)
	finChan := make(chan bool)

	// WaitGroup
	wg := sync.WaitGroup{}

	// start collector
	go ip.collectImages(workChan, finChan, config.src.dir, config.watch, config.watchDelay)

	var filters []Filter
	for _, filterOption := range config.filterOptions {
		filters = append(filters, filterOption.filter)
	}

	// start workers
	for i := 0; i < config.maxProcess; i++ {
		worker := IPWorker{workChan}
		wg.Add(1)
		go ip.work(worker, filters, config.dest.dir, &wg)
	}

	// wait for collector finish
	<-finChan

	// finish workers
	for i := 0; i < config.maxProcess; i++ {
		workChan <- IPWork{"", "", true}
	}

	wg.Wait()
}

