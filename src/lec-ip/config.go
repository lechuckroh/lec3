package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"lec/lecimg"

	"github.com/olebedev/config"
)

type SrcOption struct {
	dir       string
	recursive bool
}
type DestOption struct {
	dir string
}

type FilterOption struct {
	name   string
	filter lecimg.Filter
}

type Config struct {
	src           SrcOption
	dest          DestOption
	watch         bool
	watchDelay    int
	maxProcess    int
	filterOptions []FilterOption
}

func (c *Config) LoadYaml(filename string) {
	cfg, err := config.ParseYamlFile(filename)
	if err != nil {
		log.Printf("Error : Failed to parse %v : %v\n", filename, err)
		return
	}

	fmt.Printf("Loading %v\n", filename)

	c.src.dir = cfg.UString("src.dir", "")
	c.src.recursive = cfg.UBool("src.recursive", false)
	c.dest.dir = cfg.UString("dest.dir", "")
	c.watch = cfg.UBool("watch", false)
	c.watchDelay = cfg.UInt("watchDelay", 5)
	c.maxProcess = cfg.UInt("maxProcess", runtime.NumCPU())
	if c.maxProcess <= 0 {
		c.maxProcess = runtime.NumCPU()
	}

	// Load filters
	for i := 0; ; i++ {
		m, err := cfg.Map(fmt.Sprintf("filters.%v", i))
		if err != nil {
			break
		}
		name, ok := m["name"]
		if !ok {
			continue
		}

		c.addFilterOption(name.(string), m["options"].(map[string]interface{}))
	}
}

func (c *Config) addFilterOption(name string, options map[string]interface{}) {
	var err error
	var filter lecimg.Filter

	switch name {
	case "deskew":
		if option, err := lecimg.NewDeskewOption(options); err == nil {
			filter = lecimg.NewDeskewFilter(*option)
		}
	case "deskewED":
		if option, err := lecimg.NewDeskewEDOption(options); err == nil {
			filter = lecimg.NewDeskewEDFilter(*option)
		}
	case "autoCrop":
		if option, err := lecimg.NewAutoCropOption(options); err == nil {
			filter = lecimg.NewAutoCropFilter(*option)
		}
	case "autoCropED":
		if option, err := lecimg.NewAutoCropEDOption(options); err == nil {
			filter = lecimg.NewAutoCropEDFilter(*option)
		}
	default:
		log.Printf("Unhandled filter name : %v\n", name)
	}

	if filter != nil {
		filterOption := FilterOption{
			name:   name,
			filter: filter,
		}
		c.filterOptions = append(c.filterOptions, filterOption)
		fmt.Printf("Filter added : %v\n", name)
	}
	if err != nil {
		log.Printf("Failed to read filter : %v : %v\n", name, err)
	}
}

func (c *Config) Print() {
	fmt.Printf("src.dir : %v\n", c.src.dir)
	fmt.Printf("dest.dir : %v\n", c.dest.dir)
	fmt.Printf("watch : %v\n", c.watch)
	fmt.Printf("maxProcess : %v\n", c.maxProcess)
	fmt.Printf("filters : %v\n", len(c.filterOptions))
}

func NewConfig(cfgFilename string, srcDir string, destDir string, watch bool) *Config {
	cfg := Config{}

	if cfgFilename != "" {
		cfg.LoadYaml(cfgFilename)
	} else {
		// overwrite cfg with command line options
		if srcFlag := flag.Lookup("src"); srcFlag != nil {
			cfg.src.dir = srcDir
		}
		if destFlag := flag.Lookup("dest"); destFlag != nil {
			cfg.dest.dir = destDir
		}
		if watchFlag := flag.Lookup("watch"); watchFlag != nil {
			cfg.watch = watch
		}
	}

	return &cfg
}
