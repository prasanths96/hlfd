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

	"github.com/spf13/cobra"
)

// Flags

type CAFlags struct {
	CaName            string
	TLSEnabled        bool
	Port              int
	ExternalPort      int
	AdminUser         string
	AdminPass         string
	CAHomeVolumeMount string
	ContainerName     string
}

var caFlags CAFlags

// caCmd represents the ca command
var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "Deploys CA",
	Long:  `Deploys Hyperledfer Fabric Certificate Authority (CA)`,
	Run: func(cmd *cobra.Command, args []string) {
		deployCA()
	},
}

func init() {
	deployCmd.AddCommand(caCmd)
	caCmd.Flags().StringVarP(&caFlags.CaName, "name", "n", "", "Name of the CA to deploy. This name will be used when registering and enrolling certs (required)")
	caCmd.Flags().BoolVarP(&caFlags.TLSEnabled, "tls", "t", false, "Enable TLS")
	caCmd.Flags().IntVarP(&caFlags.Port, "port", "p", 8054, "CA server port inside docker container")
	caCmd.Flags().IntVarP(&caFlags.ExternalPort, "eport", "e", 8054, "CA server port mapping to container's host")
	caCmd.Flags().StringVarP(&caFlags.AdminUser, "admin", "a", "admin", "CA admin username")
	caCmd.Flags().StringVarP(&caFlags.AdminPass, "pass", "s", "adminpw", "CA admin password or secret")
	caCmd.Flags().StringVarP(&caFlags.CAHomeVolumeMount, "volume-mount-path", "v", "", "Host system path to mount CA home directory")
	caCmd.Flags().StringVarP(&caFlags.ContainerName, "container-name", "c", "", "Docker container name")

	// Required
	caCmd.MarkFlagRequired("name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deployCA() {
	fmt.Println("Deploying CA...", caFlags)
}
