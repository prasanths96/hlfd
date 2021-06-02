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

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy HLF components",
	Long:  `Deploy HLF components including: ca, peer, orderer`,
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
	// PreRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("prerun")
	// },
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("deploy called")
	// },
	// PostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("postrun")
	// },
}

func init() {
	rootCmd.AddCommand(deployCmd)
	// deployCmd.Flags().

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
