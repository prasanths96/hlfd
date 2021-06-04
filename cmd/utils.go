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

func execute(dir string, comdS string, args ...string) {
	comd := exec.Command(comdS, args...)
	if dir != "" {
		comd.Dir = dir
	}
	// stdin, err := comd.StdinPipe()
	// go func() {
	// 	defer stdin.Close()
	// 	io.WriteString(stdin, "an old falcon")
	// }()
	// stdout, err := comd.StdoutPipe()
	// cobra.CheckErr(err)

	out, err := comd.CombinedOutput()
	fmt.Println(string(out))
	cobra.CheckErr(err)
}
