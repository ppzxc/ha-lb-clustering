#!/bin/bash
export $(cat ./default.env)

docker stop $(docker ps -aq)
docker rm $(docker ps -aq)

if [ $(docker network ls | grep halb-cluster | wc -l) -eq 0 ];then
	docker network create halb-cluster
fi

cd ./fluentd
make d
make u

sleep 1s

cd ../proxy
make d
make u

sleep 1s

cd ../rabbitmq
make d
make u

sleep 1s

cd ../manager
make d
make u
