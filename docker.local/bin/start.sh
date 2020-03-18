#!/bin/sh
PWD=`pwd`

echo Starting 0Proxy ...

docker-compose -p 0proxy -f ../docker-compose.yml up
