package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"lec/lecimg"
	"lec/lecio"

	"github.com/olebedev/config"
)

type SrcOption struct {
	filename string
}
type DestOption struct {
	dir      string
	filename string
}

// Format returns extension of filename in lowercase
func (opt DestOption) Format() string {
	return lecio.GetExt(opt.filename)
}

type FilterOption struct {
	name   string
	filter lecimg.Filter
}

// Config defines configuration
type Config struct {
	src           SrcOption
	dest          DestOption
	width         int
	height        int
	quality       int
	showEdgePoint bool
	maxProcess    int
	filterOptions []FilterOption
}

// LoadYaml loads *.yaml file
func (c *Config) LoadYaml(filename string) {
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

func (c *Config) addFilterOption(name string, options map[string]interface{}) {
	var err error
	var filter lecimg.Filter

	switch name {
	case "changeLineSpace":
		if option, err := lecimg.NewChangeLineSpaceOption(options); err == nil {
			filter = lecimg.NewChangeLineSpaceFilter(*option)
		}
	case "resize":
		if option, err := lecimg.NewResizeOption(options); err == nil {
			filter = lecimg.NewResizeFilter(*option)
		}
	case "watermark":
		if option, err := lecimg.NewWatermarkOption(options); err == nil {
			filter = lecimg.NewWatermarkFilter(*option)
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

// FormatDestFilename formats destFilename pattern
func (c *Config) FormatDestFilename(dirname string) string {
	result := c.dest.filename

	base := filepath.Base(dirname)
	result = strings.Replace(result, "${filename}", base, -1)

	baseFilename := lecio.GetBaseWithoutExt(base)
	result = strings.Replace(result, "${baseFilename}", baseFilename, -1)
	return result
}

// Print displays configurations
func (c *Config) Print() {
	log.Printf("src.filename : %v\n", c.src.filename)
	log.Printf("dest.dir : %v\n", c.dest.dir)
	log.Printf("dest.filename : %v\n", c.dest.filename)
	log.Printf("size : (%v, %v)\n", c.width, c.height)
	log.Printf("showEdgePoint : %v\n", c.showEdgePoint)
	log.Printf("quality : %v%%\n", c.quality)
	log.Printf("maxProcess : %v\n", c.maxProcess)
	fmt.Printf("filters : %v\n", len(c.filterOptions))
}

// NewConfig creates an instance of Config
func NewConfig(cfgFilename string, srcDir string, destDir string) *Config {
	cfg := Config{}

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
