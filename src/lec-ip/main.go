package main

import (
	"flag"

	llog "lec/log"
)

func getConfig() *Config {
	cfgFilename := flag.String("cfg", "", "configuration filename")
	srcDir := flag.String("src", "./", "source directory")
	destDir := flag.String("dest", "./output", "dest directory")
	watch := flag.Bool("watch", false, "watch directory files update")
	flag.Parse()

	if *cfgFilename == "" {
		return nil
	}

	// create Config
	config := NewConfig(*cfgFilename, *srcDir, *destDir, *watch)
	return config
}

func main() {
	llog.SetLogPattern(llog.TimeOnly)

	config := getConfig()
	if config == nil || (flag.NFlag() == 1 && flag.Arg(1) == "help") {
		flag.Usage()
		return
	}

	config.Print()

	startWorks(config)
}
