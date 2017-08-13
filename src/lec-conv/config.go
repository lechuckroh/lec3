package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	limg "lec/image"

	"github.com/olebedev/config"
)

type SrcOption struct {
	dir string
}
type DestOption struct {
	dir      string
	filename string
}

type FilterOption struct {
	name   string
	filter limg.Filter
}

// Config defines configuration
type Config struct {
	src           SrcOption
	dest          DestOption
	width         int
	height        int
	quality       int
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

	c.src.dir = cfg.UString("src.dir", "./")
	c.dest.dir = cfg.UString("dest.dir", "./output")
	c.dest.filename = cfg.UString("dest.filename", "${base}.jpg")
	c.width = cfg.UInt("width", -1)
	c.height = cfg.UInt("height", -1)
	c.quality = cfg.UInt("quality", 100)
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
	var filter limg.Filter

	switch name {
	case "changeLineSpace":
		if option, err := limg.NewChangeLineSpaceOption(options); err == nil {
			filter = limg.NewChangeLineSpaceFilter(*option)
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
func (c *Config) FormatDestFilename(filename string) string {
	result := strings.Replace(c.dest.filename, "${filename}", filename, -1)
	base := strings.ToLower(limg.GetBase(filename))
	result = strings.Replace(result, "${base}", base, -1)
	return result
}

// Print displays configurations
func (c *Config) Print() {
	log.Printf("src.dir : %v\n", c.src.dir)
	log.Printf("dest.dir : %v\n", c.dest.dir)
	log.Printf("dest.filename : %v\n", c.dest.filename)
	log.Printf("size : (%v, %v)\n", c.width, c.height)
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
		cfg.src.dir = srcDir
	}
	if destDir != "" {
		cfg.dest.dir = destDir
	}

	return &cfg
}
