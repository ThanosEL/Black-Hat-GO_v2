package shodan

// The Shodan URL is defined as a constant value
const BaseURL = "https://api.shodan.io"

// Define a Client struct, used for maintaining your API token across requests
type Client struct {
	apiKey string
}

// taking the API token as input and creating and returning an initialized Client instance
func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}
