package leclog

import (
	"fmt"
	"log"
	"time"
)

type timeOnlyLogWriter struct {
}

func (writer timeOnlyLogWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().Format("15:04:05") + " " + string(bytes))
}

type Pattern int

const (
	TimeOnly Pattern = iota
)

func SetLogPattern(pattern Pattern) {
	switch pattern {
	case TimeOnly:
		log.SetFlags(0)
		log.SetOutput(new(timeOnlyLogWriter))
	default:
		log.Printf("Unhandled logFormat: %v\n", pattern)
	}
}
