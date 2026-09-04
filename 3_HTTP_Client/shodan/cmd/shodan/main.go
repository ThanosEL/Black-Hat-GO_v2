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
	hostInfo, err := s.HostIP(os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Org: %s\nISP: %s\nOS: %s\nHostnames: %v\nOpen Ports: %v\nTags: %v\nVulns: %v\nLast Update: %s\nCity: %s\nCountry: %s (%s)\nLat/Long: %.4f, %.4f\n\n",
		hostInfo.Org, hostInfo.ISP, hostInfo.OS, hostInfo.Hostnames,
		hostInfo.Ports, hostInfo.Tags, hostInfo.Vulns, hostInfo.LastUpdate,
		hostInfo.Location.City, hostInfo.Location.CountryName, hostInfo.Location.CountryCode,
		hostInfo.Location.Latitude, hostInfo.Location.Longitude)
	for _, svc := range hostInfo.Data {
		fmt.Printf("Port: %d (%s) — %s\n", svc.Port, svc.Transport, svc.Product)
	}
}
