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

var stopOrdererFlags struct {
	Name string
}

//
var stopOrdererCmd = &cobra.Command{
	Use:   "orderer",
	Short: "Stops Orderer container.",
	Long:  `Stops Hyperledger Fabric Orderer container.`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunStopOrderer()
	},
	Run: func(cmd *cobra.Command, args []string) {
		stopOrderer()
	},
}

func init() {
	stopCmd.AddCommand(stopOrdererCmd)
	stopOrdererCmd.Flags().StringVarP(&stopOrdererFlags.Name, "name", "n", "", "Orderer name")

	// Required
	stopOrdererCmd.MarkFlagRequired("name")
}

func preRunStopOrderer() {

}

func stopOrderer() {
	fmt.Println("Stopping Orderer...", stopOrdererFlags)
	// Check if folder exists
	fullPath := path.Join(hlfdPath, ordererDepFolder, stopOrdererFlags.Name)
	_, err := os.Stat(fullPath)
	cobra.CheckErr(err)

	// Run docker-compose stop
	execute(fullPath, "sudo", "docker-compose", "stop")
}
