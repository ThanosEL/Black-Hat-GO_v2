package main

import (
	"fmt"
	"net"
)

func main() {
	for i := 1; i <= 1024; i++ {
		// to runs concurrently wrap the call to Dial(network, address string) in a goroutine
		go func(j int) {
			address := fmt.Sprintf("scanme.nmap.org:%d", i)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn.Close()
			fmt.Printf("open %d\n", i)
		}(i)
	}
}

// Upon running this code, you should observe the program exiting almost immediately.
// Check with time ./concurent-scanning
