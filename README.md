# Scaling a web app

A repository for learning how to scale a web application, concentrates over using docker containers over a bridge network for application instances that are being load balanced using the Round Robin algorithm. The application instances interact with a data service and a logging service that run in their separate containers. 

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
>   The docker network inspect command should show both the containers in the Containers object of the configuration, see the eg below
- An output for docker network inspect scalenet
    ```
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
                "2003ea2193d36f79916e6b4af1eaec0ee6cce219aef3982cf09ec0a3dbf52cd3": {
                    "Name": "dataservice",
                    "EndpointID": "89c3ed067d996464141dd4923b4389d9af03caf2fa944caf2c99c6e6d0dcceb6",
                    "MacAddress": "02:42:ac:12:00:0a",
                    "IPv4Address": "172.18.0.10/16",
                    "IPv6Address": ""
                }
            },
            "Options": {},
            "Labels": {}
        }
    ]
    ```