<h1 align="center">
  Distributed-Log  
</h1>

Distributed Log for Distributed Systems course project. 

  + A Raft-based implementation of distributed log where **multiple** clients can write and read entries in the log.  
  + Log is replicated and persisted among the Raft clusters.
  + On receiving a `STOP` command from client, Leader saves its log.


**Installation and Usage**
---
```
$ git clone https://github.com/eulerkochy/Distributed-Log
$ cd Distributed-Log

# build the server code
$ go build -o simple-log

# run the server code
$ ./simple-log --id <clusterID> --cluster <hostname:port>[,<hostname:port>,..]  --port <raft_port> --mport <message_port>

# build the client code
$ cd client-code
$ go build -o client

# run the client code
$ ./client --cluster  <hostname:mport>[,<hostname:mport>,..]

Usage:
  - [W]     	Write entries in log, Leader returns the log index
  - [R]     	Read entry by index, Leader returns the logged message
  - [GET]   	Get all entries by the particular client
  - [GETLOG]	Get all entries in the log 
  - [STOP]  	Stop the servers and client

```

**Demonstration**
---

+ on server private IP `172.26.0.49`
```
./simple-log --id 1 --cluster 172.26.9.127:12345,172.26.6.197:12345,172.26.9.243:12345 --port :12345 --mport :9876
```

+ on server private IP `172.26.6.197`

```
./simple-log --id 2 --cluster 172.26.0.49:12345,172.26.9.127:12345,172.26.9.243:12345 --port :12345 --mport :9876

```

+ on server private IP `172.26.9.127`

```
./simple-log --id 3 --cluster 172.26.0.49:12345,172.26.6.197:12345,172.26.9.243:12345 --port :12345 --mport :9876

```

+ on server private IP `172.26.9.243`

```
./simple-log --id 4 --cluster 172.26.0.49:12345,172.26.6.197:12345,172.26.9.127:12345 --port :12345 --mport :9876

```

+ on clients

```
./client --cluster 172.26.0.49:9876,172.26.6.197:9876,172.26.9.127:9876,172.26.9.243:9876
```

