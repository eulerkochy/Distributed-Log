package main

import (
	"bufio"
	"fmt"
	"net"

	// "os"
	"flag"
	"strings"
	// "time"
)

func main() {
	cluster := flag.String("cluster", "127.0.0.1:9091", "comma sep")
	flag.Parse()
	clusters := strings.Split(*cluster, ",")
	fmt.Printf("%v \n", clusters)

	for {
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print(">> ")
		// text := "hi"
		var text string
		fmt.Scanf("%s", &text)
		for _, v := range clusters {
			CONNECT := v
			c, err := net.Dial("tcp", CONNECT)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Fprintf(c, text+"\n")

			message, _ := bufio.NewReader(c).ReadString('\n')
			fmt.Print("->: " + message)
			// if strings.TrimSpace(string(text)) == "STOP" {
			//         fmt.Println("TCP client exiting...")
			//         return
			// }
		}
		// time.Sleep(1 * time.Second)
	}
}
