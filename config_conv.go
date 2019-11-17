package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/olebedev/config"
)

type SrcOptionConv struct {
	filename string
}
type DestOptionConv struct {
	dir      string
	filename string
}

// Format returns extension of filename in lowercase
func (opt DestOptionConv) Format() string {
	return GetExt(opt.filename)
}

type FilterOptionConv struct {
	name   string
	filter Filter
}

// ConfigConv defines configuration
type ConfigConv struct {
	src           SrcOptionConv
	dest          DestOptionConv
	width         int
	height        int
	quality       int
	showEdgePoint bool
	maxProcess    int
	filterOptions []FilterOptionConv
}

// LoadYaml loads *.yaml file
func (c *ConfigConv) LoadYaml(filename string) {
	cfg, err := config.ParseYamlFile(filename)
	if err != nil {
		log.Printf("Error : Failed to parse %v : %v\n", filename, err)
		return
	}

	log.Printf("Load: %v\n", filename)

	c.src.filename = cfg.UString("src.filename", "./")
	c.dest.dir = cfg.UString("dest.dir", "./output")
	c.dest.filename = cfg.UString("dest.filename", "${filename}")
	c.width = cfg.UInt("width", -1)
	c.height = cfg.UInt("height", -1)
	c.quality = cfg.UInt("quality", 100)
	c.showEdgePoint = cfg.UBool("showEdgePoint", false)
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

func (c *ConfigConv) addFilterOption(name string, options map[string]interface{}) {
	var err error
	var filter Filter

	switch name {
	case "changeLineSpace":
		if option, err := NewChangeLineSpaceOption(options); err == nil {
			filter = NewChangeLineSpaceFilter(*option)
		}
	case "resize":
		if option, err := NewResizeOption(options); err == nil {
			filter = NewResizeFilter(*option)
		}
	case "watermark":
		if option, err := NewWatermarkOption(options); err == nil {
			filter = NewWatermarkFilter(*option)
		}
	default:
		log.Printf("Unhandled filter name : %v\n", name)
	}

	if filter != nil {
		filterOption := FilterOptionConv{
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

// FormatDestFilename formats destFilename pattern
func (c *ConfigConv) FormatDestFilename(dirname string) string {
	result := c.dest.filename

	base := filepath.Base(dirname)
	result = strings.Replace(result, "${filename}", base, -1)

	baseFilename := GetBaseWithoutExt(base)
	result = strings.Replace(result, "${baseFilename}", baseFilename, -1)
	return result
}

// Print displays configurations
func (c *ConfigConv) Print() {
	log.Printf("src.filename : %v\n", c.src.filename)
	log.Printf("dest.dir : %v\n", c.dest.dir)
	log.Printf("dest.filename : %v\n", c.dest.filename)
	log.Printf("size : (%v, %v)\n", c.width, c.height)
	log.Printf("showEdgePoint : %v\n", c.showEdgePoint)
	log.Printf("quality : %v%%\n", c.quality)
	log.Printf("maxProcess : %v\n", c.maxProcess)
	fmt.Printf("filters : %v\n", len(c.filterOptions))
}

// NewConfigConv creates an instance of ConfigConv
func NewConfigConv(cfgFilename string, srcDir string, destDir string) *ConfigConv {
	cfg := ConfigConv{}

	if cfgFilename != "" {
		cfg.LoadYaml(cfgFilename)
	}
	if srcDir != "" {
		cfg.src.filename = srcDir
	}
	if destDir != "" {
		cfg.dest.dir = destDir
	}

	return &cfg
}
