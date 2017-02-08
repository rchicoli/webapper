# Webapper

Webapper provides a webserver for testing purposes.

## Preparing the test plataform

Create two virtual machines with docker-machine  

```bash
docker-machine create --driver virtualbox --engine-label disk=ssd swarm1
docker-machine create --driver virtualbox swarm2
```

Log in into the machine and initialize the swarm cluster

```
docker-machine ssh swarm1
docker swarm init --advertise-addr 192.168.99.100
```

Log in into the next node and join it to the swarm cluster

```bash
docker-machine ssh swarm2
docker swarm join \
    --token SWMTKN-1-19xe2kzne8i45oh6qac9mnyhnu0u88d-bpuq0e7hbpmfo5xncxs1 \
    192.168.99.100:2377
```

## Deploying an application

Go ahead and deploy the first application on the private docker cloud

```bash
eval $(docker-machine env swarm1)
docker stack deploy --compose-file docker-stack.yml webapper
```

## Testing the scalibility

Run couple of curl requests and check the responses

```bash
INTERNAL_LOAD_BALANCE_IP=$(docker-machine config swarm1| sed -nr 's#-H=tcp://(.*):.*#\1#p')

while true; do
    curl "http://${INTERNAL_LOAD_BALANCE_IP}:8080/hostname"
    echo ""
    sleep 3
done
```

Afterwards scale up a service

```bash
docker service scale webapper_webapper=4
```
