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

	"github.com/spf13/cobra"
)

// Flags
var orgInfoFlags struct {
	Name string
}

// orgInfoCmd represents the ca command
var orgInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show info about HLF organization",
	Long:  `Show info about HLF organization`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunRoot()
		preRunOrgInfo()
	},
	Run: func(cmd *cobra.Command, args []string) {
		orgInfo()
	},
}

func init() {
	orgCmd.AddCommand(orgInfoCmd)
	orgInfoCmd.Flags().StringVarP(&orgInfoFlags.Name, "name", "n", "", "Org name (required)")

	// Required
	orgInfoCmd.MarkFlagRequired("name")
}

func preRunOrgInfo() {

}

func orgInfo() {
	// Load org
	orgInfo := loadOrgInfo(orgInfoFlags.Name)

	m, err := json.MarshalIndent(orgInfo, " ", "    ")
	cobra.CheckErr(err)

	fmt.Println(string(m))
}
