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

	"lec/lecimg"
)

// Work represents a job to do
type Work struct {
	dir      string
	filename string
	quit     bool
}

// Worker is a worker to process images.
type Worker struct {
	workChan <-chan Work
}

func collectImages(workChan chan<- Work, finChan chan<- bool, srcDir string, watch bool, watchDelay int) {
	defer func() {
		finChan <- true
	}()

	lastCheckTime := time.Unix(0, 0)
	var files []os.FileInfo
	var err error

	for {
		// List modified image files
		files, lastCheckTime, err = lecimg.ListModifiedImages(srcDir, watchDelay, lastCheckTime)
		if err != nil {
			log.Println(err)
			break
		}

		// add works
		for _, file := range files {
			workChan <- Work{srcDir, file.Name(), false}
		}

		if watch {
			// sleep for a while
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			break
		}
	}
}

func work(worker Worker, filters []lecimg.Filter, destDir string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	for {
		work := <-worker.workChan
		if work.quit {
			break
		}

		log.Printf("[R] %v\n", work.filename)

		src, err := lecimg.LoadImage(filepath.Join(work.dir, work.filename))
		if err != nil {
			log.Printf("Error : %v : %v\n", work.filename, err)
			continue
		}

		// run filters
		var dest image.Image
		for _, filter := range filters {
			result := filter.Run(lecimg.NewFilterSource(src, work.filename))
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
		err = lecimg.SaveJpeg(dest, destDir, work.filename, 80)
		if err != nil {
			log.Printf("Error : %v : %v\n", work.filename, err)
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
	go collectImages(workChan, finChan, config.src.dir, config.watch, config.watchDelay)

	var filters []lecimg.Filter
	for _, filterOption := range config.filterOptions {
		filters = append(filters, filterOption.filter)
	}

	// start workers
	for i := 0; i < config.maxProcess; i++ {
		worker := Worker{workChan}
		wg.Add(1)
		go work(worker, filters, config.dest.dir, &wg)
	}

	// wait for collector finish
	<-finChan

	// finish workers
	for i := 0; i < config.maxProcess; i++ {
		workChan <- Work{"", "", true}
	}

	wg.Wait()
}
