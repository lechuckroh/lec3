package main

import (
	"flag"

	"lec/leclog"
)

func getConfig() *Config {
	cfgFilename := flag.String("cfg", "", "configuration filename")
	srcDir := flag.String("src", "", "source directory")
	destDir := flag.String("dest", "", "dest directory")
	flag.Parse()

	if *cfgFilename == "" {
		return nil
	}

	// create Config
	config := NewConfig(*cfgFilename, *srcDir, *destDir)
	return config
}

func main() {
	leclog.SetLogPattern(leclog.TimeOnly)

	config := getConfig()
	if config == nil || (flag.NFlag() == 1 && flag.Arg(1) == "help") {
		flag.Usage()
		return
	}

	config.Print()

	startWorks(config)
}
