#!/bin/bash

killall CA-Server

go build ./

echo "start ca server"
./CA-Server -RegenCert=false &
