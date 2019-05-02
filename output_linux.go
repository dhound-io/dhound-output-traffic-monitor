// +build !windows

package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"path/filepath"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func (output *Output) Init() {
	options := output.Options

	// make directory
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

	// configure output
	if output.Options.LogFile == "syslog" {
		logwriter, e := syslog.New(syslog.LOG_NOTICE, "dhound-output-traffic-monitor ")
		if e == nil {
			//logwriter.Info()
			log.SetOutput(logwriter)
		}
	} else if output.Options.LogFile == "console" {

	} else {
		log.SetOutput(&lumberjack.Logger{
			Filename:   output.Options.LogFile,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		})
	}
}

func (output *Output) _processInput(lines []*OutputLine) {

	if len(lines) > 0 {
		for _, line := range lines {
			log.Print(line.Line)
		}
	}
}
