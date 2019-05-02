// +build windows

package main

func (output *Output) Init() {

}

func (output *Output) _processInput(lines []*OutputLine) {
	for _, line := range lines {
		debug(line.Line)
	}
}
