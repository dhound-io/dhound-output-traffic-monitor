package main

import (
	"flag"
)

type Options struct {
	LogFile          string
	NetworkInterface string
	Pprof            string
	Version          bool
	Verbose          bool
	FlushTimeOutSec  int32
}

func (options *Options) ParseArguments() {

	flag.StringVar(&options.LogFile, "log-file", "syslog", "network events output: syslog, console, <path to a custom file>; default: console")
	flag.StringVar(&options.NetworkInterface, "eth", options.NetworkInterface, "listen to a particular network interface; default: listen to all active network interfaces")
	flag.BoolVar(&options.Verbose, "verbose", options.Verbose, "log more detailed and debug information; default: false")
	flag.BoolVar(&options.Version, "version", options.Version, "dhound output traffic monitor version")
	flag.StringVar(&options.Pprof, "pprof", options.Pprof, "(for internal using) profiling option")
	flag.Parse()

	options.FlushTimeOutSec = 60
}
