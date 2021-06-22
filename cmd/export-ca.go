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
	"encoding/json"
	"fmt"
	"hlfd/cmd/os_exec_utils"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var exportCaFlags struct {
	CAName     string
	ExportPath string
}

// exportCaCmd represents the export command
var exportCaCmd = &cobra.Command{
	Use:   "ca",
	Short: "Export ca config",
	Long:  `Export ca config`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunRoot()
		preRunExportCA()
	},
	Run: func(cmd *cobra.Command, args []string) {
		exportCA()
	},
}

func init() {
	exportCmd.AddCommand(exportCaCmd)
	exportCaCmd.Flags().StringVarP(&exportCaFlags.CAName, "name", "n", "", "Name of the CA to export")
	exportCaCmd.Flags().StringVarP(&exportCaFlags.ExportPath, "out", "o", "", "Output path to save exported files")

	exportCaCmd.MarkFlagRequired("name")
	exportCaCmd.MarkFlagRequired("out")
}

func preRunExportCA() {
	if exportCaFlags.CAName == "" {
		err := fmt.Errorf("ca name cannot be empty")
		cobra.CheckErr(err)
	}

	if exportCaFlags.ExportPath == "" {
		exportCaFlags.ExportPath = path.Join(hlfdPath, exportCommonFolder, caDepFolder)
	}

	// Make export dir
	fullPath := path.Join(exportCaFlags.ExportPath, exportCaFlags.CAName)
	err := os.MkdirAll(fullPath, commonFilUmask)
	cobra.CheckErr(err)

}

func exportCA() {
	log.Println("Exporting Ca...", exportCaFlags)
	tgtPath := path.Join(exportCaFlags.ExportPath, exportCaFlags.CAName)
	// Get ca info frm src
	srcInfoPath := path.Join(hlfdPath, caDepFolder, exportCaFlags.CAName, caInfoFileName)

	caInfoB := readFileBytes(srcInfoPath)
	var u CAInfo
	err := json.Unmarshal(caInfoB, &u)
	cobra.CheckErr(err)

	if u.TLSEnabled {
		// Cpy tls cert to target
		cmds := []string{
			`cp ` + u.TlsCertPath + ` ` + path.Join(tgtPath, "."),
		}
		_, err = os_exec_utils.ExecMultiCommand(cmds)
		cobra.CheckErr(err)
	}

	// Export ca info to target
	u.TlsCertPath = "./" + caTlsCertFileName // NOTE: Make sure ./ is added, it is used programatically
	m, err := json.MarshalIndent(u, " ", "    ")
	cobra.CheckErr(err)
	writeBytesToFile(caInfoFileName, tgtPath, m)

	// Pack ca folder
	cmds := []string{
		// If not cd into dir, then all dirs inbetween gets added (just the dir, not other files)
		`cd ` + exportCaFlags.ExportPath,
		`tar -czvf ` + exportCaFlags.CAName + `.tar ` + exportCaFlags.CAName, // Only compressed from cafoldername
	}
	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)

	// Delete folder
	cmds = []string{
		`rm -rf ` + tgtPath,
	}
	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)
}
