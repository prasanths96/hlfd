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

	"github.com/spf13/cobra"
)

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
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
