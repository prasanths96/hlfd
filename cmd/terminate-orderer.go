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
	"os"
	"path"

	"github.com/spf13/cobra"
)

var terminateOrdererFlags struct {
	Name  string
	Quiet bool // Quitely ignore known acceptable errors
}

//
var terminateOrdererCmd = &cobra.Command{
	Use:   "orderer",
	Short: "Terminates Orderer container.",
	Long:  `Terminates Hyperledger Fabric Orderer container.`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunTerminateOrderer()
	},
	Run: func(cmd *cobra.Command, args []string) {
		terminateOrderer()
	},
}

func init() {
	terminateCmd.AddCommand(terminateOrdererCmd)
	terminateOrdererCmd.Flags().StringVarP(&terminateOrdererFlags.Name, "name", "n", "", "Orderer name")
	terminateOrdererCmd.Flags().BoolVarP(&terminateOrdererFlags.Quiet, "quiet", "q", false, "Quietly returns if orderer not found")

	// Required
	terminateOrdererCmd.MarkFlagRequired("name")
}

func preRunTerminateOrderer() {

}

func terminateOrderer() {
	fmt.Println("Terminating Orderer...", terminateOrdererFlags)
	// Check if folder exists
	fullPath := path.Join(hlfdPath, ordererDepFolder, terminateOrdererFlags.Name)
	_, err := os.Stat(fullPath)
	// Return quietly if no such exists (if quiet mode is active)
	if err != nil && terminateOrdererFlags.Quiet {
		return
	}
	cobra.CheckErr(err)

	// execute(fullPath, "docker-compose", "down", "-v", "--rmi", "local")
	cmds := []string{
		// Run docker-compose stop
		`cd ` + fullPath,
		`docker-compose down -v --rmi local`,

		// Remove folder
		`cd ` + path.Join(hlfdPath, ordererDepFolder),
		`sudo rm -rf ` + terminateOrdererFlags.Name,
	}
	// execute(path.Join(hlfdPath, ordererDepFolder), "sudo", "rm", "-rf", terminateOrdererFlags.Name)

	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)
}
