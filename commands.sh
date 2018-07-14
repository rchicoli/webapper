#!/bin/bash

docker container list
docker system prune
docker system df -v
docker events

docker service create --name webapper rchicoli/webapper:0.0.3
docker service update --constraint-add "engine.labels.disk==ssd" SERVICE_ID
docker service scale SERVICE_NAME=1
docker service update --update-parallelism=1 --image rchicoli/webapper:0.0.3 SERVICE_ID
docker service ls
docker service update --secret-add foo SERVICE_ID
docker service update --publish-add 8080:8080 SERVICE_ID

docker stack deploy --compose-file docker-stack.yml webapper
docker deploy -h

docker plugins (external)

docker inspect --pretty SERVICE_ID

echo secretpass123 | docker secret create foo -

docker service create \
      --mount src=<VOLUME-NAME>,dst=<CONTAINER-PATH> \
      --name myservice \
      <IMAGE>
