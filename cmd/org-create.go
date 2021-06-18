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
	"hlfd/cmd/os_exec_utils"
	"os"
	"path"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Flags
var orgCreateFlags struct {
	MSP         string
	Name        string
	CaAddr      string
	CaAdminUser string
	CaAdminPass string
	// CaClientPath        string
	CaName        string
	CaTlsCertPath string
}

var (
	orgPath    = ""
	orgMSPPath = ""
)

type OrgInfo struct {
	Name   string `json:"name"`
	MspId  string `json:"mspId"`
	CaInfo CAInfo `json:"caInfo"`
	MspDir string `json:"mspDir"`
	// Policies interface{}
}

type CAInfo struct {
	CaName           string `json:"caName"`
	CaAddr           string `json:"caAddr"`
	CaClientHomePath string `json:"caClientHomePath"`
	CaTlsCertPath    string `json:"caTlsCertPath"`
}

// orgCreateCmd represents the ca command
var orgCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create HLF organization",
	Long:  `Create HLF organization`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		// container name greater than 2 chars..
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunRoot()
		preRunOrgCreate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		orgCreate()
	},
}

func init() {
	// Options to open firewall port for this
	orgCmd.AddCommand(orgCreateCmd)
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.Name, "name", "n", "", "Org name (required)")
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.MSP, "msp-id", "m", "", "MSP ID (required)")

	// If already got msp directory, use that
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.MSP, "msp-dir", "d", "", "MSP Dir ")
	// Else, get ca infos and create msp directory
	// TODO: If local ca, ca-name should be sufficient
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.CaName, "ca-name", "", ``, "Fabric Certificate Authority name to generate certs for org (required)")
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.CaAddr, "ca-addr", "", ``, "Fabric Certificate Authority address to generate certs for org (required)")
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.CaAdminUser, "ca-admin-user", "", ``, "Fabric Certificate Authority admin user to generate certs for org (required)")
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.CaAdminPass, "ca-admin-pass", "", ``, "Fabric Certificate Authority admin pass to generate certs for org (required)")
	// orgCreateCmd.Flags().StringVarP(&orgCreateFlags.CaClientPath, "ca-client-path", "", ``, "Path to fabric-ca-client binary")
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.CaTlsCertPath, "ca-tls-cert-path", "", ``, "Path to ca's pem encoded tls certificate (if applicable)")

	// Required
	orgCreateCmd.MarkFlagRequired("name")
	orgCreateCmd.MarkFlagRequired("msp-id")
}

func preRunOrgCreate() {
	// Validate
	// Force delete existing, if flag is set
	// if orgCreateFlags.Force {
	// 	terminatePeerFlags.Name = orgCreateFlags.Name
	// 	terminatePeerFlags.Quiet = true
	// 	terminatePeer()
	// }
	// 1. Create folders
	orgPath = path.Join(hlfdPath, orgCommonFolder, orgCreateFlags.Name)
	throwIfFileExists(orgPath)
	err := os.MkdirAll(orgPath, commonFilUmask)
	throwOtherThanFileExistError(err)
	orgMSPPath = path.Join(orgPath, mspFolder)
	err = os.MkdirAll(orgMSPPath, commonFilUmask)
	throwOtherThanFileExistError(err)

	// Set variables
	caClientHomePath = path.Join(hlfdPath, caClientHomeFolder, orgCreateFlags.CaName)

}

func orgCreate() {
	// 2. Generate Org msp folder (using ca)
	generateOrgMSP()

	// 3. Store org info in json

}

func generateOrgMSP() {
	// 1. Check ca client
	dldCaBinariesIfNotExist()
	// Enroll CA
	enrollCaAdminOrg()

	// Make Admin msp dirs
	cmds := []string{
		`cd ` + orgPath,
		`cd msp`,
		`mkdir -p cacerts`,
		`mkdir -p tlscacerts`,
		`mkdir -p intermediatecerts`,
		`mkdir -p tlsintermediatecerts`,
		`mkdir -p operationscerts`,
		`mkdir -p admincerts`,
	}
	_, err := os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)

	// Copy from CA admin and make org structure
	// Cacert
	cmds = []string{
		`cd ` + caClientHomePath,
		`cd msp`,
		`cp cacerts/*.pem ` + path.Join(orgPath, "msp", "cacerts", "ca.pem"),
	}

	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)

	// Tlscacert
	cmds = []string{
		`cd ` + caClientHomePath,
		`cd msp`,
		`cp cacerts/*.pem ` + path.Join(orgPath, "msp", "tlscacerts", "ca.pem"),
	}

	_, err = os_exec_utils.ExecMultiCommand(cmds)
	cobra.CheckErr(err)

	// Node ou
	generateNodeOUConfigOrg()

}

func enrollCaAdminOrg() {
	fmt.Println("Enrolling CA Admin...")
	// Make FABRIC_CA_CLIENT_HOME folder
	err := os.MkdirAll(caClientHomePath, commonFilUmask)
	throwOtherThanFileExistError(err)

	userEncodedCaUrl := getUserEncodedCaUrl(orgCreateFlags.CaAddr, orgCreateFlags.CaAdminUser, orgCreateFlags.CaAdminPass)
	enrollCmd := `./fabric-ca-client enroll -u ` + userEncodedCaUrl + ` --caname ` + orgCreateFlags.CaName
	// Tls vs no tls command
	if orgCreateFlags.CaTlsCertPath != "" {
		enrollCmd = enrollCmd + ` --tls.certfiles ` + orgCreateFlags.CaTlsCertPath
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

func generateNodeOUConfigOrg() {
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

	writeBytesToFile("config.yaml", orgMSPPath, yamlB)
}
