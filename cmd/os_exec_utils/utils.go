package os_exec_utils

import (
	"fmt"
	"os/exec"
)

func ExecAndGetOutput(dir string, comdS string, args ...string) (out []byte, err error) {
	comd := exec.Command(comdS, args...)
	if dir != "" {
		comd.Dir = dir
	}
	out, err = comd.Output()
	return
}

func ExecMultiCommand(commands []string) (outB []byte, err error) {
	if len(commands) < 1 {
		return
	}

	// Finally exit from bash
	// commands = commands + ` && ` + `exit`

	cmd := exec.Command("bash")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	var outStream = &outstream{
		out: make([]byte, 0),
	}
	var errStream = &outstream{
		out: make([]byte, 0),
	}
	cmd.Stdout = outStream
	cmd.Stderr = errStream

	for i := 0; i < len(commands); i++ {
		_, err = fmt.Fprintf(stdin, "%s\n", commands[i])
		if err != nil {
			return
		}
	}
	stdin.Close()

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("outstream: %v, errstream: %v, err: %v", string(outStream.out), string(errStream.out), err.Error())
		return
	}

	outB = outStream.out
	return
}

func ExecMultiCommandNoWait(commands []string) (outB []byte, err error) {
	if len(commands) < 1 {
		return
	}

	cmd := exec.Command("bash")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	var outStream = &outstream{
		out: make([]byte, 0),
	}
	var errStream = &outstream{
		out: make([]byte, 0),
	}
	cmd.Stdout = outStream
	cmd.Stderr = errStream

	for i := 0; i < len(commands); i++ {
		_, err = fmt.Fprintf(stdin, "nohup %s\n", commands[i])
		if err != nil {
			return
		}
	}
	stdin.Close()

	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("outstream: %v, errstream: %v, err: %v", string(outStream.out), string(errStream.out), err.Error())
		return
	}
	outB = outStream.out
	return
}

type outstream struct {
	out []byte
}

func (out *outstream) Write(p []byte) (int, error) {
	fmt.Print(string(p))
	out.out = append(out.out, p...)
	return len(p), nil
}
