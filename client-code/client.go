package main

import (
	"bufio"
	"fmt"
	"net"
	"flag"
	"strings"
)

func main() {
	cluster := flag.String("cluster", "127.0.0.1:9091", "comma sep")
	flag.Parse()
	clusters := strings.Split(*cluster, ",")
	fmt.Printf("%v \n", clusters)

	fmt.Printf("Enter clientName >> ")
	var clientName string
	fmt.Scanf("%s", &clientName)

	for {
		var opt string
		fmt.Printf("choose [ R - Read Index ] [ W - Write Entry ] [GET - get all entries] [ STOP - to stop the server] ")
		fmt.Scanf("%s", &opt)

		var text string
		
		if opt == "R" {
			fmt.Printf("Enter index >> ")
			fmt.Scanf("%s", &text)
		} else if opt == "W" {
			fmt.Printf("Enter message >> ")
			fmt.Scanf("%s", &text)
		} else {
			text = opt
		}

		for _, v := range clusters {
			CONNECT := v
			c, err := net.Dial("tcp", CONNECT)
			if err != nil {
				fmt.Println(err)
				continue
			}
			msg := opt + "$" + text + "$" + clientName
			fmt.Fprintf(c, msg+"\n")

			message, _ := bufio.NewReader(c).ReadString('\n')
			fmt.Print("->: " + message)
			if strings.TrimSpace(string(text)) == "STOP" {
			    fmt.Println("TCP client exiting...")
			    return
			}
		}
		// time.Sleep(1 * time.Second)
	}
}
