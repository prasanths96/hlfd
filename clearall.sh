docker rm -f $(docker ps -aq)
docker volume prune -f
docker system prune -f
docker image prune -f

# sudo rm ~/.hlfd/bin -rf
sudo rm ~/.hlfd/cas -rf
sudo rm ~/.hlfd/ca-client-home -rf
sudo rm ~/.hlfd/orderers -rf
sudo rm ~/.hlfd/peers -rf
sudo rm ~/.hlfd/organizations -rf
sudo rm ~/.hlfd/imports -rf
sudo rm ~/.hlfd/exports -rf