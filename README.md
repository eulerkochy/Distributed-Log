# Distributed-Log
Distributed log for Distributed Systems course project

**Installation & Usage**
---
```
$ git clone https://github.com/eulerkochy/Distributed-Log
$ cd Distributed-Log

# build the server code
$ go build -o simple-log

# run the server code
$ ./simple-log --id <clusterID> --cluster <hostname:port>[,<hostname:port>,..]  --port <raft_port> --mport <message_port>

# build the client code
$ cd client-codes
$ go build -o client

# run the client code
$ ./client --cluster  <hostname:mport>[,<hostname:mport>,..]

Usage:
  - [W]     Write entries in log, Leader returns the log index
  - [R]     Read entry by index, Leader return the logged message
  - [GET]   Get all entries by the particular client
  - [STOP]  Stop the servers and client

```

**Demo**
---

+ on server private IP `172.26.0.49`
```
./simple-raft --id 1 --cluster 172.26.9.127:12345,172.26.6.197:12345,172.26.9.243:12345 --port :12345 --mport :9876
```

+ on server private IP `172.26.6.197`

```
./simple-raft --id 2 --cluster 172.26.0.49:12345,172.26.9.127:12345,172.26.9.243:12345 --port :12345 --mport :9876

```

+ on server private IP `172.26.9.127`

```
./simple-raft --id 3 --cluster 172.26.0.49:12345,172.26.6.197:12345,172.26.9.243:12345 --port :12345 --mport :9876

```

+ on server private IP `172.26.9.243`

```
./simple-raft --id 4 --cluster 172.26.0.49:12345,172.26.6.197:12345,172.26.9.127:12345 --port :12345 --mport :9876

```

+ on clients

```
./client --cluster 172.26.0.49:9876,172.26.6.197:9876,172.26.9.127:9876,172.26.9.243:9876
```

