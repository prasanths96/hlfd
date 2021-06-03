package cmd

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

func checkOtherThanFileExistsError(err error) {
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		cobra.CheckErr(err)
	}
}

func writeBytesToFile(fileName string, pathS string, dataB []byte) {
	// f, err := os.Create(path.Join(pathS, fileName))
	err := ioutil.WriteFile(path.Join(pathS, fileName), dataB, commonFilUmask)
	cobra.CheckErr(err)
}
