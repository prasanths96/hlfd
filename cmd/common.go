package cmd

import "os"

type Object map[string]interface{}

//
var commonFilUmask = os.FileMode(0777)

//
var caDepFolder = "ca"
var peerDepFolder = "peer"
var ordererDepFolder = "orderer"

// CA Environment
var CaAdminEnv = "HLFD_CA_ADMIN_USER"
var CaAdminPassEnv = "HLFD_CA_ADMIN_PASS"
