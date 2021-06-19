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
var depPeerFlags struct {
	PeerName            string
	TLSEnabled          bool
	Port                int
	ExternalPort        int
	PeerHomeVolumeMount string
	DockerNetwork       string
	ImageTag            string
	// MSPId               string
	PeerLogging   string
	CorePeerAddr  string
	ChaincodeAddr string
	// CaAddr              string
	CaAdminUser string
	CaAdminPass string
	// CaClientPath        string
	CaName string
	// CaTlsCertPath string
	StateDB       string
	CouchAdmin    string
	CouchPass     string
	CouchImageTag string
	CouchPort     int
	OrgName       string

	ForceTerminate bool
}

// Deployment files path
var peerDepPath = ""
var dockerComposeFileNamePeer = ""
var peerPass = uuid.New().String()
var peerMspPath = ""
var peerTlsPath = ""
var peerDepOrgInfo OrgInfo

//
var deployPeerCmd = &cobra.Command{
	Use:   "peer",
	Short: "Deploys Peer.",
	Long:  `Deploys Hyperledger Fabric Peer node.`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		// container name greater than 2 chars..
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunRoot()
		preRunDepPeer()
	},
	Run: func(cmd *cobra.Command, args []string) {
		deployPeer()
	},
}

func init() {
	// Options to open firewall port for this
	deployCmd.AddCommand(deployPeerCmd)
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.PeerName, "name", "n", "", "Name of the Peer to deploy (required)")
	deployPeerCmd.Flags().BoolVarP(&depPeerFlags.TLSEnabled, "tls", "t", false, "Enable TLS")
	deployPeerCmd.Flags().IntVarP(&depPeerFlags.Port, "port", "p", 7051, "Peer port inside docker container")
	deployPeerCmd.Flags().IntVarP(&depPeerFlags.ExternalPort, "eport", "e", -1, "Peer port mapping to container's host")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.DockerNetwork, "docker-network", "d", "hlfd", "Docker network name")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.ImageTag, "image-tag", "i", "2.2", "Hyperledger Peer docker image tag")
	deployPeerCmd.Flags().BoolVarP(&depPeerFlags.ForceTerminate, "force", "f", false, "Force deploy or terminate peer if peer with given name already exists")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.PeerLogging, "peer-log", "l", "INFO", "Peer logging spec {INFO | DEBUG}")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.CorePeerAddr, "core-peer-addr", "a", ``, "Externally accessible address of peer / CORE_PEER_ADDRESS")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.ChaincodeAddr, "chaincode-addr", "c", ``, "Externally accessible address of chaincode / CORE_PEER_CHAINCODEADDRESS")
	// deployPeerCmd.Flags().StringVarP(&depPeerFlags.MSPId, "msp-id", "m", ``, "MSP ID of peer / CORE_PEER_MSPID (required)")

	// TODO: If local ca, ca-name should be sufficient
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.CaName, "ca-name", "", ``, "Fabric Certificate Authority name to generate certs for peer (required)")
	// deployPeerCmd.Flags().StringVarP(&depPeerFlags.CaAddr, "ca-addr", "", ``, "Fabric Certificate Authority address to generate certs for peer (required)")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.CaAdminUser, "ca-admin-user", "", ``, "Fabric Certificate Authority admin user to generate certs for peer (required)")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.CaAdminPass, "ca-admin-pass", "", ``, "Fabric Certificate Authority admin pass to generate certs for peer (required)")
	// deployPeerCmd.Flags().StringVarP(&depPeerFlags.CaClientPath, "ca-client-path", "", ``, "Path to fabric-ca-client binary")
	// deployPeerCmd.Flags().StringVarP(&depPeerFlags.CaTlsCertPath, "ca-tls-cert-path", "", ``, "Path to ca's pem encoded tls certificate (if applicable)")

	// DB
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.StateDB, "state-db", "s", `goleveldb`, "World state database { goleveldb | CouchDB }. Default is goleveldb")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.CouchAdmin, "couchdb-admin-user", "", `admin`, "Admin username for couch db")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.CouchPass, "couchdb-admin-pass", "", `adminpw`, "Admin pass for couch db")
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.CouchImageTag, "couchdb-image-tag", "", `3.1`, "Image tag for couch db")
	deployPeerCmd.Flags().IntVarP(&depPeerFlags.CouchPort, "couchdb-port", "", -1, "CouchDB port")

	//
	deployPeerCmd.Flags().StringVarP(&depPeerFlags.OrgName, "org-name", "", ``, "Name of the HLFD created organization to which this peer will be a part of")

	// Required
	deployPeerCmd.MarkFlagRequired("name")
	deployPeerCmd.MarkFlagRequired("ca-name")
	deployPeerCmd.MarkFlagRequired("ca-admin-user")
	deployPeerCmd.MarkFlagRequired("ca-admin-pass")
	deployPeerCmd.MarkFlagRequired("org-name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func preRunDepPeer() {
	// Check of Org exists, else throw err
	if depPeerFlags.OrgName == "" {
		err := fmt.Errorf("org-name cannot be empty")
		cobra.CheckErr(err)
	}
	// Load orginfo
	peerDepOrgInfo = loadOrgInfo(depPeerFlags.OrgName)
	// Fill in optional flags
	// if depPeerFlags.ContainerName == "" {
	// 	depPeerFlags.ContainerName = depPeerFlags.CaName
	// }
	if depPeerFlags.ExternalPort < 0 { // Check allowed ports as per standards
		depPeerFlags.ExternalPort = depPeerFlags.Port
	}

	if depPeerFlags.CorePeerAddr == "" {
		depPeerFlags.CorePeerAddr = depPeerFlags.PeerName + `:` + strconv.Itoa(depPeerFlags.Port)
	}
	if depPeerFlags.ChaincodeAddr == "" {
		depPeerFlags.ChaincodeAddr = depPeerFlags.PeerName + `:7052`
	}
	if depPeerFlags.CouchPort < 0 {
		depPeerFlags.CouchPort = depPeerFlags.Port + 100
	}

	// Force terminate existing, if flag is set
	if depPeerFlags.ForceTerminate {
		terminatePeerFlags.Name = depPeerFlags.PeerName
		terminatePeerFlags.Quiet = true
		terminatePeer()
	}

	// Create folders for storing deployment files
	peerDepPath = path.Join(hlfdPath, peerDepFolder, depPeerFlags.PeerName)
	// Check if already exists
	throwIfFileExists(peerDepPath)
	err := os.MkdirAll(peerDepPath, commonFilUmask)
	throwOtherThanFileExistError(err)

	if depPeerFlags.PeerHomeVolumeMount == "" {
		depPeerFlags.PeerHomeVolumeMount = path.Join(peerDepPath, peerHomeFolder)
	}

	// Create volume-mount path directories
	fullPath := depPeerFlags.PeerHomeVolumeMount
	err = os.MkdirAll(fullPath, commonFilUmask)
	throwOtherThanFileExistError(err)

	// Set variables
	dockerComposeFileNamePeer = "docker-compose.yaml"
	peerMspPath = path.Join(peerDepPath, mspFolder)
	peerTlsPath = path.Join(peerDepPath, tlsFolder)
	caClientHomePath = path.Join(hlfdPath, caClientHomeFolder, depPeerFlags.CaName)

	//
	if depPeerFlags.PeerLogging != "INFO" && depPeerFlags.PeerLogging != "DEBUG" {
		err := fmt.Errorf("invalid peer-log option: %v", depPeerFlags.PeerLogging)
		cobra.CheckErr(err)
	}

	// Select CA
	selectedCA = selectCaFromList(depOrdererFlags.CaName, peerDepOrgInfo.CaInfo)

}

func deployPeer() {
	fmt.Println("Deploying Peer...", depPeerFlags)
	// 1. Generate certs and put them in right folders
	generatePeerCredentials()
	// 2. Generate yaml & env
	yamlB := generatePeerYAMLBytes()
	writeBytesToFile(dockerComposeFileNamePeer, peerDepPath, yamlB)
	if depPeerFlags.StateDB == `CouchDB` {
		envB := generatePeerEnvBytes()
		writeBytesToFile(".env", peerDepPath, envB)
	}
	// 3. Up
	_, err := os_exec_utils.ExecMultiCommand([]string{
		`cd ` + peerDepPath,
		`docker-compose up -d`,
	})
	cobra.CheckErr(err)
}

func generatePeerYAMLBytes() (yamlB []byte) {
	yamlObj := Object{
		"version": "2",
		"networks": Object{
			depPeerFlags.DockerNetwork: Object{},
		},
		"services": Object{
			depPeerFlags.PeerName: Object{
				// "env_file": ".env",
				"image": "hyperledger/fabric-peer:" + depPeerFlags.ImageTag,
				"environment": []string{
					"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
					`CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=` + `project_` + depPeerFlags.DockerNetwork,
					`FABRIC_LOGGING_SPEC=INFO`, // INFO / DEBUG
					`CORE_PEER_TLS_ENABLED=` + strconv.FormatBool(depPeerFlags.TLSEnabled),
					`CORE_PEER_PROFILE_ENABLED=false`, // Go profiling tools, must only be used non-prod
					`CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/signcerts/server.crt`,
					`CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/keystore/server.key`,
					`CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/cacerts/ca.crt`,
					//
					`CORE_PEER_ID=` + depPeerFlags.PeerName,
					`CORE_PEER_ADDRESS=` + depPeerFlags.CorePeerAddr,                     // Externally accessible (only for inside org) / peer0.org1.medisotv2.com:7051
					`CORE_PEER_LISTENADDRESS=0.0.0.0:` + strconv.Itoa(depPeerFlags.Port), // 0.0.0.0:7051
					`CORE_PEER_CHAINCODEADDRESS=` + depPeerFlags.ChaincodeAddr,           // peer0.org1.medisotv2.com:7052
					`CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:` + getPort(depPeerFlags.ChaincodeAddr),
					`CORE_PEER_GOSSIP_BOOTSTRAP=` + depPeerFlags.CorePeerAddr, // peer0.org1.medisotv2.com:7051
					// If this isn't set, the peer will not be known to other organizations.
					`CORE_PEER_GOSSIP_EXTERNALENDPOINT=` + depPeerFlags.CorePeerAddr, // peer0.org1.medisotv2.com:7051 / for outside org
					`CORE_PEER_LOCALMSPID=` + peerDepOrgInfo.MspId,
					//
					// CORE_PEER_TLS_CLIENTAUTHREQUIRED // mutual tls
				},
				"ports": []string{
					strconv.FormatInt(int64(depPeerFlags.ExternalPort), 10) + ":" + strconv.FormatInt(int64(depPeerFlags.Port), 10),
				},
				"command": `peer node start`,
				"volumes": []string{
					`/var/run/:/host/var/run/`,
					peerMspPath + `:/etc/hyperledger/fabric/msp`,
					peerTlsPath + `:/etc/hyperledger/fabric/tls`,
					depPeerFlags.PeerHomeVolumeMount + `:/var/hyperledger/production`,
				},
				"container_name": depPeerFlags.PeerName,
				"networks": []string{
					depPeerFlags.DockerNetwork,
				},
				"working_dir": "/opt/gopath/src/github.com/hyperledger/fabric/peer",
			},
		},
	}

	// Add couch db container if chosen
	if depPeerFlags.StateDB == `CouchDB` {
		couchdbName := depPeerFlags.PeerName + "_couchdb"

		// Add CouchDB container
		yamlObj["services"].(Object)[couchdbName] = Object{
			`container_name`: couchdbName,
			`image`:          `couchdb:` + depPeerFlags.CouchImageTag,
			`environment`: []string{
				`COUCHDB_USER=$` + CouchAdminEnv,
				`COUCHDB_PASSWORD=$` + CouchAdminPassEnv,
			},
			`ports`: []string{
				strconv.FormatInt(int64(depPeerFlags.CouchPort), 10) + ":" + strconv.FormatInt(int64(couchDBConstPort), 10),
			},
			`networks`: []string{
				depPeerFlags.DockerNetwork,
			},
		}

		// Modify peer container
		yamlObj["services"].(Object)[depPeerFlags.PeerName].(Object)[`environment`] = append(
			yamlObj["services"].(Object)[depPeerFlags.PeerName].(Object)[`environment`].([]string),
			[]string{
				`CORE_LEDGER_STATE_STATEDATABASE=` + depPeerFlags.StateDB,
				`CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=` + couchdbName + `:` + strconv.Itoa(couchDBConstPort),
				`CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=$` + CouchAdminEnv,
				`CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=$` + CouchAdminPassEnv,
			}...,
		)

		yamlObj["services"].(Object)[depPeerFlags.PeerName].(Object)[`depends_on`] = []string{
			couchdbName,
		}
	}

	// Parse yaml
	yamlB, err := yaml.Marshal(&yamlObj)
	cobra.CheckErr(err)

	return
}

func generatePeerEnvBytes() (envB []byte) {
	env := CouchAdminEnv + `=` + depPeerFlags.CouchAdmin + `
` + CouchAdminPassEnv + `=` + depPeerFlags.CouchPass + `
`

	envB = []byte(env)

	return
}

func generateNodeOUConfigPeer() {
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

	writeBytesToFile("config.yaml", peerMspPath, yamlB)
}

func generatePeerCredentials() {
	fmt.Println("Generating peer credentials...")
	// 0. Download CA Client binary if not found
	// if depPeerFlags.CaClientPath == "" {
	dldCaBinariesIfNotExist()
	// }
	// 1. Enroll CA Admin
	enrollCaAdminPeer()
	// 2. Register peer
	registerPeer()
	// 3. Enroll peer
	enrollPeer()
	if depPeerFlags.TLSEnabled {
		enrollPeerTls()
	}
	// 4. Directory organization
}

func enrollCaAdminPeer() {
	fmt.Println("Enrolling CA Admin...")
	// Make FABRIC_CA_CLIENT_HOME folder
	err := os.MkdirAll(caClientHomePath, commonFilUmask)
	throwOtherThanFileExistError(err)

	userEncodedCaUrl := getUserEncodedCaUrl(selectedCA.CaAddr, depPeerFlags.CaAdminUser, depPeerFlags.CaAdminPass)
	enrollCmd := `./fabric-ca-client enroll -u ` + userEncodedCaUrl + ` --caname ` + depPeerFlags.CaName
	// Tls vs no tls command
	if selectedCA.CaTlsCertPath != "" {
		enrollCmd = enrollCmd + ` --tls.certfiles ` + selectedCA.CaTlsCertPath
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

func registerPeer() {
	fmt.Println("Registering Peer with CA...")
	enrollCmd := `./fabric-ca-client register --caname ` + depPeerFlags.CaName + ` --id.name ` + depPeerFlags.PeerName + ` --id.secret ` + peerPass + ` --id.type peer`
	// Tls vs no tls command
	if selectedCA.CaTlsCertPath != "" {
		enrollCmd = enrollCmd + ` --tls.certfiles ` + selectedCA.CaTlsCertPath
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

func enrollPeer() {
	fmt.Println("Enrolling peer...")
	userEncodedCaUrl := getUserEncodedCaUrl(selectedCA.CaAddr, depPeerFlags.PeerName, peerPass)
	enrollCmd := `./fabric-ca-client enroll -u ` + userEncodedCaUrl +
		` --caname ` + depPeerFlags.CaName +
		` -M ` + peerMspPath +
		` --csr.hosts ` + getUrl(depPeerFlags.CorePeerAddr)
	// Tls vs no tls command
	if selectedCA.CaTlsCertPath != "" {
		enrollCmd = enrollCmd + ` --tls.certfiles ` + selectedCA.CaTlsCertPath
	}

	commands := []string{
		`export FABRIC_CA_CLIENT_HOME=` + caClientHomePath, // Place to put peer msp files
		`set -x`,
		`cd ` + binPath,
		enrollCmd,
	}

	_, err := os_exec_utils.ExecMultiCommand(commands)
	cobra.CheckErr(err)

	// Rename files
	cmds := []string{
		`cd ` + peerMspPath,
		`mv cacerts/*.pem cacerts/ca.pem`,
	}

	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)

	// Create ou config
	generateNodeOUConfigPeer()
}

func enrollPeerTls() {
	fmt.Println("Enrolling peer tls...")
	userEncodedCaUrl := getUserEncodedCaUrl(selectedCA.CaAddr, depPeerFlags.PeerName, peerPass)
	enrollCmd := `./fabric-ca-client enroll -u ` + userEncodedCaUrl +
		` --caname ` + depPeerFlags.CaName +
		` -M ` + peerTlsPath +
		` --csr.hosts ` + getUrl(depPeerFlags.CorePeerAddr) +
		` --csr.hosts ` + `localhost`
	// Tls vs no tls command
	if selectedCA.CaTlsCertPath != "" {
		enrollCmd = enrollCmd + ` --tls.certfiles ` + selectedCA.CaTlsCertPath
	}

	commands := []string{
		`export FABRIC_CA_CLIENT_HOME=` + caClientHomePath, // Place to put peer msp files
		`set -x`,
		`cd ` + binPath,
		enrollCmd,
	}

	_, err := os_exec_utils.ExecMultiCommand(commands)
	cobra.CheckErr(err)

	// Rename files
	cmds := []string{
		`cd ` + peerTlsPath,
		`mv cacerts/*.pem cacerts/ca.crt`,
		`mv keystore/*_sk keystore/server.key`,
		`mv signcerts/*.pem signcerts/server.crt`,
	}

	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)
}
