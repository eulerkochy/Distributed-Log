package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	port := flag.String("port", ":12345", "rpc raft port")
	mport := flag.String("mport", ":9091", "server listen port")
	cluster := flag.String("cluster", "127.0.0.1:9091", "comma sep")
	id := flag.Int("id", 1001, "node ID")

	flag.Parse()

	fmt.Printf("clusters %v \n", *cluster)

	clusters := strings.Split(*cluster, ",")

	ns := make(map[int]*node)
	for k, v := range clusters {
		ns[k] = newNode(v)
	}

	raft := &Raft{}
	raft.me = *id
	raft.nodes = ns

	PORT := *mport
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

	callMeDaddy()

	raft.rpc(*port)
	time.Sleep(10 * time.Second)
	raft.start()

	time.Sleep(1 * time.Second)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		callMeDaddy()
		clientMsg := strings.TrimSpace(string(netData))

		fmt.Printf("Client Message received : %s\n", clientMsg)

		idx := WriteEntry(raft, "ABC", clientMsg)
		if clientMsg == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Printf("-> %d %s", idx, clientMsg)
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))
	}

}
