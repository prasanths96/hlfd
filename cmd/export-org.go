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

var exportOrgFlags struct {
	OrgName    string
	ExportPath string
}

// exportOrgCmd represents the export command
var exportOrgCmd = &cobra.Command{
	Use:   "org",
	Short: "Export org config",
	Long:  `Export org config`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunRoot()
		preRunExportOrg()
	},
	Run: func(cmd *cobra.Command, args []string) {
		exportOrg()
	},
}

func init() {
	exportCmd.AddCommand(exportOrgCmd)
	exportOrgCmd.Flags().StringVarP(&exportOrgFlags.OrgName, "name", "n", "", "Name of the org to export")
	exportOrgCmd.Flags().StringVarP(&exportOrgFlags.ExportPath, "out", "o", "", "Output path to save exported files")

	exportOrgCmd.MarkFlagRequired("name")
	exportOrgCmd.MarkFlagRequired("out")
}

func preRunExportOrg() {
	if exportOrgFlags.OrgName == "" {
		err := fmt.Errorf("org name cannot be empty")
		cobra.CheckErr(err)
	}

	if exportOrgFlags.ExportPath == "" {
		exportOrgFlags.ExportPath = path.Join(hlfdPath, exportCommonFolder, orgCommonFolder)
	}

	// Make export dir
	fullPath := path.Join(exportOrgFlags.ExportPath, exportOrgFlags.OrgName)
	err := os.MkdirAll(fullPath, commonFilUmask)
	cobra.CheckErr(err)

}

func exportOrg() {
	log.Println("Exporting Org...", exportOrgFlags)
	tgtPath := path.Join(exportOrgFlags.ExportPath, exportOrgFlags.OrgName)
	// Get org info frm src
	srcMspPath := path.Join(hlfdPath, orgCommonFolder, exportOrgFlags.OrgName, mspFolder)
	srcInfoPath := path.Join(hlfdPath, orgCommonFolder, exportOrgFlags.OrgName, orgInfoFileName)

	orgInfoB := readFileBytes(srcInfoPath)
	var u OrgInfo
	err := json.Unmarshal(orgInfoB, &u)
	cobra.CheckErr(err)

	casFolder := path.Join(tgtPath, caDepFolder)
	// Copy necessary ca tls certs into target
	for i := 0; i < len(u.CaInfo); i++ {
		ca := &u.CaInfo[i]
		// Create ca folder
		caPath := path.Join(casFolder, ca.CaName)
		err := os.MkdirAll(caPath, commonFilUmask)
		cobra.CheckErr(err)

		// Clear caClientHomePath
		ca.CaClientHomePath = ""

		if ca.CaTlsCertPath == "" {
			continue
		}

		// Copy tls cert
		newTlsPath := path.Join(caPath, caTlsCertFileName)
		cmds := []string{
			`cp ` + ca.CaTlsCertPath + ` ` + newTlsPath,
		}
		_, err = os_exec_utils.ExecMultiCommand(cmds)
		cobra.CheckErr(err)

		// Change tls path
		ca.CaTlsCertPath = `./` + path.Join(caDepFolder, ca.CaName, caTlsCertFileName)
	}

	// Copy msp folder
	cmds := []string{
		`cp -rf ` + srcMspPath + ` ` + path.Join(tgtPath, "."),
	}
	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)

	// Change msp dir
	u.MspDir = `./` + mspFolder

	// Export org info to target
	m, err := json.MarshalIndent(u, " ", "    ")
	cobra.CheckErr(err)
	writeBytesToFile(orgInfoFileName, tgtPath, m)

	// Pack ca folder
	cmds = []string{
		`tar -czvf ` + tgtPath + `.tar ` + tgtPath,
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
