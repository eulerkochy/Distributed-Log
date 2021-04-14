# Use goreman to run `go get github.com/mattn/goreman`
test1: ./simple-raft --id 1 --cluster 127.0.0.1:22379,127.0.0.1:32379,127.0.0.1:42379 --port :12379 --mport :9876
test2: ./simple-raft --id 2 --cluster 127.0.0.1:12379,127.0.0.1:32379,127.0.0.1:42379 --port :22379 --mport :9877
test3: ./simple-raft --id 3 --cluster 127.0.0.1:12379,127.0.0.1:22379,127.0.0.1:42379 --port :32379 --mport :9878
test4: ./simple-raft --id 4 --cluster 127.0.0.1:12379,127.0.0.1:22379,127.0.0.1:32379 --port :42379 --mport :9879
client : ./client --cluster 127.0.0.1:9876,127.0.0.1:9877,127.0.0.1:9878,127.0.0.1:9879