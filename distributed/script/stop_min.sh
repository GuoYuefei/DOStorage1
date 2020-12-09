#!/bin/bash

cd build

echo 'before stop'
ps -aux | grep dataserver
ps -aux | grep interserver

kill `cat var/datarun1.pid`
kill `cat var/datarun2.pid`
kill `cat var/datarun3.pid`
kill `cat var/datarun4.pid`
kill `cat var/datarun5.pid`
kill `cat var/datarun6.pid`

kill `cat var/interfacerun.pid`

echo 'after stop'
ps -aux | grep dataserver
ps -aux | grep interserver

cd ..