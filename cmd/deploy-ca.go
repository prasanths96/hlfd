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
	"os/exec"
	"path"
	"strconv"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
	DockerNetwork     string
	ImageTag          string
}

var caFlags CAFlags

// Deployment files path
var caDepPath = ""
var dockerComposeFileName = ""

// caCmd represents the ca command
var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "Deploys CA",
	Long:  `Deploys Hyperledfer Fabric Certificate Authority (CA)`,
	PreRun: func(cmd *cobra.Command, args []string) {
		preProcess()
	},
	Run: func(cmd *cobra.Command, args []string) {
		deployCA()
	},
}

func init() {
	deployCmd.AddCommand(caCmd)
	caCmd.Flags().StringVarP(&caFlags.CaName, "name", "n", "", "Name of the CA to deploy. This name will be used when registering and enrolling certs (required)")
	caCmd.Flags().BoolVarP(&caFlags.TLSEnabled, "tls", "t", false, "Enable TLS")
	caCmd.Flags().IntVarP(&caFlags.Port, "port", "p", 8054, "CA server port inside docker container")
	caCmd.Flags().IntVarP(&caFlags.ExternalPort, "eport", "e", -1, "CA server port mapping to container's host")
	caCmd.Flags().StringVarP(&caFlags.AdminUser, "admin", "a", "admin", "CA admin username")
	caCmd.Flags().StringVarP(&caFlags.AdminPass, "pass", "s", "adminpw", "CA admin password or secret")
	// caCmd.Flags().StringVarP(&caFlags.CAHomeVolumeMount, "volume-mount-path", "v", "", "Host system path to mount CA home directory")
	caCmd.Flags().StringVarP(&caFlags.ContainerName, "container-name", "c", "", "Docker container name")
	caCmd.Flags().StringVarP(&caFlags.DockerNetwork, "docker-network", "d", "hlfd", "Docker network name")
	caCmd.Flags().StringVarP(&caFlags.ImageTag, "image-tag", "i", "latest", "Hyperledger CA docker image tag")

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

func preProcess() {
	// Fill in optional flags
	if caFlags.ContainerName == "" {
		caFlags.ContainerName = caFlags.CaName
	}
	if caFlags.ExternalPort < 0 { // Check allowed ports as per standards
		caFlags.ExternalPort = caFlags.Port
	}

	// Create folders for storing CA deployment files
	caDepPath = path.Join(hlfdPath, caDepFolder, caFlags.ContainerName)
	err := os.MkdirAll(caDepPath, commonFilUmask)
	checkOtherThanFileExistsError(err)

	if caFlags.CAHomeVolumeMount == "" {
		caFlags.CAHomeVolumeMount = caDepPath + "/ca-home"
	}

	// Create volume-mount path directories
	fullPath := caFlags.CAHomeVolumeMount
	err = os.MkdirAll(fullPath, commonFilUmask)
	checkOtherThanFileExistsError(err)

	// Set variables
	dockerComposeFileName = "docker-compose.yaml"

}

func deployCA() {
	fmt.Println("Deploying CA...", caFlags)
	// Create yaml file
	yamlB := generateCAYAMLBytes()
	// Create necessary dir and store file
	writeBytesToFile(dockerComposeFileName, caDepPath, yamlB)
	// Create necessary env file
	envB := generateCAEnvBytes()
	writeBytesToFile(".env", caDepPath, envB)
	// Set necessary env
	// setEnv()
	// Run docker-compose up -d
	execute(caDepPath, "docker-compose", "up", "-d")
}

func generateCAYAMLBytes() (yamlB []byte) {
	yamlObj := Object{
		"version": "2",
		"networks": Object{
			caFlags.DockerNetwork: Object{},
		},
		"services": Object{
			caFlags.ContainerName: Object{
				// "env_file": ".env",
				"image": "hyperledger/fabric-ca:" + caFlags.ImageTag,
				"environment": []string{
					"FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server",
					"FABRIC_CA_SERVER_CA_NAME=" + caFlags.CaName,
					"FABRIC_CA_SERVER_TLS_ENABLED=" + strconv.FormatBool(caFlags.TLSEnabled),
					"FABRIC_CA_SERVER_PORT=" + strconv.FormatInt(int64(caFlags.Port), 10),
				},
				"ports": []string{
					strconv.FormatInt(int64(caFlags.ExternalPort), 10) + ":" + strconv.FormatInt(int64(caFlags.Port), 10),
				},
				"command": `sh -c 'fabric-ca-server start -b $` + CaAdminEnv + `:$` + CaAdminPassEnv + ` -d'`,
				"volumes": []string{
					caFlags.CAHomeVolumeMount + ":/etc/hyperledger/fabric-ca-server",
				},
				"container_name": caFlags.ContainerName,
				"networks": []string{
					caFlags.DockerNetwork,
				},
			},
		},
	}

	// Parse yaml
	yamlB, err := yaml.Marshal(&yamlObj)
	cobra.CheckErr(err)

	return
}

func generateCAEnvBytes() (envB []byte) {
	env := `
		` + CaAdminEnv + `=` + caFlags.AdminUser + `
		` + CaAdminPassEnv + `=` + caFlags.AdminPass + `
	`

	envB = []byte(env)

	return
}

func setEnv() {
	// CA Admin User
	cmd := `export`
	arg := CaAdminEnv + `=` + caFlags.AdminUser
	execute(cmd, arg)
	arg = CaAdminPassEnv + `=` + caFlags.AdminPass
	execute(cmd, arg)
}

func execute(dir string, comdS string, args ...string) {
	comd := exec.Command(comdS, args...)
	if dir != "" {
		comd.Dir = dir
	}
	// stdin, err := comd.StdinPipe()
	// go func() {
	// 	defer stdin.Close()
	// 	io.WriteString(stdin, "an old falcon")
	// }()
	// stdout, err := comd.StdoutPipe()
	// cobra.CheckErr(err)

	out, err := comd.CombinedOutput()
	fmt.Println(string(out))
	cobra.CheckErr(err)
}
