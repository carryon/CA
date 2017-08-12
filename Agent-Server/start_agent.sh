#!/bin/bash

killall Agent-Server

go build ./

echo "start agent server"
./Agent-Server &
