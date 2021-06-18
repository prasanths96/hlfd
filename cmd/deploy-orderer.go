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
	"strconv"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"hlfd/cmd/os_exec_utils"
)

// Flags
var depOrdererFlags struct {
	OrdererName            string
	TLSEnabled             bool
	Port                   int
	ExternalPort           int
	OrdererHomeVolumeMount string
	DockerNetwork          string
	ImageTag               string
	MSPId                  string
	OrdererLogging         string
	CaAddr                 string
	CaAdminUser            string
	CaAdminPass            string
	// CaClientPath        string
	CaClientVersion string
	CaName          string
	CaTlsCertPath   string
	OrdererAddr     string

	//
	GenesisPath string
	JoinCluster bool

	ForceTerminate bool
}

// Deployment files path
var ordererDepPath = ""
var dockerComposeFileNameOrderer = ""
var ordererPass = uuid.New().String()
var ordererMspPath = ""
var ordererTlsPath = ""

//
var deployOrdererCmd = &cobra.Command{
	Use:   "orderer",
	Short: "Deploys Orderer.",
	Long:  `Deploys Hyperledger Fabric Orderer node.`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		// container name greater than 2 chars..
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunRoot()
		preRunDepOrderer()
	},
	Run: func(cmd *cobra.Command, args []string) {
		deployOrderer()
	},
}

func init() {
	// Options to open firewall port for this
	deployCmd.AddCommand(deployOrdererCmd)
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.OrdererName, "name", "n", "", "Name of the Orderer to deploy (required)")
	deployOrdererCmd.Flags().BoolVarP(&depOrdererFlags.TLSEnabled, "tls", "t", false, "Enable TLS")
	deployOrdererCmd.Flags().IntVarP(&depOrdererFlags.Port, "port", "p", 7051, "Orderer port inside docker container")
	deployOrdererCmd.Flags().IntVarP(&depOrdererFlags.ExternalPort, "eport", "e", -1, "Orderer port mapping to container's host")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.DockerNetwork, "docker-network", "d", "hlfd", "Docker network name")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.ImageTag, "image-tag", "i", "2.2", "Hyperledger Orderer docker image tag")
	deployOrdererCmd.Flags().BoolVarP(&depOrdererFlags.ForceTerminate, "force", "f", false, "Force deploy or terminate orderer if orderer with given name already exists")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.OrdererLogging, "orderer-log", "l", "INFO", "Orderer logging spec {INFO | DEBUG}")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.MSPId, "msp-id", "m", ``, "MSP ID of orderer / ORDERER_GENERAL_LOCALMSPID (required)")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.OrdererAddr, "orderer-addr", "a", ``, "Externally accessible address of orderer")

	// TODO: If local ca, ca-name should be sufficient
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.CaName, "ca-name", "", ``, "Fabric Certificate Authority name to generate certs for orderer (required)")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.CaAddr, "ca-addr", "", ``, "Fabric Certificate Authority address to generate certs for orderer (required)")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.CaAdminUser, "ca-admin-user", "", ``, "Fabric Certificate Authority admin user to generate certs for orderer (required)")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.CaAdminPass, "ca-admin-pass", "", ``, "Fabric Certificate Authority admin pass to generate certs for orderer (required)")
	// deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.CaClientPath, "ca-client-path", "", ``, "Path to fabric-ca-client binary")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.CaClientVersion, "ca-client-version", "", `1.5.0`, "Version of fabric-ca-client binary (same as CA docker image version). Default: 1.5.0")
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.CaTlsCertPath, "ca-tls-cert-path", "", ``, "Path to ca's pem encoded tls certificate (if applicable)")

	//Passing Genesis block, if not passed, generate new
	deployOrdererCmd.Flags().StringVarP(&depOrdererFlags.GenesisPath, "genesis-path", "", ``, "Path to orderer genesis block to bootstrap with")
	// If no genesis block, either create new ord-system-channel (cluster) or join existing cluster
	// deployOrdererCmd.Flags().BoolVarP(&depOrdererFlags.JoinCluster, "join-existing-cluster", "", false, "Path to orderer genesis block to bootstrap with")

	// Required
	deployOrdererCmd.MarkFlagRequired("name")
	deployOrdererCmd.MarkFlagRequired("msp-id")
	deployOrdererCmd.MarkFlagRequired("ca-name")
	deployOrdererCmd.MarkFlagRequired("ca-addr")
	deployOrdererCmd.MarkFlagRequired("ca-admin-user")
	deployOrdererCmd.MarkFlagRequired("ca-admin-pass")
	deployOrdererCmd.MarkFlagRequired("orderer-addr")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func preRunDepOrderer() {
	// Fill in optional flags
	// if depOrdererFlags.ContainerName == "" {
	// 	depOrdererFlags.ContainerName = depOrdererFlags.CaName
	// }
	if depOrdererFlags.ExternalPort < 0 { // Check allowed ports as per standards
		depOrdererFlags.ExternalPort = depOrdererFlags.Port
	}

	// Force terminate existing, if flag is set
	if depOrdererFlags.ForceTerminate {
		// terminateOrdererFlags.Name = depOrdererFlags.OrdererName
		// terminateOrdererFlags.Quiet = true
		// terminateOrderer()
	}

	// Create folders for storing deployment files
	ordererDepPath = path.Join(hlfdPath, ordererDepFolder, depOrdererFlags.OrdererName)
	// Check if already exists
	throwIfFileExists(ordererDepPath)
	err := os.MkdirAll(ordererDepPath, commonFilUmask)
	throwOtherThanFileExistError(err)

	if depOrdererFlags.OrdererHomeVolumeMount == "" {
		depOrdererFlags.OrdererHomeVolumeMount = path.Join(ordererDepPath, ordererHomeFolder)
	}

	// Create volume-mount path directories
	fullPath := depOrdererFlags.OrdererHomeVolumeMount
	err = os.MkdirAll(fullPath, commonFilUmask)
	throwOtherThanFileExistError(err)

	// Set variables
	dockerComposeFileNameOrderer = "docker-compose.yaml"
	ordererMspPath = path.Join(ordererDepPath, mspFolder)
	ordererTlsPath = path.Join(ordererDepPath, tlsFolder)
	caClientHomePath = path.Join(hlfdPath, caClientHomeFolder, depOrdererFlags.CaName)

	//
	if depOrdererFlags.OrdererLogging != "INFO" && depOrdererFlags.OrdererLogging != "DEBUG" {
		err := fmt.Errorf("invalid orderer-log option: %v", depOrdererFlags.OrdererLogging)
		cobra.CheckErr(err)
	}
}

func deployOrderer() {
	fmt.Println("Deploying Orderer...", depOrdererFlags)
	// 1. Generate certs and put them in right folders
	generateOrdererCredentials()
	// 2. Generate yaml & env
	yamlB := generateOrdererYAMLBytes()
	writeBytesToFile(dockerComposeFileNameOrderer, ordererDepPath, yamlB)
	// 3. Up
	_, err := os_exec_utils.ExecMultiCommand([]string{
		`cd ` + ordererDepPath,
		`docker-compose up -d`,
	})
	cobra.CheckErr(err)
}

func generateOrdererYAMLBytes() (yamlB []byte) {
	yamlObj := Object{
		"version": "2",
		"networks": Object{
			depOrdererFlags.DockerNetwork: Object{},
		},
		"services": Object{
			depOrdererFlags.OrdererName: Object{
				// "env_file": ".env",
				"image": "hyperledger/fabric-orderer:" + depOrdererFlags.ImageTag,
				"environment": []string{
					`FABRIC_LOGGING_SPEC=INFO`, // INFO / DEBUG
					`ORDERER_GENERAL_LISTENADDRESS=0.0.0.0`,
					`ORDERER_GENERAL_LISTENPORT=` + strconv.Itoa(depOrdererFlags.Port),
					// `ORDERER_GENERAL_GENESISMETHOD=file`,
					// ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
					`ORDERER_GENERAL_GENESISMETHOD=none`,
					`ORDERER_GENERAL_LOCALMSPID=` + depOrdererFlags.MSPId,
					`ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp`,
					`ORDERER_GENERAL_TLS_ENABLED=` + strconv.FormatBool(depOrdererFlags.TLSEnabled),
					`ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/keystore/server.key`,
					`ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/signcerts/server.crt`,
					`ORDERER_GENERAL_TLS_ROOTCAS=/var/hyperledger/orderer/tls/cacerts/ca.crt`,
					//
					`ORDERER_KAFKA_TOPIC_REPLICATIONFACTOR=1`,
					`ORDERER_KAFKA_VERBOSE=true`,
					//
					`ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/signcerts/server.crt`,
					`ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/keystore/server.key`,
					`ORDERER_GENERAL_CLUSTER_ROOTCAS=/var/hyperledger/orderer/tls/cacerts/ca.crt`,
				},
				"ports": []string{
					strconv.FormatInt(int64(depOrdererFlags.ExternalPort), 10) + ":" + strconv.FormatInt(int64(depOrdererFlags.Port), 10),
				},
				"command": `orderer`,
				"volumes": []string{
					// `../system-genesis-block/genesis.block:/var/hyperledger/orderer/orderer.genesis.block`,
					ordererMspPath + `:/var/hyperledger/orderer/msp`,
					ordererTlsPath + `:/var/hyperledger/orderer/tls`,
					depOrdererFlags.OrdererHomeVolumeMount + `:/var/hyperledger/production/orderer`,
				},
				"container_name": depOrdererFlags.OrdererName,
				"networks": []string{
					depOrdererFlags.DockerNetwork,
				},
				"working_dir": "/opt/gopath/src/github.com/hyperledger/fabric",
			},
		},
	}

	// Parse yaml
	yamlB, err := yaml.Marshal(&yamlObj)
	cobra.CheckErr(err)

	return
}

func generateNodeOUConfigOrderer() {
	yamlObj := Object{
		`NodeOUs`: Object{
			`Enable`: true,
			`ClientOUIdentifier`: Object{
				`Certificate`:                  `cacerts/ca.pem`,
				`OrganizationalUnitIdentifier`: `client`,
			},
			`PeerOUIdentifier`: Object{
				`Certificate`:                  `cacerts/ca.pem`,
				`OrganizationalUnitIdentifier`: `peer`,
			},
			`AdminOUIdentifier`: Object{
				`Certificate`:                  `cacerts/ca.pem`,
				`OrganizationalUnitIdentifier`: `admin`,
			},
			`OrdererOUIdentifier`: Object{
				`Certificate`:                  `cacerts/ca.pem`,
				`OrganizationalUnitIdentifier`: `orderer`,
			},
		},
	}

	// Parse yaml
	yamlB, err := yaml.Marshal(&yamlObj)
	cobra.CheckErr(err)

	writeBytesToFile("config.yaml", ordererMspPath, yamlB)
}

func generateOrdererCredentials() {
	fmt.Println("Generating orderer credentials...")
	// 0. Download CA Client binary if not found
	// if depOrdererFlags.CaClientPath == "" {
	dldCaBinariesIfNotExist()
	// }
	// 1. Enroll CA Admin
	enrollCaAdminOrderer()
	// 2. Register orderer
	registerOrderer()
	// 3. Enroll orderer
	enrollOrderer()
	if depOrdererFlags.TLSEnabled {
		enrollOrdererTls()
	}
	// 4. Directory organization
}

func enrollCaAdminOrderer() {
	fmt.Println("Enrolling CA Admin...")
	// Make FABRIC_CA_CLIENT_HOME folder
	err := os.MkdirAll(caClientHomePath, commonFilUmask)
	throwOtherThanFileExistError(err)

	userEncodedCaUrl := getUserEncodedCaUrl(depOrdererFlags.CaAddr, depOrdererFlags.CaAdminUser, depOrdererFlags.CaAdminPass)
	enrollCmd := `./fabric-ca-client enroll -u ` + userEncodedCaUrl + ` --caname ` + depOrdererFlags.CaName
	// Tls vs no tls command
	if depOrdererFlags.CaTlsCertPath != "" {
		enrollCmd = enrollCmd + ` --tls.certfiles ` + depOrdererFlags.CaTlsCertPath
	}

	commands := []string{
		`export FABRIC_CA_CLIENT_HOME=` + caClientHomePath,
		`set -x`,
		`cd ` + binPath,
		enrollCmd,
	}

	_, err = os_exec_utils.ExecMultiCommand(commands)
	cobra.CheckErr(err)
}

func registerOrderer() {
	fmt.Println("Registering Orderer with CA...")
	enrollCmd := `./fabric-ca-client register --caname ` + depOrdererFlags.CaName + ` --id.name ` + depOrdererFlags.OrdererName + ` --id.secret ` + ordererPass + ` --id.type orderer`
	// Tls vs no tls command
	if depOrdererFlags.CaTlsCertPath != "" {
		enrollCmd = enrollCmd + ` --tls.certfiles ` + depOrdererFlags.CaTlsCertPath
	}
	commands := []string{
		`export FABRIC_CA_CLIENT_HOME=` + caClientHomePath,
		`set -x`,
		`cd ` + binPath,
		enrollCmd,
	}

	_, err := os_exec_utils.ExecMultiCommand(commands)
	cobra.CheckErr(err)
}

func enrollOrderer() {
	fmt.Println("Enrolling Orderer...")
	userEncodedCaUrl := getUserEncodedCaUrl(depOrdererFlags.CaAddr, depOrdererFlags.OrdererName, ordererPass)
	enrollCmd := `./fabric-ca-client enroll -u ` + userEncodedCaUrl +
		` --caname ` + depOrdererFlags.CaName +
		` -M ` + ordererMspPath +
		` --csr.hosts ` + getUrl(depOrdererFlags.OrdererAddr)
	// Tls vs no tls command
	if depOrdererFlags.CaTlsCertPath != "" {
		enrollCmd = enrollCmd + ` --tls.certfiles ` + depOrdererFlags.CaTlsCertPath
	}

	commands := []string{
		`export FABRIC_CA_CLIENT_HOME=` + caClientHomePath, // Place to put orderer msp files
		`set -x`,
		`cd ` + binPath,
		enrollCmd,
	}

	_, err := os_exec_utils.ExecMultiCommand(commands)
	cobra.CheckErr(err)

	// Rename files
	cmds := []string{
		`cd ` + ordererMspPath,
		`mv cacerts/*.pem cacerts/ca.pem`,
	}

	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)

	// Create ou config
	generateNodeOUConfigOrderer()
}

func enrollOrdererTls() {
	fmt.Println("Enrolling orderer tls...")
	userEncodedCaUrl := getUserEncodedCaUrl(depOrdererFlags.CaAddr, depOrdererFlags.OrdererName, ordererPass)
	enrollCmd := `./fabric-ca-client enroll -u ` + userEncodedCaUrl +
		` --caname ` + depOrdererFlags.CaName +
		` -M ` + ordererTlsPath +
		` --csr.hosts ` + getUrl(depOrdererFlags.OrdererAddr) +
		` --csr.hosts ` + `localhost`
	// Tls vs no tls command
	if depOrdererFlags.CaTlsCertPath != "" {
		enrollCmd = enrollCmd + ` --tls.certfiles ` + depOrdererFlags.CaTlsCertPath
	}

	commands := []string{
		`export FABRIC_CA_CLIENT_HOME=` + caClientHomePath, // Place to put orderer msp files
		`set -x`,
		`cd ` + binPath,
		enrollCmd,
	}

	_, err := os_exec_utils.ExecMultiCommand(commands)
	cobra.CheckErr(err)

	// Rename files
	cmds := []string{
		`cd ` + ordererTlsPath,
		`mv cacerts/*.pem cacerts/ca.crt`,
		`mv keystore/*_sk keystore/server.key`,
		`mv signcerts/*.pem signcerts/server.crt`,
	}

	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)
}
