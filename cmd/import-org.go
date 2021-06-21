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

var importOrgFlags struct {
	FilePath string
}

// importOrgCmd represents the import command
var importOrgCmd = &cobra.Command{
	Use:   "org",
	Short: "Import org config",
	Long:  `Import org config`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunRoot()
		preRunImportOrg()
	},
	Run: func(cmd *cobra.Command, args []string) {
		importOrg()
	},
}

func init() {
	importCmd.AddCommand(importOrgCmd)
	importOrgCmd.Flags().StringVarP(&importOrgFlags.FilePath, "file", "f", "", "Path to exported org .tar file")
	importOrgCmd.MarkFlagRequired("file")
}

func preRunImportOrg() {
	if importOrgFlags.FilePath == "" {
		err := fmt.Errorf("file path cannot be empty")
		cobra.CheckErr(err)
	}
	// Resolving full path
	var err error
	importOrgFlags.FilePath, err = filepath.Abs(importOrgFlags.FilePath)
	cobra.CheckErr(err)

	// Make import dir
	fullPath := path.Join(hlfdPath, importCommonFolder, orgCommonFolder)
	err = os.MkdirAll(fullPath, commonFilUmask)
	cobra.CheckErr(err)
}

func importOrg() {
	log.Println("Importing Org...", importOrgFlags)
	// Export file to imports folder
	fullPath := path.Join(hlfdPath, importCommonFolder, orgCommonFolder)
	cmds := []string{
		`cd ` + fullPath,
		`tar -xvf ` + importOrgFlags.FilePath,
	}
	_, err := os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)
}
