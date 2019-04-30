# Scaling a web app

This repository is for learning how a web app can be scaled.

Subnet makes service discovery easier, as all the microservices are gonna run on this subnet, and it will be easier to
manually assign ip addresses.

## Installation
### Create subnet
- `docker network create --subnet=172.18.0.0/16 scalenet`

### Building image for the dataservice
- `docker build -f Dockerfile-dataservice -t scale/dataservice .`

### Building image for web backend
- `docker build -f Dockerfile-web -t scale/web .`

### Create containers
- Containers:
    ```
    docker run --name dataservice --ip=172.18.0.10 --net=scalenet -P --rm -it scale/dataservice
    docker run --name web --ip=172.18.0.11 --net=scalenet -P --rm -it -- scale/web --dataservice=http://172.18.0.10:4000
    ```
