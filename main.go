package main

import (
	"flag"
	"strings"
	"time"
	"fmt"
)

func main() {
	port := flag.String("port", ":9091", "rpc listen port")
	cluster := flag.String("cluster", "127.0.0.1:9091", "comma sep")
	id := flag.Int("id", 1, "node ID")

	flag.Parse()
	clusters := strings.Split(*cluster, ",")

	ns := make(map[int]*node)
	for k, v := range clusters {
		ns[k] = newNode(v)
	}

	raft := &Raft{}
	raft.me = *id
	raft.nodes = ns


	raft.rpc(*port)
	callMeDaddy()
	time.Sleep(10 * time.Second)
	raft.start()

	time.Sleep(1 * time.Second)
	str := "Hello"
	clientName := "BSDK"
	for {
		idx := WriteEntry(raft, clientName,  str)
		str += "!"
		time.Sleep(1 * time.Second)
		msgWritten := ReadEntry(raft, 2)

		fmt.Printf("idx %d stringWritten %s\n", idx, msgWritten)

		time.Sleep(500 * time.Millisecond)

	}
}
