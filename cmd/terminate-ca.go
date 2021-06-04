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
	"os"
	"path"

	"github.com/spf13/cobra"
)

var terminateCaFlags struct {
	Name string
}

// caCmd represents the ca command
var terminateCaCmd = &cobra.Command{
	Use:   "ca",
	Short: "Terminates CA container",
	Long:  `Terminates Hyperledger Fabric Certificate Authority (CA) container`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunTerminateCa()
	},
	Run: func(cmd *cobra.Command, args []string) {
		terminateCA()
	},
}

func init() {
	terminateCmd.AddCommand(terminateCaCmd)
	terminateCaCmd.Flags().StringVarP(&terminateCaFlags.Name, "name", "n", "", "CA name")

	// Required
	terminateCaCmd.MarkFlagRequired("name")
}

func preRunTerminateCa() {

}

func terminateCA() {
	fmt.Println("Terminating CA...", terminateCaFlags)
	// Check if Ca folder exists
	fullPath := path.Join(hlfdPath, caDepFolder, terminateCaFlags.Name)
	_, err := os.Stat(fullPath)
	cobra.CheckErr(err)

	// Run docker-compose stop
	execute(fullPath, "docker-compose", "down", "-v", "--rmi", "local")

	// Remove ca folder
	execute(path.Join(hlfdPath, caDepFolder), "sudo", "rm", "-rf", terminateCaFlags.Name)
}
