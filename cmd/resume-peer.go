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

var resumePeerFlags struct {
	Name string
}

//
var resumePeerCmd = &cobra.Command{
	Use:   "peer",
	Short: "Resumes Peer container.",
	Long:  `Resumes Hyperledger Fabric Peer container.`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunResumePeer()
	},
	Run: func(cmd *cobra.Command, args []string) {
		resumePeer()
	},
}

func init() {
	resumeCmd.AddCommand(resumePeerCmd)
	resumePeerCmd.Flags().StringVarP(&resumePeerFlags.Name, "name", "n", "", "Peer name")

	// Required
	resumePeerCmd.MarkFlagRequired("name")
}

func preRunResumePeer() {

}

func resumePeer() {
	fmt.Println("Resuming Peer...", resumePeerFlags)
	// Check if folder exists
	fullPath := path.Join(hlfdPath, peerDepFolder, resumePeerFlags.Name)
	_, err := os.Stat(fullPath)
	cobra.CheckErr(err)

	// Run docker-compose
	execute(fullPath, "sudo", "docker-compose", "up", "-d")
}
