package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type OutputLine struct {
	EventTimeUtcNumber int64
	Line               string
}

type Output struct {
	Options *Options
	Input   chan []*OutputLine
}

func (output *Output) Init() {
	options := output.Options

	if len(options.LogFile) > 0 && options.LogFile != "syslog" && options.LogFile != "console" {
		fmt.Println("See output in  " + options.LogFile)

		logDir := filepath.Dir(options.LogFile)
		var err error
		if _, err = os.Stat(logDir); os.IsNotExist(err) {
			err = os.MkdirAll(logDir, 0765)
		}

		if err != nil {
			emitLine(logLevel.important, "Failed to create dir '%s' for log file. Error: %s", logDir, err)
		}
	}
}

func (output *Output) Run() {
	for logs := range output.Input {
		output._processInput(logs)
	}
}
