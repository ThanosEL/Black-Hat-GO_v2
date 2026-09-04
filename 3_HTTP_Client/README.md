## HTTP Fundamentals with Go
1. HTTP is a stateless protocol: the server doesn’t inherently maintain state and status for each request. Instead, state is tracked through a variety of means, which may include session identifiers, cookies, HTTP headers, and more.

2. Communications between clients and servers can occur either synchronously or asynchronously, but they operate on a request/response cycle.

3. APIs commonly communicate via more structured data encoding, such as XML, JSON, or MSGRPC.

## Calling HTTP APIs
*3 basic functions:*
- http.Get(url)
- http.Head(url)
- http.Post(url, contentType, body)

Always close the responcse body:
```go
defer r1.Body.Close()
```

*POST shortcut:*
```go
// Better this:
http.PostForm(url, form)
// Than this: http.Post() with manual encoding
```

# HTTP_Client
Examples of building and sending HTTP requests in Go using `net/http`,
and handling structured responses.

#### basic-requests.go
Demonstrates GET, HEAD, POST, PostForm, DELETE, and PUT requests.
Convenience functions (`http.Get`, `http.Post`) cover the common verbs,
while `http.NewRequest` + `client.Do()` handles the rest (DELETE, PUT, PATCH).
Always `defer resp.Body.Close()` to avoid memory leaks.

#### response-parsing.go
Reads and prints the HTTP status code and response body from a GET request
to `google.com/robots.txt`. Uses `io.ReadAll()` to read from the `resp.Body`
`io.ReadCloser`.

#### json-parsing.go
Decodes a JSON response from an API endpoint into a Go struct using
`json.NewDecoder(res.Body).Decode(&struct)`. Requires the struct fields
to match the JSON keys.

#### simple-api-server.go
Local test server on `:8080` that serves `/ping` — returns a JSON response
`{"Message":"...","Status":"..."}` to test `json-parsing.go` against a
real endpoint.