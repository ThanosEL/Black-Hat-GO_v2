package main

import (
	"fmt"
	"net"
)

func main() {
	for i := 1; i <= 1024; i++ {
		// use the string conversion package, strconv or to use Sprintf
		address := fmt.Sprintf("scanme.nmap.org:%d", i)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			// port is closed or filtered
			continue
		}
		// close the connection if it was successful that way, connections aren’t left open
		conn.Close()
		fmt.Printf("%d open\n", i)
	}
}

// TODO: implement timeout
