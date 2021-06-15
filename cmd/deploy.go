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
	// "fmt"

	"fmt"
	"hlfd/cmd/os_exec_utils"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var binPath = ""
var caClientPath = ""
var caClientHomePath = ""

var depCommonFlags struct {
	CaClientVersion string
}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys components of HLF.",
	Long: `Deploys Hyperledger Fabric components deployed in docker containers such as:
	CA
	Peer
	Orderer.`,
	// Args: func(cmd *cobra.Command, args []string) error {
	// 	// err := cobra.NoArgs(cmd, args)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	// ArbitraryArgs
	// 	// OnlyValidArgs
	// 	// MinimumNArgs
	// 	// MaximumNArgs
	// 	// ExactArgs
	// 	// ExactValidArgs
	// 	// RangeArgs
	// 	fmt.Println("args validated", len(args))
	// 	return nil
	// },
	// PersistentPreRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("DEPLOY: PersistentPreRun")
	// 	preRunDeploy()
	// },
	Run: func(cmd *cobra.Command, args []string) {
		err := fmt.Errorf("deploy command cannot be used separately")
		cobra.CheckErr(err)
	},
	// PostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("postrun")
	// },
}

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.Flags().StringVarP(&depCommonFlags.CaClientVersion, "ca-client-version", "", `1.5.0`, "Version of fabric-ca-client binary (same as CA docker image version). Default: 1.5.0")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func preRunDeploy() {
	// For peer & orderer deployment
	binPath = path.Join(hlfdPath, binFolder)
	caClientPath = path.Join(binPath, caClientName)
}

func dldCaBinariesIfNotExist() {
	// https://github.com/hyperledger/fabric/releases/download/v2.2.3/hyperledger-fabric-linux-amd64-2.2.3.tar.gz
	// https://github.com/hyperledger/fabric-ca/releases/download/v1.5.0/hyperledger-fabric-ca-linux-amd64-1.5.0.tar.gz
	// Check and download
	exists := isFileExists(caClientPath)
	if exists {
		return
	}

	fmt.Println("Downloading fabric-ca-client binary...")

	// Make folders
	fmt.Println("Creating, ", binPath)
	err := os.MkdirAll(binPath, commonFilUmask)
	throwOtherThanFileExistError(err)

	// Download binaries
	caBinDldFileName := `hyperledger-fabric-ca-linux-amd64-` + depCommonFlags.CaClientVersion + `.tar.gz`
	caBinDldurl := `https://github.com/hyperledger/fabric-ca/releases/download/v` + depCommonFlags.CaClientVersion + `/` + caBinDldFileName
	execute(binPath, "wget", caBinDldurl)

	// Extract
	execute(binPath, "tar", "xvf", caBinDldFileName)

	// Move files
	execute(binPath, "mv", "bin/"+caClientName, ".")

	// Delete files
	delCmd := []string{
		"cd " + binPath,
		`rm -rf bin`,
		`rm -rf ` + caBinDldFileName,
		`rm -rf ` + caBinDldFileName + `*`,
	}
	_, err = os_exec_utils.ExecMultiCommand(delCmd)
	cobra.CheckErr(err)
}

// Utils
func getUserEncodedCaUrl(caAddr string, username string, pass string) (fullUrl string) {
	// parts[0] = https: , parts[1] = localhost:27054
	parts := strings.SplitN(caAddr, `//`, 2)
	if len(parts) < 2 {
		err := fmt.Errorf(`invalid ca address format. Correct format sample: "https://localhost:27054"`)
		cobra.CheckErr(err)
	}

	fullUrl = parts[0] + `//` + username + `:` + pass + `@` + parts[1]

	return
}

func getPort(fullAddr string) (port string) {
	parts := strings.Split(fullAddr, ":")
	port = parts[1] // Expecting correct address <address>:<port>
	return
}

func getUrl(fullAddr string) (url string) {
	parts := strings.Split(fullAddr, ":")
	url = parts[0] // Expecting correct address <address>:<port>
	return
}
