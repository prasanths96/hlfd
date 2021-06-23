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
	"path"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var installPrereqsFlags struct {
}

// caCmd represents the ca command
var installPrereqsCmd = &cobra.Command{
	Use:   "prereqs",
	Short: "Installs pre-requirements.",
	Long:  `Installs pre-requirements needed for deploying Hyperledger Fabric components.`,
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
	var wg sync.WaitGroup
	execute("", "sudo", "apt", "update", "-y") // Don't forget to un-comment

	wg.Add(7)
	installGit(&wg)
	installWget(&wg)
	installCurl(&wg)
	installBuildEssential(&wg)
	installDocker(&wg)
	installDockerCompose(&wg)
	// installGo(&wg) // why is it installing everytime
	wg.Wait()

	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("Installation success!")
	fmt.Println(`Run "source ~/.profile" or re-login to access the following commands, if not previously installed:`)
	fmt.Println("\t- go")
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
}

// Different pre-reqs
func installGit(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Checking git...")
	if isCmdExists("git") {
		return
	}

	fmt.Println("Installing git...")
	execute("", "sudo", "apt", "install", "git", "-y")
	execute("", "git", "--version")
}

func installWget(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Checking wget...")
	if isCmdExists("wget") {
		return
	}

	fmt.Println("Installing wget...")
	execute("", "sudo", "apt", "install", "wget", "-y")
	execute("", "wget", "--version")
}

func installCurl(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Checking curl...")
	if isCmdExists("curl") {
		return
	}

	fmt.Println("Installing curl...")
	execute("", "sudo", "apt", "install", "curl", "-y")
	execute("", "curl", "--version")
}

func installBuildEssential(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Installing build-essential...")
	execute("", "sudo", "apt", "install", "build-essential", "-y")
}

func installDocker(wg *sync.WaitGroup) {
	defer wg.Done()
	// Notes: Problems with WSL: systemd not present in wsl so systemctl wont work. Instead, init is present.
	// use "service docker start", "service docker status"
	// So check if systemd is present
	fmt.Println("Checking docker...")
	if isCmdExists("docker") {
		return
	}

	fmt.Println("Installing docker...")
	execute("", "sudo", "apt", "install", "apt-transport-https", "ca-certificates", "curl", "software-properties-common", "-y")
	// execute("", "sudo", "curl", "-fsSL", "https://download.docker.com/linux/ubuntu/gpg", "|", "sudo", "apt-key", "add", "-")
	// execute("", "sudo", "curl", "-fsSL", "https://download.docker.com/linux/ubuntu/gpg>>/etc/apt/trusted.gpg")
	// execute("", "curl", "-fsSL", "https://download.docker.com/linux/ubuntu/gpg>>~/dockergpg")
	gpgBytes := execAndGetOutput("", "curl", "-fsSL", "https://download.docker.com/linux/ubuntu/gpg")
	writeBytesToFile("dockergpg", hlfdPath, gpgBytes)
	execute("", "sudo", "apt-key", "add", path.Join(hlfdPath, "dockergpg"))
	execute("", "sudo", "add-apt-repository", "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable")
	execute("", "sudo", "apt", "update", "-y")
	execute("", "sudo", "apt", "install", "docker-ce", "-y")
	username := strings.Trim(string(execAndGetOutput("", "whoami")), "\n")
	executeIgnoreErr("", "sudo", "groupadd", "docker")
	execute("", "sudo", "usermod", "-aG", "docker", username)
	execute("", "newgrp", "docker")

	execute("", "docker", "--version")
}

func installDockerCompose(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Checking docker-compose...")
	if isCmdExists("docker-compose") {
		return
	}

	fmt.Println("Installing docker-compose...")
	kernelname := strings.Trim(string(execAndGetOutput("", "uname", "-s")), "\n")
	machineHardwareName := strings.Trim(string(execAndGetOutput("", "uname", "-m")), "\n")
	execute("", "sudo", "curl", "-L", "https://github.com/docker/compose/releases/download/"+dockerComposeVersion+"/docker-compose-"+kernelname+"-"+machineHardwareName, "-o", "/usr/local/bin/docker-compose")
	execute("", "sudo", "chmod", "+x", "/usr/local/bin/docker-compose")

	execute("", "docker-compose", "--version")
}

func installGo(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Checking go...")
	if isCmdExists("go") {
		return
	}

	fmt.Println("Installing go...")
	execute(hlfdPath, "wget", "https://dl.google.com/go/go"+goVersion+".linux-amd64.tar.gz")
	execute(hlfdPath, "tar", "-xvf", "go"+goVersion+".linux-amd64.tar.gz")
	execute("", "sudo", "rm", "/usr/local/go", "-rf")
	execute(hlfdPath, "sudo", "mv", "go", "/usr/local")
	updateProfile := `
export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
`
	appendStringToFile(".profile", homeDir, updateProfile)
	// execute(hlfdPath, "echo", "'export GOROOT=/usr/local/go'>>~/.profile")
	// execute(hlfdPath, "echo", "'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH'>>~/.profile")
	// execute(hlfdPath, "source", path.Join(homeDir, ".profile"))
	// execute("", "go", "version")
}
