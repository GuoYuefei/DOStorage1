# Distributed Object Storage
[![Build Status](https://travis-ci.com/GuoYuefei/DOStorage1.svg?branch=master)](https://travis-ci.com/GuoYuefei/DOStorage1) [![GitHub license](https://img.shields.io/github/license/GuoYuefei/DOStorage1)](https://github.com/GuoYuefei/DOStorage1/blob/master/LICENSE) ![language](https://img.shields.io/github/languages/top/GuoYuefei/DOStorage1) ![last](https://img.shields.io/github/last-commit/GuoYuefei/DOStorage1.svg)

- [Getting Started](#getting-started)
	+ [Prepare Components](#prepare-components)
  + [Environmental configuration](#environmental-configuration)
  + [Usage](#usage)

- [Lisence](#License)
- [Preface](#Preface)

## Getting Started

### Prepare Components
The first component to be laid out is **RabbitMQ**

Download RabbitMQ and erlang.

In rabbitmqctl command

```powershell
rabbitmqctl add_user test test
rabbitmqctl set_permissions -p / test ".*" ".*" ".*"
```
The second component to be laid out is **ElasticSearch 7.x**

Refer to this URL to install and use the default configuration to open it. https://www.elastic.co/guide/en/elastic-stack-get-started/current/get-started-elastic-stack.html

Then <code>go run readyfordistributed.go</code> for prepare MQ and ES.

Use <code>rabbitmqctl list_exchanges</code> to view the results. If there are two exchanges, apiServers and dataServers, the above program has been successfully executed.

Use [Kibana](https://www.elastic.co/products/kibana) to view and manage your ES.

### Environmental configuration

|Variable name| interface | data |
|:---:|:---------:|:----:|
|RABBITMQ_SERVER|need / default "amqp://test:test@localhost:5672/"|need / default same with interface|
|LISTEN_ADDRESS|need / default ":23333"|need / default ":23334"|
|STORAGE_ROOT|no need|need / default "./data"|
|ES_SERVER|need / "localhost:9200"|Not sure / "localhost:9200"|

-----
You can use default values **or** set environment variables. 

```powershell
// powershell
$env:RABBITMQ_SERVER = "amqp://test:test@192.168.1.68:5672/"
$env:LISTEN_ADDRESS = "192.168.1.68:23333"
$env:STORAGE_ROOT = YOU_PATH(Where you want to store the object)
./interface.exe or ./data.exe 
```
```shell script
// shell
export RABBITMQ_SERVER="amqp://test:test@localhost:5672/"
export LISTEN_ADDRESS="localhost:23333"
export STORAGE_ROOT=YOU_PATH(Where you want to store the object)
```

```shell script
// shell
LISTEN_ADDRESS=:12345 STORAGE_ROOT=./tmp RABBITMQ_SERVER="amqp://test:test@localhost:5672/" go run distributed/data/main.go
```

### Usage

~~Please **create a directory of STORAGE_ROOT**  in advance and **create an folders named *objects* ** in that directory.~~

The configuration content mentioned above needs to be configured on each node machine as needed.

Then you can find the corresponding binary files data.exe and interface.exe in the release. Run the two programs on different node machines through the command line. 

**After this project version v0.2,  six data server nodes must be equipped, otherwise the service will be unavailable.**

You can compile it yourself with the source code.

**This is the external RESTful interface table of the interface server：**

| Http Method | URL                                     | Param                                                        | Effect                                                       |
| ----------- | --------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| GET         | "http://"+APIHOST+"/objects/"+"[name]"  | [version=「int」]                                            | Get object named [name]， if version is empty, will return object that is Latest Version |
| PUT         | "http://"+APIHOST+"/objects/"+"[name]"  | 1. Request.Body = object content 2. Content-length=len(object) 3. Digest="SHA-256=「object's hash base64」" | Put object                                                   |
| DELETE      | "http://"+APIHOST+"/objects/"+"[name]"  |                                                              | Delte object                                                 |
| Get         | "http://"+APIHOST+"/versions/"+"[name]" |                                                              | If name is empty, all version information for all objects is returned. Otherwise, all versions of the corresponding object are returned |
| Get         | "http://"+APIHOST+"/locate/"+"[name]"   |                                                              | Locate the data server on which the object named [name] is located |

## License

Mozilla Public License Version 2.0

## Preface

I write a distributed object storage program based on the book and some modifications of my own. I believe that when the project is completed, I will have a new understanding of the architecture. I hope this is of some use to others and myself.