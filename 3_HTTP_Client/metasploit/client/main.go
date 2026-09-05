package main

import (
	"fmt"
	"log"
	"os"

	"metasploit/rpc"
)

func main() {
	// Read connection details from environment variables, so credentials
	// never get hardcoded into the source file.
	host := os.Getenv("MSFHOST")
	pass := os.Getenv("MSFPASS")
	user := "msf"

	if host == "" || pass == "" {
		log.Fatalln("Missing required environment variable MSFHOST or MSFPASS")
	}

	// New() connects AND logs in in one step (see rpc/msf.go). If this
	// succeeds, msf already has a valid auth token ready to use.
	msf, err := rpc.New(host, user, pass)
	if err != nil {
		log.Panicln(err)
	}

	// defer schedules Logout() to run automatically when main() exits,
	// whether it exits normally or via a panic — so we always clean up
	// the session token instead of leaving it dangling on the server.
	defer msf.Logout()

	// Fetch the current list of active Meterpreter/shell sessions.
	sessions, err := msf.SessionList()
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("Sessions:")
	for _, session := range sessions {
		fmt.Printf("%5d %s\n", session.ID, session.Info)
	}
}
