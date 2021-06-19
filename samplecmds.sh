# CA
hlfd deploy ca -n test -t

# Org
hlfd org create -n org1 -m Org1MSP --ca-name test --ca-addr https://localhost:8054 --ca-admin-user admin --ca-admin-pass adminpw --ca-tls-cert-path ~/.hlfd/cas/test/ca-home/tls-cert.pem

# Peer
hlfd deploy peer -n first-peer -t --org-name org1 --ca-name test --ca-admin-user admin --ca-admin-pass adminpw --state-db CouchDB 

# Orderer    !!!MUST USE EXTERNALLY ACCESSIBLE ADDRESS, not localhost (--orderer-addr)
hlfd deploy orderer -n first-orderer -t --org-name org1 --orderer-addr=localhost:7050 --ca-name test --ca-admin-user admin --ca-admin-pass adminpw 