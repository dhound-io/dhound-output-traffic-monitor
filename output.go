package main

type OutputLine struct {
	EventTimeUtcNumber int64
	Line               string
}

type Output struct {
	Options *Options
	Input   chan []*OutputLine
}

func (output *Output) Run() {
	for logs := range output.Input {
		output._processInput(logs)
	}
}
