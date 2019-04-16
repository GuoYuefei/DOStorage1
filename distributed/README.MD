# Distributed Object Storage

- [Getting Started](#getting-started)
  + [Environmental configuration](#Environmental-configuration)
  + [Usage](#usage)

- [Lisence](#License)
- [Preface](#Preface)

## Getting Started

### Environmental configuration

|Variable name| interface | data |
|:---:|:---------:|:----:|
|RABBITMQ_SERVER|need|need|
|LISTEN_ADDRESS|need|need|
|STORAGE_ROOT|no need|need|

-----
```powershell
$env:RABBITMQ_SERVER = "amqp://test:test@192.168.1.68:5672/"
$env:LISTEN_ADDRESS = "192.168.1.68:23333"
$env:STORAGE_ROOT = YOU_PATH(Where you want to store the object)
./interface.exe or ./data.exe 
```

In rabbitmqctl command
```powershell
rabbitmqctl add_user test test
rabbitmqctl set_permissions -o / test ".*" ".*" ".*"
```

Then <code>go run declareExchange.go</code> in Rabbit Message Queue Server machine.

Use <code>rabbitmqctl list_exchanges</code> to view the results. If there are two exchanges, apiServers and dataServers, the above program has been successfully executed.

### Usage

Please **create a directory of STORAGE_ROOT**  in advance and **create an folders named *objects* ** in that directory.

The configuration content mentioned above needs to be configured on each node machine as needed.

Then you can find the corresponding binary files data.exe and interface.exe in the release. Run the two programs on different node machines through the command line.

You can compile it yourself with the source code.

## License

Mozilla Public License Version 2.0

## Preface

I write a distributed object storage program based on the book and some modifications of my own. I believe that when the project is completed, I will have a new understanding of the architecture. I hope this is of some use to others and myself.