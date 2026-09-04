package main

import (
	"fmt"
	"log"
	"os"
	"shodan/shodan"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: shodan searchterm")
	}
	// Start by reading your API key from the SHODAN_API_KEY environment variable
	apiKey := os.Getenv("SHODAN_API_KEY")
	// Then use that value to initialize a new Client struct
	s := shodan.New(apiKey)
	// subsequently using it to call your APIInfo() method
	info, err := s.APIInfo()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf(
		"Query Credits: %d\nScan Credits: %d\n\n",
		info.QueryCredits,
		info.ScanCredits)

	// passing in a search string captured as a command line argument
	hostSearch, err := s.HostSearch(os.Args[1])
	if err != nil {
		log.Panic(err)
	}

	// loop through the results to display the IP and port
	for _, host := range hostSearch.Matches {
		fmt.Printf("%18s%8d\n", host.IPString, host.Port)
	}
}
