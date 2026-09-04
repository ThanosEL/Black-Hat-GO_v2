package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HostLocation struct {
	City         string  `json:"city"`
	RegionCode   string  `json:"region_code"`
	AreaCode     int     `json:"area_code"`
	Longitude    float32 `json:"longitude"`
	CountryCode3 string  `json:"country_code3"`
	CountryName  string  `json:"country_name"`
	PostalCode   string  `json:"postal_code"`
	DMACode      int     `json:"dma_code"`
	CountryCode  string  `json:"country_code"`
	Latitude     float32 `json:"latitude"`
}

type Host struct {
	OS        string       `json:"os"`
	Timestamp string       `json:"timestamp"`
	ISP       string       `json:"isp"`
	ASN       string       `json:"asn"`
	Hostnames []string     `json:"hostnames"`
	Location  HostLocation `json:"location"`
	IP        int64        `json:"ip"`
	Domains   []string     `json:"domains"`
	Org       string       `json:"org"`
	Data      string       `json:"data"`
	Port      int          `json:"port"`
	IPString  string       `json:"ip_str"`
}

type HostSearch struct {
	Matches []Host `json:"matches"`
}

// take the search query string as a parameter
func (s *Client) HostSearch(q string) (*HostSearch, error) {
	res, err := http.Get(fmt.Sprintf("%s/shodan/host/search?key=%s&query=%s", BaseURL, s.apiKey, q))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ret HostSearch
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

type Service struct {
	Port      int      `json:"port"`
	Transport string   `json:"transport"`
	Data      string   `json:"data"`
	Product   string   `json:"product"`
	Version   string   `json:"version"`
	CPE       []string `json:"cpe"`
	Timestamp string   `json:"timestamp"`
	Hostnames []string `json:"hostnames"`
}

type HostIPInfo struct {
	IP         int64        `json:"ip"`
	IPString   string       `json:"ip_str"`
	OS         string       `json:"os"`
	Org        string       `json:"org"`
	ISP        string       `json:"isp"`
	ASN        string       `json:"asn"`
	Hostnames  []string     `json:"hostnames"`
	Ports      []int        `json:"ports"`
	Tags       []string     `json:"tags"`
	Vulns      []string     `json:"vulns"`
	LastUpdate string       `json:"last_update"`
	Location   HostLocation `json:"location"`
	Data       []Service    `json:"data"`
}

// look up details for a single IP address
func (s *Client) HostIP(ip string) (*HostIPInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/shodan/host/%s?key=%s", BaseURL, ip, s.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ret HostIPInfo
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
