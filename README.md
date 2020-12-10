# 0proxy

0proxy is used to do CRUD operations on 0chain via web interface using REST APIs. It uses GoSDK internally and exposes the SDK methods in the form of [APIs](#api).
You can find the [API documentation](https://0chain.net/page-documentation.html) on 0chain website.

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

## API

### Upload

To upload OR update a file to 0chain network.

Path : `/upload`

Details:

- https://0chain.net/page-documentation.html#tag/0proxy/paths/~1upload/post
- https://0chain.net/page-documentation.html#tag/0proxy/paths/~1upload/put

### Download

To download a file from 0chain network

Path: `/download`

Details: https://0chain.net/page-documentation.html#tag/0proxy/paths/~1download/get

### Stream

To stream a file from 0chain network

Path: `/stream`

Details: https://0chain.net/page-documentation.html#tag/0proxy/paths/~1stream/get

### Delete

To delete a file from 0chain network

Path: `/delete`

Details: https://0chain.net/page-documentation.html#tag/0proxy/paths/~1delete/delete

### Copy

To copy a file on 0chain network

Path: `/copy`

Details: https://0chain.net/page-documentation.html#tag/0proxy/paths/~1copy/put

### Rename

To rename a file on 0chain network

Path: `/rename`

Details: https://0chain.net/page-documentation.html#tag/0proxy/paths/~1rename/put

### Move

To move a file on 0chain network

Path: `/move`

Details: https://0chain.net/page-documentation.html#tag/0proxy/paths/~1move/put

### Share

To share a file on 0chain network

Path: `/share`

Details: https://0chain.net/page-documentation.html#tag/0proxy/paths/~1share/put
