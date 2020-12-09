#!/bin/bash

cd build

LISTEN_ADDRESS=:23334 STORAGE_ROOT=./data1 nohup ./dataserver > dataserver1.out 2>&1 & echo $! > ./var/datarun1.pid
LISTEN_ADDRESS=:23335 STORAGE_ROOT=./data2 nohup ./dataserver > dataserver2.out 2>&1 & echo $! > ./var/datarun2.pid
LISTEN_ADDRESS=:23336 STORAGE_ROOT=./data3 nohup ./dataserver > dataserver3.out 2>&1 & echo $! > ./var/datarun3.pid
LISTEN_ADDRESS=:23337 STORAGE_ROOT=./data4 nohup ./dataserver > dataserver4.out 2>&1 & echo $! > ./var/datarun4.pid
LISTEN_ADDRESS=:23338 STORAGE_ROOT=./data5 nohup ./dataserver > dataserver5.out 2>&1 & echo $! > ./var/datarun5.pid
LISTEN_ADDRESS=:23339 STORAGE_ROOT=./data6 nohup ./dataserver > dataserver6.out 2>&1 & echo $! > ./var/datarun6.pid

nohup ./interserver > interserver.out 2>&1 & echo $! > ./var/interfacerun.pid

ps -aux | grep dataserver
ps -aux | grep interserver

cd ..