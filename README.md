# Scaling a web app

This repository is for learning how a web app can be scaled.

Subnet makes service discovery easier, as all the microservices are gonna run on this subnet, and it will be easier to
manually assign ip addresses.

docker network create --subnet=172.18.0.0/16 scalenet

Run the data service using the below mentioned docker commands
docker build -f Dockerfile-dataservice -t scale/dataservice .


Run the image and create a container
docker run --name dataservice --ip=172.18.0.10 --net=scalenet -P --rm -it scale/dataservice


Building image for web backend
docker build -f Dockerfile-web -t scale/web .

docker run --name web --ip=172.18.0.11 --net=scalenet -P --rm -it -- scale/web --dataservice=http://172.18.0.10:4000
