/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"hlfd/cmd/os_exec_utils"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

var importCaFlags struct {
	FilePath string
}

// importCaCmd represents the import command
var importCaCmd = &cobra.Command{
	Use:   "ca",
	Short: "Import ca config",
	Long:  `Import ca config`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunRoot()
		preRunImportCA()
	},
	Run: func(cmd *cobra.Command, args []string) {
		importCA()
	},
}

func init() {
	importCmd.AddCommand(importCaCmd)
	importCaCmd.Flags().StringVarP(&importCaFlags.FilePath, "file", "f", "", "Path to exported ca .tar file")
	importCaCmd.MarkFlagRequired("file")
}

func preRunImportCA() {
	if importCaFlags.FilePath == "" {
		err := fmt.Errorf("file path cannot be empty")
		cobra.CheckErr(err)
	}
	// Resolving full path
	var err error
	importCaFlags.FilePath, err = filepath.Abs(importCaFlags.FilePath)
	cobra.CheckErr(err)

	// Make import dir
	fullPath := path.Join(hlfdPath, importCommonFolder, caDepFolder)
	err = os.MkdirAll(fullPath, commonFilUmask)
	cobra.CheckErr(err)
}

func importCA() {
	log.Println("Importing Ca...", importCaFlags)
	// Export file to imports folder
	fullPath := path.Join(hlfdPath, importCommonFolder, caDepFolder)
	cmds := []string{
		`cd ` + fullPath,
		`tar -xvf ` + importCaFlags.FilePath,
	}
	_, err := os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)
}
