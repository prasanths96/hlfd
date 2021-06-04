package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

func throwOtherThanFileExistError(err error) {
	if err != nil && !os.IsExist(err) {
		cobra.CheckErr(err)
	}
}

func throwIfFileExists(path string) {
	_, err := os.Stat(path)
	if err == nil {
		fmt.Println(path, "already exists")
		os.Exit(1)
	}
}

func writeBytesToFile(fileName string, pathS string, dataB []byte) {
	// f, err := os.Create(path.Join(pathS, fileName))
	err := ioutil.WriteFile(path.Join(pathS, fileName), dataB, commonFilUmask)
	cobra.CheckErr(err)
}

func appendStringToFile(fileName string, pathS string, data string) {
	file, err := os.OpenFile(path.Join(pathS, fileName), os.O_APPEND|os.O_WRONLY, 0644)
	cobra.CheckErr(err)
	defer file.Close()
	_, err = file.WriteString(data)
	cobra.CheckErr(err)
}

func execute(dir string, comdS string, args ...string) {
	comd := exec.Command(comdS, args...)
	if dir != "" {
		comd.Dir = dir
	}
	var out outstream
	comd.Stdout = out
	comd.Stderr = out
	// out, err := comd.CombinedOutput()
	// fmt.Println(string(out))
	// cobra.CheckErr(err)

	// err := comd.Run()
	// cobra.CheckErr(err)

	// stdout, err := comd.StdoutPipe()
	// cobra.CheckErr(err)

	err := comd.Start()
	cobra.CheckErr(err)

	// _, err = ioutil.ReadAll(stdout)
	// cobra.CheckErr(err)

	err = comd.Wait()
	cobra.CheckErr(err)
}

func execAndGetOutput(dir string, comdS string, args ...string) (out []byte) {
	comd := exec.Command(comdS, args...)
	if dir != "" {
		comd.Dir = dir
	}
	out, err := comd.Output()
	if err != nil {
		fmt.Println(out)
		cobra.CheckErr(err)
	}
	return
}

type outstream struct{}

func (out outstream) Write(p []byte) (int, error) {
	fmt.Print(string(p))
	return len(p), nil
}

func isCmdExists(comdS string) (ok bool) {
	comd := exec.Command("which", comdS)
	out, _ := comd.CombinedOutput()
	// cobra.CheckErr(err) // Exit status 1 for command not exist
	if len(out) != 0 { // If not empty
		ok = true
	}
	return
}
