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
	"hlfd/cmd/os_exec_utils"
	"os"
	"path"
	"strings"

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
	//Policies
	ReaderPolicyS string
	WriterPolicyS string
	AdminPolicyS  string
}

var (
	orgPath    = ""
	orgMSPPath = ""

	loadedOrgInfo OrgInfo
)

type OrgInfo struct {
	Name     string      `json:"name"`
	MspId    string      `json:"mspId"`
	CaInfo   []OrgCAInfo `json:"caInfo"`
	MspDir   string      `json:"mspDir"`
	Policies Object      `json:"policies"` // map[Readers]Policy , map[Writers]Policy, map[Admins]Policy ...
}

type OrgCAInfo struct {
	CaName           string `json:"caName"`
	CaAddr           string `json:"caAddr"`
	CaClientHomePath string `json:"caClientHomePath"`
	CaTlsCertPath    string `json:"caTlsCertPath"`
}

type Policy struct {
	Type string
	Rule string
}

// orgCreateCmd represents the ca command
var orgCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create HLF organization",
	Long:  `Create HLF organization`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
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

	// Policies
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.ReaderPolicyS, "reader-policy", "", ``, `Org Reader Policy. Format: "<type> <policy>" Eg: "Signature OR('OrdererMSP.member')"`)
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.WriterPolicyS, "writer-policy", "", ``, `Org Writer Policy. Format: "<type> <policy>" Eg: "Signature OR('OrdererMSP.member')"`)
	orgCreateCmd.Flags().StringVarP(&orgCreateFlags.AdminPolicyS, "admin-policy", "", ``, `Org Admin Policy. Format: "<type> <policy>" Eg: "Signature OR('OrdererMSP.admin')"`)

	// Required
	orgCreateCmd.MarkFlagRequired("name")
	orgCreateCmd.MarkFlagRequired("msp-id")
	orgCreateCmd.MarkFlagRequired("ca-name")
	orgCreateCmd.MarkFlagRequired("ca-admin-user")
	orgCreateCmd.MarkFlagRequired("ca-admin-pass")
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
	orgImportPath := path.Join(hlfdPath, importCommonFolder, orgCommonFolder, orgCreateFlags.Name)
	throwIfFileExists(orgImportPath)

	err := os.MkdirAll(orgPath, commonFilUmask)
	throwOtherThanFileExistError(err)
	orgMSPPath = path.Join(orgPath, mspFolder)
	err = os.MkdirAll(orgMSPPath, commonFilUmask)
	throwOtherThanFileExistError(err)

	// Set variables
	caClientHomePath = path.Join(hlfdPath, caClientHomeFolder, orgCreateFlags.CaName)

	// Load CA
	if orgCreateFlags.CaName == "" {
		err := fmt.Errorf("ca-name cannot be empty")
		cobra.CheckErr(err)
	}
	if orgCreateFlags.CaAddr == "" {
		ca := loadCAInfo(orgCreateFlags.CaName)
		orgCreateFlags.CaAddr = getCaAddrFromCAInfo(ca) 
		if ca.TLSEnabled {
			orgCreateFlags.CaTlsCertPath = ca.TlsCertPath
		}
	}

}

func orgCreate() {
	// 2. Generate Org msp folder (using ca)
	generateOrgMSP()

	// 3. Store org info in json
	storeOrgInfo()
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

func storeOrgInfo() {
	// Parse policies
	readerPolicy := Policy{
		Type: "Signature",
		Rule: `OR('` + orgCreateFlags.MSP + `.member')`,
	}
	writerPolicy := Policy{
		Type: "Signature",
		Rule: `OR('` + orgCreateFlags.MSP + `.member')`,
	}
	adminPolicy := Policy{
		Type: "Signature",
		Rule: `OR('` + orgCreateFlags.MSP + `.admin')`,
	}

	if orgCreateFlags.ReaderPolicyS != "" {
		typee, rule := parsePolicyInput(orgCreateFlags.ReaderPolicyS)
		readerPolicy.Type = typee
		readerPolicy.Rule = rule
	}
	if orgCreateFlags.WriterPolicyS != "" {
		typee, rule := parsePolicyInput(orgCreateFlags.WriterPolicyS)
		writerPolicy.Type = typee
		writerPolicy.Rule = rule
	}
	if orgCreateFlags.AdminPolicyS != "" {
		typee, rule := parsePolicyInput(orgCreateFlags.AdminPolicyS)
		adminPolicy.Type = typee
		adminPolicy.Rule = rule
	}

	orgInfo := OrgInfo{
		Name:  orgCreateFlags.Name,
		MspId: orgCreateFlags.MSP,
		CaInfo: []OrgCAInfo{
			{
				CaName:           orgCreateFlags.CaName,
				CaAddr:           orgCreateFlags.CaAddr,
				CaClientHomePath: caClientHomePath,
				CaTlsCertPath:    orgCreateFlags.CaTlsCertPath,
			},
		},
		MspDir: path.Join(orgPath, "msp"),
		Policies: map[string]interface{}{
			"Readers": readerPolicy,
			"Writers": writerPolicy,
			"Admins":  adminPolicy,
		},
	}

	m, err := json.MarshalIndent(orgInfo, "", "    ")
	cobra.CheckErr(err)

	writeBytesToFile(orgInfoFileName, orgPath, m)
}

func parsePolicyInput(policyS string) (typee, rule string) {
	parts := strings.Split(policyS, " ")
	if len(parts) < 2 {
		err := fmt.Errorf("bad policy syntax")
		cobra.CheckErr(err)
	}
	typee = parts[0]
	rule = strings.Join(parts[1:], " ")
	return
}
