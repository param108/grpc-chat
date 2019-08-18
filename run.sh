#!/bin/bash -x
kill -9 `cat out/PID`; 
rm -rf out
tar -zxvf bundle.tgz;
cd out
./grpc_chat migrate

if [ $? -ne 0 ]
then
    echo "Failed to run Migrations"
    exit
fi

nohup ./grpc_chat server >> ~/nohup.out < /dev/null  &
sleep 10
