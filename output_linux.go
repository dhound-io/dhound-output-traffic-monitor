// +build !windows

package main

func (output *Output) _processInput(lines []string) {
	// debug("Output started %d", len(lines))
	if lines != nil && len(lines) > 0 {
		if output.Options.Out == "syslog" {
			logwriter, e := syslog.New(syslog.LOG_NOTICE, "Dhound Output Traffic Monitor")
			if e == nil {
				log.SetOutput(logwriter)
			}
		} else {
			log.SetOutput(&lumberjack.Logger{
				Filename:   output.Options.Out,
				MaxSize:    100, // megabytes
				MaxBackups: 3,
				MaxAge:     28, // days
			})
		}
		for _, line := range lines {
			log.Print(line)
		}
	}
	// debug("Output finished %d", len(lines))
}
