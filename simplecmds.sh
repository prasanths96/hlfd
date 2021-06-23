# CA
hlfd deploy ca -n test -t

# Org
hlfd org create -n org1 -m Org1MSP --ca-name test --ca-admin-user admin --ca-admin-pass adminpw 

# Peer
hlfd deploy peer -n first-peer -t --org-name org1 --ca-name test --ca-admin-user admin --ca-admin-pass adminpw --state-db CouchDB

# Orderer
hlfd deploy orderer -n first-orderer -t --org-name org1 --ca-name test --ca-admin-user admin --ca-admin-pass adminpw 

# export ca
hlfd export ca -n test -o .

# import ca
hlfd import ca -f test.tar
rm -rf test.tar

# export org
hlfd export org -n org1 -o .

# import org
hlfd import org -f org1.tar
rm -rf org1.tar