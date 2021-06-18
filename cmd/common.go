package cmd

import (
	"os"
	"runtime"
)

type Object map[string]interface{}

var (
	GOOS = runtime.GOOS
	ARCH = runtime.GOARCH
)

var (
	//
	hlfdPath string

	//
	binFolder          = "bin"
	caClientHomeFolder = "ca-client-home"
	caClientName       = "fabric-ca-client"
	mspFolder          = "msp"
	tlsFolder          = "tls"
	//
	configTxGenName = "configtxgen"

	//
	commonFilUmask = os.FileMode(0700)

	//
	caDepFolder      = "cas"
	peerDepFolder    = "peers"
	ordererDepFolder = "orderers"
	orgCommonFolder  = "organizations"

	//
	caHomeFolder      = "ca-home"
	peerHomeFolder    = "peer"
	ordererHomeFolder = "orderer"
	//
	genesisFileName = "genesis.block"
	//
	orgInfoFileName = "info.json"

	// CA Environment
	CaAdminEnv     = "HLFD_CA_ADMIN_USER"
	CaAdminPassEnv = "HLFD_CA_ADMIN_PASS"

	// Couch Environment
	CouchAdminEnv     = "HLFD_COUCH_ADMIN_USER"
	CouchAdminPassEnv = "HLFD_COUCH_ADMIN_PASS"

	// Install prereqs
	//  dockerVersion = ""
	dockerComposeVersion = "1.29.2"
	goVersion            = "1.16.5"

	// HLF container constant internal configs (cannot change these by env variables)
	//  peerConstPort = 7054
	couchDBConstPort = 5984
)
