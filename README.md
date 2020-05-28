# 0proxy

## Setup

Clone the repo and run the following command inside the cloned directory

```
$ ./docker.local/bin/init.sh
```

## Building and Starting the Node

If there is new code, do a git pull and run the following command

```
$ ./docker.local/bin/build.sh
```

Go to the bin directory (cd docker.local/bin) and run the container using

```
$ ./start.sh
```

## Point to another blockchain

You can point the server to any instance of 0chain blockchain you like, Just go to config (docker.local/config) and update the 0proxy.yaml.

```
block_worker: http://198.18.0.98:9091
```

We use blockWorker to connect to the network instead of giving network details directly, It will fetch the network details automatically from the blockWorker's network API.

There are other configurable properties as well which you can update as per the requirement.

### Cleanup

Get rid of old data when the blockchain is restarted or if you point to a different network:

```
$ ./docker.local/bin/clean.sh
```

### Network issue

If there is no test network, run the following command

```
docker network create --driver=bridge --subnet=198.18.0.0/15 --gateway=198.18.0.255 testnet0
```
