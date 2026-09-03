package main

import (
	"fmt"
	"net"
)

func main() {
	// The first argument is a string that identifies the kind of connection to initiate.
	// The second argument tells Dial(network, address string) the host to which you wish to connect
	_, err := net.Dial("tcp", "scanme.nmap.org:80")
	if err == nil {
		fmt.Println("Connection successful")
	}
}
