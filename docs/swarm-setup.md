# Swarm setup

Example commands to setup a swarm for development

## Docker machine

```bash

#Create manger+nodes
docker-machine create manager1
docker-machine create manager2
docker-machine create node1
docker-machine create node2

docker-machine ssh manager1
# Init swarm
docker swarm init --advertise-addr $(ifconfig eth1 | grep "inet addr" | awk '{print $2}' | sed s/addr://)

# Join command for worker
docker swarm join-token worker | grep "join"

# Join command for manager
docker swarm join-token worker | grep "join"



# Join swarm
docker-machine ssh node1
docker swarm join --token *** 192.168.99.100:2377
docker-machine ssh node2
docker swarm join --token *** 192.168.99.100:2377

# On manager1

docker node inspect node1

# set node labels
docker node update --label-add  virhal.color=blue node1
docker node update --label-add  virhal.color=red node2
docker service create  --constraint 'node.labels.virhal.color == blue' --name helloworld alpine ping docker.com
docker service inspect --pretty helloworld
docker service ps helloworld
# on node1

docker node update --label-add  virhal.color=red node1
docker node update --label-add  virhal.color=blue node2

docker service ps helloworld
# on node2

docker node update --label-add  virhal.color=white node2
docker service ps helloworld
# service state is pending

```
