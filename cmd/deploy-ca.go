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
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Flags
var depCaFlags struct {
	CaName            string
	TLSEnabled        bool
	Port              int
	ExternalPort      int
	AdminUser         string
	AdminPass         string
	CAHomeVolumeMount string
	DockerNetwork     string
	ImageTag          string
	// ContainerName     string
	ForceTerminate bool
}

type CAInfo struct {
	CaName      string `json:"caName"`
	CaHost      string `json:"caHost"`
	CaPort      int    `json:"caPort"`
	TLSEnabled  bool   `json:"tlsEnabled"`
	TlsCertPath string `json:"tlsCertPath"`
}

// Deployment files path
var caDepPath = ""
var dockerComposeFileNameCa = ""

// deployCaCmd represents the ca command
var deployCaCmd = &cobra.Command{
	Use:   "ca",
	Short: "Deploys CA.",
	Long:  `Deploys Hyperledger Fabric Certificate Authority (CA).`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		// container name greater than 2 chars..
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunRoot()
		preRunDepCa()
	},
	Run: func(cmd *cobra.Command, args []string) {
		deployCA()
	},
}

func init() {
	// Options to open firewall port for this
	deployCmd.AddCommand(deployCaCmd)
	deployCaCmd.Flags().StringVarP(&depCaFlags.CaName, "name", "n", "", "Name of the CA to deploy. This name will be used when registering and enrolling certs (required)")
	deployCaCmd.Flags().BoolVarP(&depCaFlags.TLSEnabled, "tls", "t", false, "Enable TLS")
	deployCaCmd.Flags().IntVarP(&depCaFlags.Port, "port", "p", 8054, "CA server port inside docker container")
	deployCaCmd.Flags().IntVarP(&depCaFlags.ExternalPort, "eport", "e", -1, "CA server port mapping to container's host")
	deployCaCmd.Flags().StringVarP(&depCaFlags.AdminUser, "admin", "a", "admin", "CA admin username")
	deployCaCmd.Flags().StringVarP(&depCaFlags.AdminPass, "pass", "s", "adminpw", "CA admin password or secret")
	// deployCaCmd.Flags().StringVarP(&depCaFlags.CAHomeVolumeMount, "volume-mount-path", "v", "", "Host system path to mount CA home directory")
	// deployCaCmd.Flags().StringVarP(&depCaFlags.ContainerName, "container-name", "c", "", "Docker container name")
	deployCaCmd.Flags().StringVarP(&depCaFlags.DockerNetwork, "docker-network", "d", "hlfd", "Docker network name")
	deployCaCmd.Flags().StringVarP(&depCaFlags.ImageTag, "image-tag", "i", "latest", "Hyperledger CA docker image tag")
	deployCaCmd.Flags().BoolVarP(&depCaFlags.ForceTerminate, "force", "f", false, "Force deploy or terminate ca if ca with given name already exists")

	// Required
	deployCaCmd.MarkFlagRequired("name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func preRunDepCa() {
	// Fill in optional flags
	// if depCaFlags.ContainerName == "" {
	// 	depCaFlags.ContainerName = depCaFlags.CaName
	// }
	if depCaFlags.ExternalPort < 0 { // Check allowed ports as per standards
		depCaFlags.ExternalPort = depCaFlags.Port
	}

	// Force terminate existing ca, if flag is set
	if depCaFlags.ForceTerminate {
		quietTerminateCa(depCaFlags.CaName)
	}

	// Create folders for storing CA deployment files
	caDepPath = path.Join(hlfdPath, caDepFolder, depCaFlags.CaName)
	// Check if CA already exists
	throwIfFileExists(caDepPath)
	err := os.MkdirAll(caDepPath, commonFilUmask)
	throwOtherThanFileExistError(err)

	if depCaFlags.CAHomeVolumeMount == "" {
		depCaFlags.CAHomeVolumeMount = path.Join(caDepPath, caHomeFolder)
	}

	// Create volume-mount path directories
	fullPath := depCaFlags.CAHomeVolumeMount
	err = os.MkdirAll(fullPath, commonFilUmask)
	throwOtherThanFileExistError(err)

	// Set variables
	dockerComposeFileNameCa = "docker-compose.yaml"

}

func deployCA() {
	fmt.Println("Deploying CA...", depCaFlags)
	// Create yaml file
	yamlB := generateCAYAMLBytes()
	// Create necessary dir and store file
	writeBytesToFile(dockerComposeFileNameCa, caDepPath, yamlB)
	// Create necessary env file
	envB := generateCAEnvBytes()
	writeBytesToFile(".env", caDepPath, envB)
	// Set necessary env
	// setEnv()
	// Run docker-compose up -d
	execute(caDepPath, "docker-compose", "up", "-d")

	// Store info
	storeCAInfo()
}

func generateCAYAMLBytes() (yamlB []byte) {
	yamlObj := Object{
		"version": "2",
		"networks": Object{
			depCaFlags.DockerNetwork: Object{},
		},
		"services": Object{
			depCaFlags.CaName: Object{
				// "env_file": ".env",
				"image": "hyperledger/fabric-ca:" + depCaFlags.ImageTag,
				"environment": []string{
					"FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server",
					"FABRIC_CA_SERVER_CA_NAME=" + depCaFlags.CaName,
					"FABRIC_CA_SERVER_TLS_ENABLED=" + strconv.FormatBool(depCaFlags.TLSEnabled),
					"FABRIC_CA_SERVER_PORT=" + strconv.FormatInt(int64(depCaFlags.Port), 10),
					"FABRIC_CA_SERVER_CSR_HOSTS=" + GetOutboundIP(),
				},
				"ports": []string{
					strconv.FormatInt(int64(depCaFlags.ExternalPort), 10) + ":" + strconv.FormatInt(int64(depCaFlags.Port), 10),
				},
				"command": `sh -c 'fabric-ca-server start -b $` + CaAdminEnv + `:$` + CaAdminPassEnv + ` -d'`,
				"volumes": []string{
					depCaFlags.CAHomeVolumeMount + ":/etc/hyperledger/fabric-ca-server",
				},
				"container_name": depCaFlags.CaName,
				"networks": []string{
					depCaFlags.DockerNetwork,
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
	env := CaAdminEnv + `=` + depCaFlags.AdminUser + `
` + CaAdminPassEnv + `=` + depCaFlags.AdminPass + `
`

	envB = []byte(env)

	return
}

func quietTerminateCa(peerName string) {
	terminateCaFlags.Name = peerName
	terminateCaFlags.Quiet = true
	terminateCA()
}

func storeCAInfo() {
	caInfo := CAInfo{
		CaName:      depCaFlags.CaName,
		CaHost:      GetOutboundIP(),
		CaPort:      depCaFlags.ExternalPort,
		TLSEnabled:  depCaFlags.TLSEnabled,
		TlsCertPath: path.Join(depCaFlags.CAHomeVolumeMount, caTlsCertFileName),
	}

	m, err := json.MarshalIndent(caInfo, "", "    ")
	cobra.CheckErr(err)

	writeBytesToFile(caInfoFileName, caDepPath, m)
}

func loadCAInfo(caName string) (loadedCaInfo CAInfo) {
	// See if available locally
	caPath := path.Join(hlfdPath, caDepFolder, caName)
	_, err := os.Stat(caPath)
	if err != nil {
		if !os.IsNotExist(err) {
			cobra.CheckErr(err)
		}

		// Check if in imports
		caPath = path.Join(hlfdPath, importCommonFolder, caDepFolder, caName)
		throwIfFileNotExist(caPath)
	}

	caInfoPath := path.Join(caPath, caInfoFileName)
	err = json.Unmarshal(readFileBytes(caInfoPath), &loadedCaInfo)
	cobra.CheckErr(err)

	// Resolve paths to absolute paths
	if loadedCaInfo.TLSEnabled {

		if loadedCaInfo.TlsCertPath[0] == '.' {
			loadedCaInfo.TlsCertPath = strings.Replace(loadedCaInfo.TlsCertPath, ".", caPath, 1)
		}
	}

	return
}

func getCaAddrFromCAInfo(caInfo CAInfo) (addr string) {
	addr = `http://`
	if caInfo.TLSEnabled {
		addr = `https://`
	}
	host := caInfo.CaHost

	// if GetOutboundIP() == host {
	// 	host = `localhost`
	// }

	addr = addr + host + `:` + strconv.Itoa(caInfo.CaPort)
	return
}
