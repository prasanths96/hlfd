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

# stop 
hlfd stop ca -n test 
hlfd stop peer -n first-peer
hlfd stop orderer -n first-orderer

# You can delete the docker containers, volumes, networks all you want, it wont matter

# resume
hlfd resume ca -n test
hlfd resume peer -n first-peer
hlfd resume orderer -n first-orderer

# terminate
hlfd terminate ca -n test
hlfd terminate peer -n first-peer
hlfd terminate orderer -n first-orderer

# now all data is also deleted, theres no going back