package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	//  sync.WaitGroup acts as a synchronized counter
	var wg sync.WaitGroup
	for i := 1; i <= 1024; i++ {
		// increment this counter via wg.Add(1) each time you create a goroutine to scan a port
		wg.Add(1)
		go func(j int) {
			// a deferred call to wg.Done() decrements the counter whenever one unit of work has been performed
			defer wg.Done()
			address := fmt.Sprintf("scanme.nmap.org:%d", i)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn.Close()
			fmt.Printf("open %d\n", i)
		}(i)
	}
	// Your main() function calls wg.Wait(), which blocks until all the work has been done and your counter has returned to zero
	wg.Wait()
}
