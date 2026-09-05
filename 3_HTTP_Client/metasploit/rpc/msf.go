package rpc

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/vmihailenco/msgpack/v5"
)

// sessionListReq is the request sent to Metasploit's "session.list" RPC method.
// The "asArray" tag tells msgpack to encode this as a positional array
// (e.g. ["session.list", "mytoken"]) instead of a map/object, because
// that's the wire format Metasploit's RPC API expects for requests.
type sessionListReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Token    string
}

// SessionListRes represents a single Meterpreter session, as returned by
// "session.list". Unlike the request, this IS decoded from a map (session
// details are key/value pairs), so each field needs a tag naming the exact
// key it maps to.
type SessionListRes struct {
	ID          uint32 `msgpack:",omitempty"` // populated manually later, not sent by Metasploit itself (see SessionList())
	Type        string `msgpack:"type"`
	TunnelLocal string `msgpack:"tunnel_local"`
	TunnelPeer  string `msgpack:"tunnel_peer"`
	ViaExploit  string `msgpack:"via_exploit"`
	ViaPayload  string `msgpack:"via_payload"`
	Description string `msgpack:"desc"`
	Info        string `msgpack:"info"`
	Workspace   string `msgpack:"workspace"`
	SessionHost string `msgpack:"session_host"`
	SessionPort int    `msgpack:"session_port"`
	Username    string `msgpack:"username"`
	UUID        string `msgpack:"uuid"`
	ExploitUUID string `msgpack:"exploit_uuid"`
}

// loginReq is the request for Metasploit's "auth.login" method.
// Same idea as sessionListReq: encoded as a positional array
// ["auth.login", "username", "password"].
type loginReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Username string
	Password string
}

// loginRes is the response from "auth.login". Go only fills in the fields
// that are actually present in the response, so this single struct works
// for BOTH a successful login (Result + Token populated) and a failed one
// (Error, ErrorClass, ErrorMessage populated instead).
type loginRes struct {
	Result       string `msgpack:"result"`
	Token        string `msgpack:"token"`
	Error        bool   `msgpack:"error"`
	ErrorClass   string `msgpack:"error_class"`
	ErrorMessage string `msgpack:"error_message"`
}

// logoutReq is the request for "auth.logout". Takes the method name, the
// current token, and a logout token (kept simple here — we just reuse
// the same token for both).
type logoutReq struct {
	_msgpack    struct{} `msgpack:",asArray"`
	Method      string
	Token       string
	LogoutToken string
}

// logoutRes is the response from "auth.logout" — just a simple success flag.
type logoutRes struct {
	Result string `msgpack:"result"`
}

// Metasploit holds everything needed to talk to a remote Metasploit RPC
// server: where it is (host), how to authenticate (user/pass), and the
// active session token once logged in. Building methods on this struct
// means we don't have to keep passing host/credentials into every call.
type Metasploit struct {
	host  string
	user  string
	pass  string
	token string
}

// send is the shared plumbing used by every RPC call. It takes ANY request
// struct and ANY response struct (via interface{}), so we only have to
// write the HTTP + encoding/decoding logic once instead of repeating it
// in Login(), Logout(), SessionList(), etc.
func (msf *Metasploit) send(req interface{}, res interface{}) error {
	// Encode the request struct into MessagePack binary format into buf.
	buf := new(bytes.Buffer)
	msgpack.NewEncoder(buf).Encode(req)

	// Build the target URL using this client's configured host.
	dest := fmt.Sprintf("http://%s/api", msf.host)

	// POST the binary-encoded request. Metasploit's RPC API expects the
	// content type "binary/message-pack" so it knows how to parse the body.
	r, err := http.Post(dest, "binary/message-pack", buf)
	if err != nil {
		return err
	}
	defer r.Body.Close() // always close the response body when done

	// Decode the binary response directly into whatever "res" points to.
	// This is why callers pass a pointer (&res) — so this function can
	// write the decoded data back into the caller's variable.
	if err := msgpack.NewDecoder(r.Body).Decode(&res); err != nil {
		return err
	}
	return nil
}

// Login authenticates against the Metasploit RPC server using auth.login,
// and stores the returned token on the struct for use in later calls.
func (msf *Metasploit) Login() error {
	ctx := &loginReq{
		Method:   "auth.login",
		Username: msf.user,
		Password: msf.pass,
	}
	var res loginRes
	if err := msf.send(ctx, &res); err != nil {
		return err
	}
	msf.token = res.Token // save token for future authenticated calls
	return nil
}

// Logout invalidates the current session token via auth.logout, then
// clears it locally so this client can no longer make authenticated calls
// until it logs in again.
func (msf *Metasploit) Logout() error {
	ctx := &logoutReq{
		Method:      "auth.logout",
		Token:       msf.token,
		LogoutToken: msf.token,
	}
	var res logoutRes
	if err := msf.send(ctx, &res); err != nil {
		return err
	}
	msf.token = "" // token is no longer valid, clear it
	return nil
}

// SessionList calls session.list to retrieve all active Meterpreter/shell
// sessions. Requires a valid token (i.e. Login() must have succeeded first).
func (msf *Metasploit) SessionList() (map[uint32]SessionListRes, error) {
	req := &sessionListReq{Method: "session.list", Token: msf.token}

	// Metasploit returns sessions as a map: session ID (number) -> session details.
	res := make(map[uint32]SessionListRes)
	if err := msf.send(req, &res); err != nil {
		return nil, err
	}

	// The session ID only exists as the MAP KEY in the raw response, not
	// inside the SessionListRes struct itself. This loop copies that key
	// into the .ID field on each struct, so callers can access it directly
	// (session.ID) instead of having to track map keys separately.
	for id, session := range res {
		session.ID = id
		res[id] = session
	}
	return res, nil
}

// New bootstraps a Metasploit client: stores connection details and
// immediately attempts to log in. If login fails, it returns an error
// instead of a half-working client — so any Metasploit struct that
// successfully comes out of New() is guaranteed to already have a valid token.
func New(host, user, pass string) (*Metasploit, error) {
	msf := &Metasploit{
		host: host,
		user: user,
		pass: pass,
	}
	if err := msf.Login(); err != nil {
		return nil, err
	}
	return msf, nil
}
