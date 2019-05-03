# Scaling a web app

This repository is for learning how a web app can be scaled.

## Installation
### Create subnet
> Subnet makes service discovery easier, as all the microservices are gonna run on this subnet, and it will be easier to manually assign ip addresses.
- `docker network create --subnet=172.18.0.0/16 scalenet`

### Building image for the dataservice
- `docker build -f Dockerfile-dataservice -t scale/dataservice .`

### Building image for web backend
- `docker build -f Dockerfile-web -t scale/web .`

### Create containers
- Containers:
    ```
    docker run --name dataservice --ip=172.18.0.10 --net=scalenet -p 4000:4000 --rm -it -d scale/dataservice
    docker run --net=scalenet --ip=172.18.0.2 -p3000:3000 --name webtest --rm -it scale/web --dataservice=https://172.18.0.10:4000
    ```
- Check that both the containers are up and running using the below-mentioned commands:
    ```
    docker ps
    docker network inspect scalenet
    ```
>   The network inspect should show both the containers in the Containers object of the configuration