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

var terminatePeerFlags struct {
	Name  string
	Quiet bool // Quitely ignore known acceptable errors
}

//
var terminatePeerCmd = &cobra.Command{
	Use:   "peer",
	Short: "Terminates Peer container.",
	Long:  `Terminates Hyperledger Fabric Peer container.`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunTerminatePeer()
	},
	Run: func(cmd *cobra.Command, args []string) {
		terminatePeer()
	},
}

func init() {
	terminateCmd.AddCommand(terminatePeerCmd)
	terminatePeerCmd.Flags().StringVarP(&terminatePeerFlags.Name, "name", "n", "", "Peer name")
	terminatePeerCmd.Flags().BoolVarP(&terminatePeerFlags.Quiet, "quiet", "q", false, "Quietly returns if peer not found")

	// Required
	terminatePeerCmd.MarkFlagRequired("name")
}

func preRunTerminatePeer() {

}

func terminatePeer() {
	fmt.Println("Terminating Peer...", terminatePeerFlags)
	// Check if folder exists
	fullPath := path.Join(hlfdPath, peerDepFolder, terminatePeerFlags.Name)
	_, err := os.Stat(fullPath)
	// Return quietly if no such exists (if quiet mode is active)
	if err != nil && terminatePeerFlags.Quiet {
		return
	}
	cobra.CheckErr(err)

	// execute(fullPath, "docker-compose", "down", "-v", "--rmi", "local")
	cmds := []string{
		// Run docker-compose stop
		`cd ` + fullPath,
		`sudo docker-compose down -v --rmi local`,

		// Remove folder
		`cd ` + path.Join(hlfdPath, peerDepFolder),
		`sudo rm -rf ` + terminatePeerFlags.Name,
	}
	// execute(path.Join(hlfdPath, peerDepFolder), "sudo", "rm", "-rf", terminatePeerFlags.Name)

	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)
}
