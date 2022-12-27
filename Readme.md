Steps to clone the docker setup

**Docker pull :**
```
$ docker pull dockeralphazero/tf-apache-prior-todocker-push-v6:firstdockerpush
```

**Clone the golang repo:**
```
$ git clone git@bitbucket.org:sachins/connect4zero-golang.git connect4-selfplay-golang
```

**Clone the python tensorflow backend:**
```
$git clone git@bitbucket.org:sachins/connect4zero.git connect4-selfplay-python
```

**Start docker with the python and golang repo mounted.**

- Finding the docker image id which was pulled from dockerhub

```
$ docker images
REPOSITORY                                         TAG               IMAGE ID       CREATED         SIZE
dockeralphazero/tf-apache-prior-todocker-push-v6   firstdockerpush   ad94de7a61ae   5 months ago    7.25GB
```

- Start the docker image as a container with the local golang and python directories mounted

```
docker run --publish=0.0.0.0:8888:8888/tcp --publish=0.0.0.0:1234:80/tcp --publish=6006:6006 --env="DISPLAY" --volume="/tmp/.X11-unix:/tmp/.X11-unix:rw" --gpus all -u 1000:1000 -d --name tf-volume-test --mount type=bind,source=**PATH_TO_PYTHON_REPO**,target=/python  --mount type=bind,source=**PATH_TO_GOLOANG_REPO**,target=/golang -it **ad94de7a61ae**
```

Replace the bold parts with your local folders of Python and goloang repos we pulled from bitbucket earlier.
Then replace ad94de7a61ae with the docker image id you get from 'docker images' command which was run before.

------------------------------------------------------------------------------------------------------------------------
**Golang client setup**

Once thats done You would want to install go on your linux system first,

https://golang.org/doc/install

Then :
```
cd /path/to/connect4-selfplay
export GOPATH="/path/to/connect4-selfplay"
```
**Build a automated test program, which tests standard connect-4 positions and checks AlphaZero responses**
```
go build src/test/test-connect4-moves-zero.go
```
**run the test program**
```
./test-connect4-moves-zero
```
**If you would like to play AlphaZero then build this code**
```
go build src/humanPlay-connect4zero.go
```
**Play yourself agains AlphaZero**
```
./humanPlay-connect4zero
```