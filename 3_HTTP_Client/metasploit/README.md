# Interacting with Metasploit
Metasploit is a framework used to perform a variety of adversarial techniques,
including reconnaissance, exploitation, command and control, persistence,
lateral network movement, payload creation and delivery, privilege escalation, and more.

Build a client that interacts with a remote Metasploit instance.
Directory structure:
```
|---client
| |---main.go
|---rpc
|---msf.go
```

### Setting Up Your Environment
1. Start the Metasploit RPC daemon on the lab VM:
```bash
   msfrpcd -P yourpassword -U msf -a 127.0.0.1 -p 55553 -S
```
2. Set environment variables for the client to authenticate:
```bash
   export MSFHOST="127.0.0.1:55553"
   export MSFPASS="yourpassword"
```
3. Install the MessagePack dependency:
```bash
   go get github.com/vmihailenco/msgpack/v5
```
