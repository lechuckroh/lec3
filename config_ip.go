package main

import (
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
	var filter Filter

	switch name {
	case "deskew":
		option, err := NewDeskewOption(options)
		if err == nil {
			filter = NewDeskewFilter(*option)
		} else {
			log.Printf("Failed to read filter : %v : %v\n", name, err)
		}
	case "deskewED":
		option, err := NewDeskewEDOption(options)
		if err == nil {
			filter = NewDeskewEDFilter(*option)
		} else {
			log.Printf("Failed to read filter : %v : %v\n", name, err)
		}
	case "autoCrop":
		option, err := NewAutoCropOption(options)
		if err == nil {
			filter = NewAutoCropFilter(*option)
		} else {
			log.Printf("Failed to read filter : %v : %v\n", name, err)
		}
	case "autoCropED":
		option, err := NewAutoCropEDOption(options)
		if err == nil {
			filter = NewAutoCropEDFilter(*option)
		} else {
			log.Printf("Failed to read filter : %v : %v\n", name, err)
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
}

func (c *ConfigIP) Print() {
	fmt.Printf("src.dir : %v\n", c.src.dir)
	fmt.Printf("dest.dir : %v\n", c.dest.dir)
	fmt.Printf("watch : %v\n", c.watch)
	fmt.Printf("maxProcess : %v\n", c.maxProcess)
	fmt.Printf("filters : %v\n", len(c.filterOptions))
}
