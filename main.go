package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
	"strconv"
	"os"
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
	time.Sleep(1 * time.Second)
	raft.start()

	time.Sleep(1 * time.Second)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		msg := strings.TrimSpace(string(netData))

		clientMsgArr := strings.Split(msg, "$")

		opt, clientMsg, clientName := clientMsgArr[0], clientMsgArr[1], clientMsgArr[2]

		var idx int

		t := time.Now()

		if opt == "W" {
			idx = WriteEntry(raft, clientName, clientMsg)
			if (idx == -1) {
				myMsg := t.Format(time.RFC3339) + " :: error: server is not a Leader" + "\n"
				c.Write([]byte(myMsg))
			} else {
				myMsg := t.Format(time.RFC3339) + " :: " + fmt.Sprintf("stored at log index %d", idx) +"\n"
				c.Write([]byte(myMsg))
			}
		} else if opt == "R" {
			idx, _ = strconv.Atoi(clientMsg)
			logMsg := ReadEntry(raft, idx)
			myMsg := t.Format(time.RFC3339) + " :: " + fmt.Sprintf("log at index %d is %s", idx, logMsg) +"\n"
			c.Write([]byte(myMsg))

		} else if opt == "GET" { // get all entries by that particular client
			logMsgs := GetEntries(raft, clientName)
			myMsg := t.Format(time.RFC3339) + " :: " + fmt.Sprintf("log entries by %s are %s", clientName, logMsgs) +"\n"
			c.Write([]byte(myMsg))			
		} else if opt == "GETLOG" {
			logMsgs := GetAllEntries(raft)
			myMsg := t.Format(time.RFC3339) + " :: " + fmt.Sprintf("all log entries : %s", logMsgs) +"\n"
			c.Write([]byte(myMsg))
		} else if opt == "STOP" { 
			// save the log 
			logArr := GetAllEntriesArray(raft)
			if len(logArr) > 0 {
				filename := "log-id-" + strconv.Itoa(*id) + "-"+ t.Format(time.RFC3339) + ".txt"
				f, err := os.Create(filename)
				if err != nil {
			        fmt.Println("error: cannot create " + filename)
			    }

			    defer f.Close()

			    for _, log := range logArr {

			        _, err := f.WriteString(log + "\n")

			        if err != nil {
			        	fmt.Println("error: cannot Write " + log)
			        }
			    }
			}


			fmt.Println("Stopping TCP server!")
			myMsg := "server stopped"
			c.Write([]byte(myMsg))
			return
		} else {
			myMsg := "error :: wrong command"
			c.Write([]byte(myMsg))
		}
	}

}
