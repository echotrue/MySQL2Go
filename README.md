# MySQL2Go
Generate struct for Golang


#### Install

```shell script
go get -u github.com/echotrue/MySQL2Go
```

#### Download
    [Downloading](https://www.baidu.com)
    
    
#### Usage
```shell script
./GoSQL generate -H 192.168.1.1 -U root -P password -D database_name
```

#### Get help
```shell script
$ ./GoSQL generate -h
Generate MySQL tables to Golang struct.

Usage:
  GoSQL generate [flags]

Flags:
  -D, --db string     db name
  -h, --help          help for generate
  -H, --host string   host (default "127.0.0.1")
      --path string   path to be saved struct (default "./struct")
  -p, --port uint16   port (default 3306)
  -P, --pwd string    password
  -U, --user string   user (default "root")

Global Flags:
      --config string   config file (default is $HOME/.MySQL2Go.yaml)

```