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




