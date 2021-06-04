package cmd

import "os"

type Object map[string]interface{}

//
var hlfdPath string

//
var commonFilUmask = os.FileMode(0777)

//
var caDepFolder = "ca"
var peerDepFolder = "peer"
var ordererDepFolder = "orderer"

//
var caHomeFolder = "ca-home"

// CA Environment
var CaAdminEnv = "HLFD_CA_ADMIN_USER"
var CaAdminPassEnv = "HLFD_CA_ADMIN_PASS"

// Install prereqs
// var dockerVersion = ""
var dockerComposeVersion = "1.29.2"
var goVersion = "1.16.5"
