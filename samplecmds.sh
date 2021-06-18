# CA
hlfd deploy ca -n test -t

# Org
hlfd org create -n org1 -m Org1MSP --ca-name test --ca-addr https://localhost:8054 --ca-admin-user admin --ca-admin-pass adminpw --ca-tls-cert-path ~/.hlfd/cas/test/ca-home/tls-cert.pem

# Peer
hlfd deploy peer -n first-peer -t -m Org1MSP --ca-name test --ca-addr https://localhost:8054 --ca-admin-user admin --ca-admin-pass adminpw --ca-tls-cert-path ~/.hlfd/cas/test/ca-home/tls-cert.pem --state-db CouchDB 

# Orderer
hlfd deploy orderer -n first-orderer -t --org-name org1 --orderer-addr=localhost:7051 --ca-name test --ca-admin-user admin --ca-admin-pass adminpw 