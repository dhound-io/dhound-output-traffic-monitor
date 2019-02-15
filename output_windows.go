// +build windows

package main

import (
	"log"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func (output *Output) _processInput(lines []string) {
	// debug("Output started %d", len(lines))
	if lines != nil && len(lines) > 0 {
		log.SetOutput(&lumberjack.Logger{
			Filename:   output.Options.Out,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		})

		for _, line := range lines {
			log.Print(line)
		}
	}
	// debug("Output finished %d", len(lines))
}
