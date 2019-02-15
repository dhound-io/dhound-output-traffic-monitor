package main

import (
	"time"
)

type Output struct {
	Options *Options
	Input   chan []string
}

func (output *Output) Init() {

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			for logs := range output.Input {
				output._processInput(logs)
			}
		}
	}()
}

func (output *Output) Run() {

}
