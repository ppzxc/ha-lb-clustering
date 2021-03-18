#!/bin/bash

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
