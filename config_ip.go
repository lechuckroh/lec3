package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/olebedev/config"
)

type SrcOptionIP struct {
	dir       string
	recursive bool
}
type DestOptionIP struct {
	dir string
}

type FilterOptionIP struct {
	name   string
	filter Filter
}

type ConfigIP struct {
	src           SrcOptionIP
	dest          DestOptionIP
	watch         bool
	watchDelay    int
	maxProcess    int
	filterOptions []FilterOptionIP
}

func (c *ConfigIP) LoadYaml(filename string) {
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

func (c *ConfigIP) addFilterOption(name string, options map[string]interface{}) {
	var err error
	var filter Filter

	switch name {
	case "deskew":
		if option, err := NewDeskewOption(options); err == nil {
			filter = NewDeskewFilter(*option)
		}
	case "deskewED":
		if option, err := NewDeskewEDOption(options); err == nil {
			filter = NewDeskewEDFilter(*option)
		}
	case "autoCrop":
		if option, err := NewAutoCropOption(options); err == nil {
			filter = NewAutoCropFilter(*option)
		}
	case "autoCropED":
		if option, err := NewAutoCropEDOption(options); err == nil {
			filter = NewAutoCropEDFilter(*option)
		}
	default:
		log.Printf("Unhandled filter name : %v\n", name)
	}

	if filter != nil {
		filterOption := FilterOptionIP{
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

func (c *ConfigIP) Print() {
	fmt.Printf("src.dir : %v\n", c.src.dir)
	fmt.Printf("dest.dir : %v\n", c.dest.dir)
	fmt.Printf("watch : %v\n", c.watch)
	fmt.Printf("maxProcess : %v\n", c.maxProcess)
	fmt.Printf("filters : %v\n", len(c.filterOptions))
}

func NewConfigIP(cfgFilename string, srcDir string, destDir string, watch bool) *ConfigIP {
	cfg := ConfigIP{}

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
