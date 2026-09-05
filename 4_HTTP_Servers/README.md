#### simple-server.go
The most basic possible HTTP server. Exposes a single route, `/hello`, 
that reads a `name` query parameter and echoes back a greeting.

#### simple-router.go
Builds a custom router from scratch instead of relying on 
`DefaultServeMux`, to show what a router actually does under the hood.

#### middleware.go
Middleware = a wrapper that runs on *every* request before/after the 
real handler, regardless of which route was hit. This example logs the 
start and finish time of each request.

#### gorilla-router/
Replaces Go's built-in `DefaultServeMux` with the third-party `gorilla/mux` 
router, which adds matching on HTTP method, host header, and path 
parameters (with optional regex constraints) — not just the URL path.

#### negroni/
Uses `urfave/negroni` on top of `gorilla/mux` to chain multiple middleware 
together (not just one, like `middleware.go` above). Demonstrates a basic 
auth middleware that checks credentials from the query string, stores the 
authenticated username in the request context, and passes it down the 
chain to the handler. Note: negroni's last release was 2020 and it now 
sees infrequent maintenance, but it still works fine for this example.

#### template.go
Demonstrates Go's `html/template` package, which auto-escapes dynamic 
content based on where it's rendered in the page (HTML body, an href 
attribute, etc.) — built-in XSS protection, unlike the lower-level 
`text/template` package. The `{{.}}` placeholder renders the entire 
value passed to `Execute()`; `{{.FieldName}}` would reach into a specific 
struct field if a struct were passed instead of a plain string.

#### credential-harvester. 
Clones a real login page (e.g. Roundcube webmail), rewrites the form's 
`action` to point at a local `/login` handler instead of the real server, 
and logs whatever credentials get submitted. Serves the cloned static 
files via `http.FileServer` + `gorilla/mux`'s `PathPrefix("/")`, with 
structured logging (`logrus`) writing captured attempts to a file. 
Classic single-factor phishing technique.

#### websocket_keylogger. 
Uses `gorilla/websocket` to upgrade an HTTP connection and stream every 
keystroke from a victim's browser back to the server in real time, via a 
small JS payload (delivered as a Go `html/template` so the WebSocket 
address can be injected dynamically). Relevant in an XSS or compromised-
web-server scenario.

#### c2_mutliplexing
A reverse HTTP proxy (`net/http/httputil` + `gorilla/mux`) that routes 
incoming Meterpreter reverse-HTTP connections to different backend 
listeners based on the `Host` header — the same mechanism virtual hosting 
uses. Lets a single exposed IP/port (e.g. 80/443, likely allowed egress) 
front multiple C2 listeners without exposing them directly, and hints at 
domain fronting as a related technique.

