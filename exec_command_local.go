package main

import (
	"os/exec"
	"path/filepath"
	"sync"
)

var (
	// The path to the statically compiled binary (CLI command) based on github.com/crewjam/go-xmlsec
	// See https://github.com/vstojic/cloudfoundry-cflinuxfs3-go-xmlsec
	XmldsigCmdPath, _ = filepath.Abs("./xmldsig")
)

// For the sake of ExecCmd.
type CommandOutput struct {
	output []byte
	err    error
}

func execCommand(cmd string, //wg *sync.WaitGroup,
	output chan<- *CommandOutput) {

	//defer wg.Done() // Need to signal to waitgroup that this goroutine is done

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		output <- &CommandOutput{nil, err}
	} else {
		// do stuff that generates output, err; then when ready to exit function:
		output <- &CommandOutput{out, nil}
	}
}

// ExecCommands executes an arbitrary shell commands/scripts.
func ExecCommand(cmd string) ([]byte, error) {
	output := make(chan *CommandOutput, 1) // initialize an unbuffered channel
	var wg sync.WaitGroup
	wg.Add(1)
	go execCommand(cmd,
		//&wg,
		output)
	wg.Wait()
	o := <-output
	if o.err != nil {
		return nil, o.err
	}
	return o.output, nil
}
