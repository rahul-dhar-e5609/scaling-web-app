# Scaling a web app

A repository for learning how to scale a web application, concentrates over using docker containers over a bridge network for application instances that are being load balanced using the Round Robin algorithm. The application instances interact with a data service and a logging service that run in their separate containers. 

## Installation
### Create subnet
> Subnet makes service discovery easier, as all the microservices are gonna run on this subnet, and it will be easier to manually assign ip addresses.
- `docker network create --subnet=172.18.0.0/16 scalenet`

### Building image for the dataservice
- `docker build -f .docker/dataservice.dockerfile -t scale/dataservice .`

### Building image for the loadbalancer
- `docker build -f .docker/loadbalancer.dockerfile -t scale/loadbalancer .`

### Building image for cache service
- `docker build -f .docker/cacheservice.dockerfile -t scale/cacheservice .`

### Building image for log service
- `docker build -f .docker/logservice.dockerfile -t scale/logservice .`

### Building image for web backend
- `docker build -f .docker/web.dockerfile -t scale/web .`

### Create containers
- Containers:
    ```
    docker run --name dataservice --ip=172.18.0.10 --net=scalenet -p 4000:4000 --rm -it -d scale/dataservice
    docker run --name loadbalancer --ip=172.18.0.12 --net=scalenet -p 2000:2000 --rm -it -d scale/loadbalancer
    docker run --name cacheservice --ip=172.18.0.13 --net=scalenet -p 5000:5000 --rm -it scale/cacheservice
    docker run --name logservice --ip=172.18.0.14 --net=scalenet -p 6000:6000 --rm -it -v /Users/rahuldhar/log:/log scale/logservice
    docker run --net=scalenet --ip=172.18.0.2 -p3000:3000 --name webtest --rm -it scale/web --dataservice=https://172.18.0.10:4000 --loadbalancer=https://172.18.0.12:2001 --cachingservice=https://172.18.0.13:5000 --logservice=https://172.18.0.14:6000
    ```
- Check that both the containers are up and running using the below-mentioned commands:
    ```
    docker ps
    docker network inspect scalenet
    ```
>   The docker network inspect command should show both the containers in the Containers object of the configuration, see the eg below
- An output for docker network inspect scalenet
    ```
    Rahuls-MacBook-Pro:scaling-web-app rahuldhar$ docker network inspect scalenet
    [
        {
            "Name": "scalenet",
            "Id": "01f8a6ee533b66b0775b4ad6fb99c825dbafd8aba963459d00fb835940ae9eab",
            "Created": "2019-04-30T19:45:53.380381186Z",
            "Scope": "local",
            "Driver": "bridge",
            "EnableIPv6": false,
            "IPAM": {
                "Driver": "default",
                "Options": {},
                "Config": [
                    {
                        "Subnet": "172.18.0.0/16",
                        "Gateway": "172.18.0.1"
                    }
                ]
            },
            "Internal": false,
            "Attachable": false,
            "Ingress": false,
            "ConfigFrom": {
                "Network": ""
            },
            "ConfigOnly": false,
            "Containers": {
                "a7c1fe39219c3eb41f69ef6965dd8b0911eb69067cc22e8171e30958f47d1cf5": {
                    "Name": "dataservice",
                    "EndpointID": "e3696f2c570ac122526c94e82fe258faa58d58a2f669e513b3c433c8f2b81df2",
                    "MacAddress": "02:42:ac:12:00:0a",
                    "IPv4Address": "172.18.0.10/16",
                    "IPv6Address": ""
                },
                "c7bffd2989d8534671005e43f23b9ad45fc2d67d351a6828c12a049bfceacc6d": {
                    "Name": "webtest",
                    "EndpointID": "000995e02495a6ea1455588af0747a9423bafba6761dddffd7f31513eab1c571",
                    "MacAddress": "02:42:ac:12:00:02",
                    "IPv4Address": "172.18.0.2/16",
                    "IPv6Address": ""
                },
                "f402bfa13e2491a5de1e08af880c80a59ffd2c4919400cf9c9a841f17202925c": {
                    "Name": "loadbalancer",
                    "EndpointID": "b8d56b9813d55ddcac5a12853ccd14ba2b911b4c8609537519b1c56d85904f9e",
                    "MacAddress": "02:42:ac:12:00:0c",
                    "IPv4Address": "172.18.0.12/16",
                    "IPv6Address": ""
                }
            },
            "Options": {},
            "Labels": {}
        }
    ]
    ```