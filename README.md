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

By default the are pointing to one network
```
miners:
  - http://one.devnet-0chain.net:31071
  - http://one.devnet-0chain.net:31072
  - http://one.devnet-0chain.net:31073
  - http://one.devnet-0chain.net:31074
  - http://one.devnet-0chain.net:31075
  - http://one.devnet-0chain.net:31076
  - http://one.devnet-0chain.net:31077
  - http://one.devnet-0chain.net:31078
  - http://one.devnet-0chain.net:31079
sharders:
  - http://one.devnet-0chain.net:31171
  - http://one.devnet-0chain.net:31172
  - http://one.devnet-0chain.net:31173
  - http://one.devnet-0chain.net:31174
  - http://one.devnet-0chain.net:31175
  - http://one.devnet-0chain.net:31176
  - http://one.devnet-0chain.net:31177
  - http://one.devnet-0chain.net:31178
  - http://one.devnet-0chain.net:31179
  ```

You need to set miners and sharders of the blockchain you want to connect in 0proxy.yaml, There are other configurable properties as well which you can update as per the requirement.

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