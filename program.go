package main

import (
	"net/http"
	"sync"

	svc "github.com/judwhite/go-svc/svc"
)

import _ "net/http/pprof"

type Program struct {
	Options *Options
	Wg      sync.WaitGroup
	Quit    chan struct{}
}

func (program *Program) Init(env svc.Environment) error {

	return nil
}

func (program *Program) InternalRun() {

	options := program.Options

	emitLine(logLevel.important, "Dhound output traffic monitor %s started. Options: log-file='%s' eth='%s' verbose='%t' protocol='%s'", Version, options.LogFile, options.NetworkInterface, options.Verbose, options.Protocol)

	if len(options.Pprof) > 0 {
		go func() {
			emit(logLevel.verbose, "run profiler on http://%s/debug/pprof/\n", options.Pprof)
			err := http.ListenAndServe(options.Pprof, nil)
			if err != nil {
				emit(logLevel.important, "failed running profiler: %s \n", err.Error())
			}
		}()
	}

	netStat := &NetStatManager{
		Options: options,
	}
	netStat.Init()

	sysProcessManager := &SysProcessManager{}
	sysProcessManager.Init()

	output := &Output{
		Input:   make(chan []*OutputLine),
		Options: options,
	}
	output.Init()

	networkEventEnricher := &NetworkEventEnricher{
		Input:      make(chan *NetworkEvent),
		Output:     output.Input,
		NetStat:    netStat,
		SysManager: sysProcessManager,
	}
	networkEventEnricher.Init()

	networkMonitor := &NetworkMonitor{
		Options: options,
		Output:  networkEventEnricher.Input,
	}

	sysProcessManager.Run()
	netStat.Run()
	networkMonitor.Run()
	networkEventEnricher.Run()
	output.Run()
}

func (program *Program) Start() error {
	// The Start method must not block, or Windows may assume your service failed
	// to start. Launch a Goroutine here to do something interesting/blocking.

	program.Quit = make(chan struct{})

	program.Wg.Add(1)
	go func() {

		program.InternalRun()

		<-program.Quit
		// debug("Quit signal received...")
		program.Wg.Done()
	}()

	return nil
}

func (program *Program) Stop() error {
	// The Stop method is invoked by stopping the Windows service, or by pressing Ctrl+C on the console.
	// This method may block, but it's a good idea to finish quickly or your process may be killed by
	// Windows during a shutdown/reboot. As a general rule you shouldn't rely on graceful shutdown.

	// emitLine(logLevel.verbose, "Stopping...")
	close(program.Quit)
	program.Wg.Wait()
	emitLine(logLevel.verbose, "Stopped.")

	return nil
}
