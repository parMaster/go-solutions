    $ docker run <image> 
Pull **image** and run **container** from docker hub

Command locates image priority:
Host \> Docker hub


    docker run docker/whalesay cowsay Hello-world!
    docker ps
Currently running containers

    docker ps -a
Currently and previously running containers

    docker stop <container id|name> 
    docker rm <containter>
To remove a container 

    docker images
List of available (pulled) locally images

    docker rmi <image>
Delete local image

Must stop all dependent containers to delete an image 

    docker pull <image>
To preheat image - download so docker run will run faster

    docker run ubuntu
Runs ubuntu image and exits immediately

Containers aint't suppose to run OS, but to run a process (service, db, api, script)

    docker run ubuntu sleep 5
Will sleep for 5 seconds and exit

    docker exec <container_name> cat /etc/hosts
Runs cat /etc/hosts on container <container_name> and returns result

    docker run <image>:<tag>
<tag> is a version of image to run.
Latest by default
    docker run redis:4.0 --- will run redis version 4.0  

## Attached vs Detached modes

    docker run kodekloud/simple-webapp
Runs in attached mode.

DOESN'T listen to STDIN. See **Interactive mode**

    docker run -d kodekloud/simple-webapp
Detached mode (-d) 

    docker attach <name or id of container> 
<container_id> can be shortened to first N chars

## Interactive mode

    docker run -i <image>
Maps local STDIN to application

    docker run -it <image>
Maps local STDIN to the container terminal

## Run PORT mapping
    docker run -p <ext>:<int> <image>
<ext> - docker host port
<int> - application listening port 

## Run VOLUME mapping
    docker run -v <ext>:<int> <image>
<ext> - host path
<int> - container path

Like:

    docker run -v /opt/datadir:/var/lib/mysql mysql

Runs mysql and maps internal container storage path to container host path /opt/datadir

## Inspect container
    docker inspect <container id|name>
Provides info about container: state, mounts, config

## Container Logs
    docker logs <container id|name>

## Run Environment Variables
    docker run -e <VAR>=<value> <image>

# Docker Build
    docker build -t <tag> <path>
Build a container from <path> and name it <tag>


# Dockerfile

    sudo apt --assume-yes install <package>
    sudo apt install -y <package>
Assume Yes to all queries and do not prompt


## CMD vs ENTRYPOINT

### CMD

    CMD <command param>
    CMD sleep 5
or

    CMD ["command", "param"]
    CMD ["sleep", "5"]
Command to run in a container. **Overwritten** if called like:
    
    docker run ubuntu-sleeper sleep 10

### ENTRYPOINT

    ENTRYPOINT ["command"]
Allows to specify param during launch

    ENTRYPOINT ["sleep"] 
    docker run ubuntu-sleeper 10
will run
    sleep 10
at container startup. Command line param will be **appended**

### Combining both
Allows to set a **value on startup and a default value** as well.

    ENTRYPOINT ["sleep"]
    CMD ["5"]
will run

    sleep 5
at startup. But if argument is specified, like:
    
    docker run ubuntu-sleeper 10
then startup command be like:

    sleep 10

### Overwriting command AND param at startup

    docker run --entrypoint <new_entrypoint> <image> <param>
    docker run --entrypoint sleep2.0 ubuntu-sleeper 10

# Networking

Three networks:
- Bridge (default network). Internal IPs. Mapping internal ports to extrnal required to acces container apps
- none - not attached, no network access, isolated network
- host - container uses host network. No network isolation, no port mapping required.

Can be specified at launch:

    docker run ubuntu --network=none
    docker run ubuntu --network=host

Create new internal network

    docker network create \
        --driver bridge \
        --subnet 182.18.00/16
        custom-isolated-network

    docker network ls
to list networks

    docker inspect <container>
to check network config of the container

## built-in DNS server 
    mysql.connect(<container>)
will resolve \<container\> on internal DNS server


# Docker storage
Layered architecture, building from cache, reuse layers etc.

**Container layer** - Writable layer on top of the image layers. 

## Copy-on-write 
When app mades change to the image file, docker makes a copy in the container layer. Alive while container is alive, then destroyed when container stopped

## Volumes
    docker volume create <volume>
    docker volume create data_volume
Creates volume and folder /var/lib/docker/volumes

### volume mounting
    docker run -v <image_volume>:<container_folder> <image>
    docker run -v data_volume:/var/lib/mysql mysql
Mounts **"image" volume** (writable layer on top of the image layer) to **container folder**

### bind mounting
    docker run -v <image_folder>:<container_folder> <image>
    docker run -v /data/mysql:/var/lib/mysql mysql
mysql image will run and use /var/lib/mysql to store data, so container have the data in /data/mysql folder as long as it runs

## New mounting style (!)
    docker run \
        --mount type=bind,source=/data/mysql,target=/var/lib/mysql mysql
much more readable

# Docker Compose
To run application stack on a single docker host
Yaml config files like this:

    services:
        web:
            image: "mmushad/single-webapp"
        database:
            image: "mongodb"
        messaging:
            image: "redis:alpine"
        Orchestration:
            image: "ansible"
    
    
    docker-compose up


## Sample application - voting application
Will be used to demonstrate

    voting app -> in-memory DB -> worker -> db -> result-app

Would require to run containers like this:

    docker run -d --name=redis redis
    docker run -d --name=db postgres:9.4 result-app
    docker run -d --name=vote -p 5000:80 voting-app
    docker run -d --name=result -p 5001:80
    docker run -d --name=worker worker

Adding **DEPRECATED** link option to voting app:

    docker run -d --name=vote -p 5000:80 --link redis:redis voting-app

--link redis:redis - creates an **entry in the hosts file** so redis becomes resolvable host name

Finally, **DEPRECATED** links for every service in this example:

    docker run -d --name=redis redis
    docker run -d --name=db postgres:9.4 --link db:db result-app
    docker run -d --name=vote -p 5000:80 --link redis:redis voting-app
    docker run -d --name=result -p 5001:80
    docker run -d --name=worker --link db:db --link redis:redis worker

Since all this is just an illustration of the concept, we'll be using docker-compose instead.

# Docker compose

## Building from docker images:

**docker-compose.yml version:1**

    redis:
        image: redis
    db:
        image: postgres:9.4
    vote:
        image:voting-app
        ports:
            - 5000:80
        links:
            - redis
    result:
        image: result-app
        ports:
            - 5001:80
        links:
            - db:db
    worker: 
        image: worker
        links:
            - redis
            - db

Asuming links names is happening for worker container:
    links: 
        - db:db
            is the same as
        - db

## Building from application folders

**docker-compose.yml version:1**

    redis:
        image: redis
    db:
        image: postgres:9.4
    vote:
        build: ./vote
        ports:
            - 5000:80
        links:
            - redis
    result:
        build: ./result
        ports:
            - 5001:80
        links:
            - db:db
    worker: 
        build: ./worker
        links:
            - redis
            - db

_./vote, ./result, ./worker - application folders with Dockerfiles_
version 1 is on the default Bridged network


## Docker-compose versions

### V2
    - links automatically created
    - images run on Bridged network and can talk to each other

    version: 2
    services:
        redis:
            image: redis

            networks:
                - back-end
        db:
            image: postgres:9.4
            networks:
                - front-end
                - back-end
        vote:
            image: voting-app
            ports:
                - 5000:80
            depends_on:
                - redis
        result:
            image: result
            networks:
                - front-end
                - back-end
    networks:
        front-end:
        back-end:

Here we also specified two networks 
    - **front-end** is resolved from the outside - to service user requests
    - **back-end** is for services to talk through

# Coding exercices

# Registry
Center repository of docker images

    docker run nginx
    or
    docker run nginx/nginx
    assumed to be
    docker run docker.io/nginx/nginx

Private registry can be provided by service provider (Amazon ,DO)

    docker login private-registry.io
    docker run private-registry.io/apps/internal-app

obviously login and run a container form a prvate registry image

### Run a private registry
    
    docker run -d -p 5000:5000 --name registry registry:2

    docker image tag   my-image localhost:5000/my-image

    docker push localhost:5000?my-image

Then the image is accessible if host is accessible:

    docker pull 192.168.56.100:5000/my-image

# Docker Engine

- Docker CLI
- REST API
- Docker Deamon

    docker -H=remote-docker-engine:2375
runs a Docker CLI and connects to engine on a remote host

    docker -H=10.123.2.1:2375 run nginx
runs nginx image container on 10.123.2.1 

Isolation with namespaces
- PID name spaces - PID namespace on the host maps to PID namespace on container

    docker run --cpus=.5 ubuntu
limits CPU usage for container to 50%

    docker run --memory=100m ubuntu
limits memory for container to 100MB

# Container orchestration

    docker service create --replicas=100 nodejs
command for Docker Swarm

- Docker swarm - lacks autoscaling features, easy to setup
- Kubernetes - standart 
- Mesos - hard to setup and get started, many features

## Docker swarm. Quick look
Combiune multiple Docker machines together in a single cluster. Swarm will take care distributing and load balancing.

To **start a Swarm Manager Host**:

    docker swarm init 

To **join Nodes** (Workers) to swarm, run on the node:

    docker swarm join --token \<token\>

# Kubernetes. Quick look

Powerful shit. 

Kubernetes uses Docker host to host applications in the form of Docker containers.

Nodes - worker machines

Nodes cluster managed by Master node and responsible for orchestration

    kubectl run
    kubectl cluster-info
    kubectl get nodes
