// +build windows

package main

// 	"log"

//lumberjack "gopkg.in/natefinch/lumberjack.v2"

func (output *Output) _processInput(lines []*OutputLine) {

	// debug("Output started %d", len(lines))
	/*if lines != nil && len(lines) > 0 {


		log.SetOutput(&lumberjack.Logger{
			Filename:   output.Options.Out,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		})

		for _, line := range lines {
			log.Print(line)
		}
	}*/

	// time.Unix(event.EventTimeUtcNumber, 0).Format(time.RFC3339)

	for _, line := range lines {
		debug(line.Line)
	}

	// debug("Output finished %d", len(lines))
}
