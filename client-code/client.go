package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
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
		fmt.Printf("choose [ R - Read Index ] [ W - Write Entry ] [ GET - get all entries by client ] [ GETLOG - get all log entries  ] [ STOP - to stop the server] ")
		fmt.Scanf("%s", &opt)

		var text string

		if opt == "R" {
			fmt.Printf("Enter index >> ")
			fmt.Scanf("%s", &text)
		} else if opt == "W" {
			fmt.Printf("Enter message >> ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan() // use `for scanner.Scan()` to keep reading
			text = scanner.Text()
		} else {
			text = opt
			if !((text == "STOP") || (text == "GET") || (text == "GETLOG")) {
				fmt.Println("Wrong Command entered")
				continue
			}
		}

		for _, v := range clusters {
			CONNECT := v
			c, err := net.Dial("tcp", CONNECT)
			if err != nil {
				fmt.Printf("->: %s   %v\n", v, err)
				continue
			}
			msg := opt + "$" + text + "$" + clientName
			fmt.Fprintf(c, msg+"\n")

			message, _ := bufio.NewReader(c).ReadString('\n')
			fmt.Printf("->: " + v + "  " + message)

		}
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
		// time.Sleep(1 * time.Second)
	}
}
