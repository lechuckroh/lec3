package main

import (
	"fmt"
	"os"
)

func printUsage() {
	fmt.Printf("Usage: %s ip|conv [arguments...]", os.Args[0])
}

func main() {
	SetLogPattern(TimeOnly)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "ip":
		ip := ImageProcess{}
		ip.run()
	case "conv":
		conv := Convert{}
		conv.run()
	}
}
