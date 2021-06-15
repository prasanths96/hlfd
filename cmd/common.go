package cmd

import "os"

type Object map[string]interface{}

//
var hlfdPath string

//
var binFolder = "bin"
var caClientHomeFolder = "ca-client-home"
var caClientName = "fabric-ca-client"
var mspFolder = "msp"
var tlsFolder = "tls"

//
var commonFilUmask = os.FileMode(0700)

//
var caDepFolder = "ca"
var peerDepFolder = "peer"
var ordererDepFolder = "orderer"

//
var caHomeFolder = "ca-home"
var peerHomeFolder = "peer"
var ordererHomeFolder = "orderer"

// CA Environment
var CaAdminEnv = "HLFD_CA_ADMIN_USER"
var CaAdminPassEnv = "HLFD_CA_ADMIN_PASS"

// Couch Environment
var CouchAdminEnv = "HLFD_COUCH_ADMIN_USER"
var CouchAdminPassEnv = "HLFD_COUCH_ADMIN_PASS"

// Install prereqs
// var dockerVersion = ""
var dockerComposeVersion = "1.29.2"
var goVersion = "1.16.5"

// HLF container constant internal configs (cannot change these by env variables)
// var peerConstPort = 7054
var couchDBConstPort = 5984
