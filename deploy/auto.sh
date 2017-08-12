#!/bin/bash
rm -rf /tmp/leveldb/
killall deploy

go build ./

./deploy  newagent --aid=123 --remark="this is a test"
./deploy  newconfig --nodeType=lcnd --nodeID=0001_abc  --config=nodeConfigs/lcnd.yaml
./deploy  newconfig --nodeType=lcnd --nodeID=0001_def  --config=nodeConfigs/lcnd.yaml

./deploy  relate  --nodeType=lcnd --aid=123 --nodeIDs=0001_abc,0001_def
./deploy  updateVersion --nodeType=lcnd --version=v1.0.9

# msg-net

./deploy  newagent --aid=456 --remark="this is a test of msg-net"
./deploy  newconfig --nodeType=msg-net --nodeID=msg-net_abc  --config=msgNetConfigs/msg-net.yaml
./deploy  newconfig --nodeType=msg-net --nodeID=msg-net_def  --config=msgNetConfigs/msg-net.yaml
./deploy  relate --nodeType=msg-net --aid=456 --nodeIDs=msg-net_abc,msg-net_def
./deploy  updateVersion --nodeType=msg-net --version=v1.0.9



echo "start server"
./deploy serve --port=8080&
sleep 2s

curl -X POST --data '{"jsonrpc":"2.0","method":"nodes-config","params":["123"],"id":1}' http://localhost:8080
curl -X POST --data '{"jsonrpc":"2.0","method":"config-timestamp","params":["123"],"id":1}' http://localhost:8080
curl -X POST --data '{"jsonrpc":"2.0","method":"lcnd-version","params":["123"],"id":1}' http://localhost:8080



curl -X POST --data '{"jsonrpc":"2.0","method":"msgnet-config","params":["456"],"id":1}' http://localhost:8080
curl -X POST --data '{"jsonrpc":"2.0","method":"msgnet-timestamp","params":["456"],"id":1}' http://localhost:8080
curl -X POST --data '{"jsonrpc":"2.0","method":"msgnet-version","params":["456"],"id":1}' http://localhost:8080

