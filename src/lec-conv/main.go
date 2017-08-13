package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"path"
	"runtime"
	"sync"
	"time"

	img "lec/image"
)

//-----------------------------------------------------------------------------
// Log
//-----------------------------------------------------------------------------

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().Format("15:04:05") + " " + string(bytes))
}

//-----------------------------------------------------------------------------
// Work
//-----------------------------------------------------------------------------

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
	files, err := img.ListImages(srcDir)
	if err != nil {
		log.Println(err)
		return
	}

	// add works
	for _, file := range files {
		workChan <- Work{srcDir, file.Name(), false}
	}
}

func work(worker Worker, config *Config, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	changeLineSpaceFilter := img.NewChangeLineSpaceFilter(img.ChangeLineSpaceOption{
		WidthRatio:         float64(config.width),
		HeightRatio:        float64(config.height),
		LineSpaceScale:     0.1,
		MinSpace:           1,
		MaxRemove:          9999,
		Threshold:          180,
		EmptyLineThreshold: config.emptyLineThreshold,
	})

	for {
		work := <-worker.workChan
		if work.quit {
			break
		}

		log.Printf("[READ] %v\n", work.filename)

		src, err := img.LoadImage(path.Join(work.dir, work.filename))
		if err != nil {
			log.Printf("Error : %v : %v\n", work.filename, err)
			continue
		}

		// process image
		// --------------
		var dest image.Image

		// change line space
		result := changeLineSpaceFilter.Run(img.NewFilterSource(src, work.filename))
		result.Log()
		dest = result.Img()

		// resize
		dest = img.ResizeImage(dest, config.width, config.height)

		// save dest Img
		// ---------------
		filename := config.FormatDestFilename(work.filename)
		err = img.SaveJpeg(dest, config.destDir, filename, 80)
		if err != nil {
			log.Printf("Error : %v : %v\n", filename, err)
			continue
		}
	}
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	cfgFilename := flag.String("cfg", "", "configuration filename")
	srcDir := flag.String("src", "", "source directory")
	destDir := flag.String("dest", "", "dest directory")
	flag.Parse()

	// Print usage
	if flag.NFlag() == 1 && flag.Arg(1) == "help" {
		flag.Usage()
		return
	}

	// create Config
	config := NewConfig(*cfgFilename, *srcDir, *destDir)
	config.Print()

	// set maxCPU
	runtime.GOMAXPROCS(config.maxCPU)

	// Create channels
	workChan := make(chan Work, 100)
	finChan := make(chan bool)

	// WaitGroup
	wg := sync.WaitGroup{}

	// start collector
	go collectImages(workChan, finChan, config.srcDir)

	// start workers
	for i := 0; i < config.maxCPU; i++ {
		worker := Worker{workChan}
		wg.Add(1)
		go work(worker, config, &wg)
	}

	// wait for collector finish
	<-finChan

	// finish workers
	for i := 0; i < config.maxCPU; i++ {
		workChan <- Work{"", "", true}
	}

	wg.Wait()
}
