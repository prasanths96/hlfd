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
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// orgCmd represents the create command
var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Org Command is used to do operations related to HLF orgs.",
	Long:  `Org Command is used to do operations related to HLF orgs.`,
}

func init() {
	rootCmd.AddCommand(orgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// orgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// orgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loadOrgInfo(orgName string) (loadedOrgInfo OrgInfo) {
	// See if org available locally
	orgPath := path.Join(hlfdPath, orgCommonFolder, orgName)
	// throwIfFileNotExist(orgPath)
	_, err := os.Stat(orgPath)
	if err != nil {
		if !os.IsNotExist(err) {
			cobra.CheckErr(err)
		}

		// Check if org in imports
		orgPath = path.Join(hlfdPath, importCommonFolder, orgCommonFolder, orgName)
		throwIfFileNotExist(orgPath)
	}

	orgInfoPath := path.Join(orgPath, orgInfoFileName)
	err = json.Unmarshal(readFileBytes(orgInfoPath), &loadedOrgInfo)
	cobra.CheckErr(err)

	// Resolve ca paths to absolute paths
	for i := 0; i < len(loadedOrgInfo.CaInfo); i++ {
		ca := &loadedOrgInfo.CaInfo[i]
		if ca.CaTlsCertPath == "" {
			continue
		}

		if ca.CaTlsCertPath[0] == '.' {
			ca.CaTlsCertPath = strings.Replace(ca.CaTlsCertPath, ".", orgPath, 1)
		}
	}
	// Resolve msp path to absolute path
	if loadedOrgInfo.MspDir[0] == '.' {
		loadedOrgInfo.MspDir = strings.Replace(loadedOrgInfo.MspDir, ".", orgPath, 1)
	}

	return
}

func selectCaFromList(caName string, caInfo []OrgCAInfo) (selectedCA OrgCAInfo) {
	switch caName {
	case "":
		// Choose first ca from list
		selectedCA = caInfo[0]
	default:
		found := false
		for _, v := range caInfo {
			if v.CaName == caName {
				selectedCA = v
				found = true
				break
			}
		}
		if !found {
			err := fmt.Errorf("ca-name: %v not found in the org: %v", caName, depOrdererFlags.OrgName)
			cobra.CheckErr(err)
		}
	}
	return
}
