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

var installPrereqsFlags struct {
}

// caCmd represents the ca command
var installPrereqsCmd = &cobra.Command{
	Use:   "prereqs",
	Short: "Installs pre-requirements",
	Long:  `Installs pre-requirements needed for deploying Hyperledger Fabric components`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		preRunInstallPrereqs()
	},
	Run: func(cmd *cobra.Command, args []string) {
		installPrereqs()
	},
}

func init() {
	installCmd.AddCommand(installPrereqsCmd)
	// installPrereqsCmd.Flags().StringVarP(&installPrereqs.Name, "name", "n", "", "CA name")

	// Required
	// installPrereqsCmd.MarkFlagRequired("name")
}

func preRunInstallPrereqs() {

}

func installPrereqs() {
	// Install indiviual things

	// Install all
	execute("", "sudo", "apt", "update", "-y") // Don't forget to un-comment
	installGit()
	installCurl()
	installBuildEssential()
	installDocker()
	installDockerCompose()
	installGo()
}

// Different pre-reqs
func installGit() {
	fmt.Println("Checking git...")
	if isCmdExists("git") {
		return
	}

	fmt.Println("Installing git...")
	execute("", "sudo", "apt", "install", "git", "-y")
	execute("", "git", "--version")
}

func installCurl() {
	fmt.Println("Checking curl...")
	if isCmdExists("curl") {
		return
	}

	fmt.Println("Installing curl...")
	execute("", "sudo", "apt", "install", "curl", "-y")
	execute("", "curl", "--version")
}

func installBuildEssential() {
	fmt.Println("Installing build-essential...")
	execute("", "sudo", "apt", "install", "build-essential", "-y")
}

func installDocker() {
	fmt.Println("Checking docker...")
	if isCmdExists("docker") {
		return
	}

	fmt.Println("Installing docker...")
	execute("", "sudo", "apt", "install", "apt-transport-https", "ca-certificates", "curl", "software-properties-common", "-y")
	fmt.Println("Running CURL.............................................")
	execute("", "sudo", "curl", "-fsSL", "https://download.docker.com/linux/ubuntu/gpg", "|", "sudo", "apt-key", "add", "-")
	fmt.Println("Running ADD_APT_REPOSITORY.............................................")
	execute("", "sudo", "add-apt-repository", "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable")
	fmt.Println("Running APT UPDATE.............................................")
	execute("", "sudo", "apt", "update", "-y")
	fmt.Println("Running APT INSTALL.............................................")
	execute("", "sudo", "apt", "install", "docker-ce", "-y")
	fmt.Println("Running APT USERMOD.............................................")
	execute("", "sudo", "usermod", "-aG", "docker ${USER}")

	fmt.Println("Running DOCKER VERSION.............................................")
	execute("", "docker", "--version")
}

func installDockerCompose() {
	fmt.Println("Checking docker-compose...")
	if isCmdExists("docker-compose") {
		return
	}

	fmt.Println("Installing docker-compose...")
	execute("", "curl", "-L", "https://github.com/docker/compose/releases/download/"+dockerComposeVersion+"/docker-compose-`uname -s`-`uname -m`", "-o", "/usr/local/bin/docker-compose")
	execute("", "sudo", "chmod", "+x", "/usr/local/bin/docker-compose")

	execute("", "docker-compose", "--version")
}

func installGo() {
	fmt.Println("Checking go...")
	if isCmdExists("go") {
		return
	}

	fmt.Println("Installing go...")
	execute("~/", "wget", "https://dl.google.com/go/go"+goVersion+".linux-amd64.tar.gz")
	execute("~/", "tar", "-xvf", "go"+goVersion+".linux-amd64.tar.gz")
	execute("~/", "rm", "/usr/local/go", "-rf")
	execute("~/", "mv", "go", "/usr/local")
	execute("~/", "echo", "'export GOROOT=/usr/local/go'>>~/.profile")
	execute("~/", "echo", "'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH'>>~/.profile")
	execute("~/", "source", "~/.profile")

	execute("", "go", "version")
}
