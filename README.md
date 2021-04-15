# Distributed-Log
Distributed log for Distributed Systems course project

### build
```
go build -o simple-raft
```

### run
```
go get github.com/mattn/goreman
goreman start
```

### server terminal codes

- on server private ip 172.26.0.49

```
./simple-raft --id 1 --cluster 172.26.9.127:12345,172.26.6.197:12345,172.26.9.243:12345 --port :12345 --mport :9876

```

- on server private ip 172.26.6.197

```
./simple-raft --id 2 --cluster 172.26.0.49:12345,172.26.9.127:12345,172.26.9.243:12345 --port :12345 --mport :9876

```

- on server private ip 172.26.9.127

```
./simple-raft --id 3 --cluster 172.26.0.49:12345,172.26.6.197:12345,172.26.9.243:12345 --port :12345 --mport :9876

```

- on server private ip 172.26.9.243

```
./simple-raft --id 4 --cluster 172.26.0.49:12345,172.26.6.197:12345,172.26.9.127:12345 --port :12345 --mport :9876

```

- on clients

```
./client --cluster 172.26.0.49:9876,172.26.6.197:9876,172.26.9.127:9876,172.26.9.243:9876
```
