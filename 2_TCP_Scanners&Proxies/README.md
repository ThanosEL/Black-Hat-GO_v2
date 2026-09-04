## Understanding the TCP Handshake
The below screenshots shows how TCP uses a handshake process 
when querying a port to determine whether the port is open, closed, or filtered.

#### Open Port
If the port is open, a three-way handshake takes place. First, the client sends a syn packet,
or acknowledgment of the syn packet it received, prompting the client to finish with an ack, or acknowledgment of the servers response.
![alt text](images/image.png)

#### Closed Port
If the port is closed, the server responds with a rst packet instead of a syn-ack
![alt text](images/image2.png)

#### Filterd Port
If the traffic is being filtered by a firewall, the client will typically receive no
response from the server.
![alt text](images/image3.png)


## Bypassing Firewalls with Port Forwarding
People can configure firewalls to prevent a client from connecting to certain
servers and ports, while allowing access to others.
In some cases, you can circumvent these restrictions by using an intermediary 
system to proxy the connection around or through a firewall, a technique known as port forwarding.

1. Imagine a nefarious site called evil.com.
2. If an employee attempts to browse evil.com directly, a firewall blocks the request
3. However, should an employee own an external system thats allowed through the firewall (for example, stacktitan.com)
4. That employee can leverage the allowed domain to bounce connectionsto evil.com
![alt text](images/image4.png)
---


# Scanners
In the `scanners` folder are different kinds of imprementations, each one 
improving on the previous. This progression shows why the naive approach 
isn't enough, and how we get to a correct concurrent scanner.

#### single-port-scanner.go
Dials a single TCP port (`scanme.nmap.org:80`) using `net.Dial`. If `err` 
is `nil`, the connection succeeded. The most basic building block — no 
loop, no concurrency.

#### nonconcurrent-scanning.go
Loops through ports 1–1024 sequentially, dialing each one and closing the 
connection if successful. Works, but slow — each port is scanned one at a 
time, so it takes as long as the sum of all connection attempts.

#### concurent-scanning.go
Wraps each `Dial` call in a goroutine to scan all ports at once. Much 
faster, but broken: `main()` doesn't wait for the goroutines to finish, so 
the program exits almost instantly, before most connections even complete. 
Results are unreliable.

#### synchronized-scanning.go
Fixes the previous version using `sync.WaitGroup`. `wg.Add(1)` before each 
goroutine, `wg.Done()` when it finishes, and `wg.Wait()` in `main()` blocks 
until all scans are done. Correct and concurrent — but scanning too many 
ports at once can still overwhelm the network/system and skew results.

#### workerPool-scanner.go
Introduces a worker pool using a buffered channel and `sync.WaitGroup`. 
100 worker goroutines consume ports via `range`. No actual scanning — 
demonstrates the pool pattern before adding network logic.

#### tcp-scanner-final.go
Complete port scanner using two channels: `ports` (buffered) for work 
distribution and `results` for collecting outcomes. Workers send `0` for 
closed, port number for open. `WaitGroup` is replaced by receiving exactly 
1024 results — same synchronization, no counter. Output is sorted.


# proxies
#### io-example.go
Demonstrates `io.Reader` and `io.Writer` interfaces by wrapping `stdin` and `stdout`
in custom types `FooReader` and `FooWriter`. Shows both manual `Read()`/`Write()` 
calls and the `io.Copy()` shorthand that replaces them.

#### echo-server.go
TCP server on port 20080 that echoes data back to the client. Uses `net.Listen()` 
to bind, `listener.Accept()` to handle incoming connections, and `conn.Read()`/
`conn.Write()` for raw I/O. Each connection runs in a goroutine via `go echo(conn)`.

#### echo-server-improved.go
Refactors the echo handler using `bufio.NewReader`/`bufio.NewWriter` for buffered I/O,
replacing manual byte slice management. Final version simplifies further with 
`io.Copy(conn, conn)` — passing `conn` as both source and destination.

#### tcp-proxy.go
A TCP port forwarder that proxies connections through an intermediary host.
Listens on port 8080 and forwards all traffic to `localhost:9999`, using two
`io.Copy` calls for bidirectional forwarding — one in a goroutine to prevent
blocking.

**Tested locally:**
- `Target` = `nc -lvp 9999` (Terminal 1)
- `Proxy`  = `go run proxy.go` — listens :8080, forwards to :9999 (Terminal 2)
- `Client` = `nc localhost 8080` (Terminal 3)

Confirmed bidirectional traffic: client→proxy→target and target→proxy→client.

#### netcat.go
Replicates Netcat's gaping security hole — a TCP listener on port 13337 that
grants shell access to any connecting client. Uses `io.Pipe()` to synchronously
connect `/bin/sh` stdout to the TCP connection, and assigns `conn` directly to
`cmd.Stdin`. Any command sent by the client is executed on the server and the
output is returned over the connection.

**Tested with:**
- `go run netcat.go` as the listener
- `nc localhost 13337` as the connecting client
- Confirmed remote command execution: `ls` returned server-side directory listing