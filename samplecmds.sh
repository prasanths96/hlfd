hlfd deploy peer -n first-peer -t -m Org1MSP --ca-name test --ca-addr https://localhost:8054 --ca-admin-user admin --ca-admin-pass adminpw --ca-tls-cert-path ~/.hlfd/ca/test/ca-home/tls-cert.pem --state-db CouchDB 

hlfd deploy orderer -n first-orderer -t -m OrdererMSP --orderer-addr=orderer.example.com --ca-name test --ca-addr https://localhost:8054 --ca-admin-user admin --ca-admin-pass adminpw --ca-tls-cert-path ~/.hlfd/ca/test/ca-home/tls-cert.pem

hlfd org create -n org1 -m Org1MSP --ca-name test --ca-addr https://localhost:8054 --ca-admin-user admin --ca-admin-pass adminpw --ca-tls-cert-path ~/.hlfd/ca/test/ca-home/tls-cert.pem