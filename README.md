# Webapper

Webapper provides a webserver for testing purposes.

## Docker Cluster

### Preparing the test plataform

Create two virtual machines with docker-machine  

```bash
docker-machine create --driver virtualbox --engine-label disk=ssd swarm1
docker-machine create --driver virtualbox --engine-label disk=hdd swarm2
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

### Deploying an application

Go ahead and deploy the first application on the private docker cloud

```bash
eval $(docker-machine env swarm1)
docker stack deploy --compose-file docker-stack.yml webapper
```

### Testing the scalibility

Run couple of curl requests and check the responses

```bash
INTERNAL_LOAD_BALANCE_IP=$(docker-machine ip swarm1)

while true; do
    curl "http://${INTERNAL_LOAD_BALANCE_IP}:8080/hostname"
    sleep 3
done
```

Afterwards scale up a service

```bash
docker service scale webapper_webapper=4
```

### Testing the placement

In order to move a running container to another swarm node

```bash
docker service update --constraint-add "engine.labels.disk==ssd" [SERVICE_ID]
```

### Deploying a new release

Note that the update process will be accord to the parallelism option 

```bash
docker service update --image rchicoli/webapper:0.0.3 [SERVICE_NAME]
```

## Kubernetes Cluster

### Setting up locally a single node Kubernetes cluster

To create a kubernetes cluster, run following command:

```bash
$ minikube start --vm-driver=virtualbox
```

For usability, make sure to load the kubectl completion code for a given shell

```
source <(kubectl completion zsh)
```

The cluster should be up and running, now it is time to start creating a service:

```bash
$ kubectl create -f k8s-service.yaml
```

Afterwards let's create a deployment:

```bash
$ kubectl create -f k8s-deployment.yaml
```

### Testing the running pods

```bash
INTERNAL_LOAD_BALANCE_IP=$(minikube ip)
curl "http://${INTERNAL_LOAD_BALANCE_IP}:30080/hostname"
```