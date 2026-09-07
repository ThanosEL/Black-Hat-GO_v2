# SMB & NTLM Protocol Reference
## What is SMB?

**SMB (Server Message Block)** is an application-layer network protocol 
used primarily by Windows systems to share access to files, printers, 
named pipes, and other resources over a network. Conceptually it plays a 
similar role to HTTP (client sends a request, server responds), but the 
similarities stop there — SMB is a **binary** protocol, not a 
human-readable text protocol. Messages are built from a mix of 
fixed-length fields, variable-length fields, and positional data laid out 
in a specific byte order (little-endian), rather than key-value headers 
you can read in a terminal.

### SMB Dialects (Versions)
SMB has evolved through several versions, called "dialects." A client and 
server must agree on a dialect before they can communicate, and they do 
this negotiation as the very first step of any connection.

| Dialect | Introduced with |
|---|---|
| SMB 1.0 | Legacy Windows / early networking (now considered insecure, largely deprecated) |
| SMB 2.0 | Windows Vista / Server 2008 |
| SMB 2.1 | Windows 7 / Server 2008 R2 |
| SMB 3.0 | Windows 8 / Server 2012 |
| SMB 3.02 | Windows 8.1 / Server 2012 R2 |
| SMB 3.1.1 | Windows 10 / Server 2016+ |

Each newer dialect adds performance and security improvements. A modern 
Windows client and server will always negotiate up to the highest dialect 
both sides support.

## How Authentication Works: NTLMSSP
SMB itself doesn't dictate *how* users are authenticated — that's handled 
by a separate negotiated mechanism. In Active Directory environments, the 
most common one is **NTLMSSP** (NTLM Security Support Provider), a 
challenge-response authentication scheme.

### Why challenge-response instead of just sending a password?
Instead of transmitting a password (or even its hash) directly over the 
wire, the server sends the client a random **challenge** value. The 
client must prove it knows the correct password by combining that 
challenge with a value derived from the password (the NTLM hash) and 
sending back a computed **response**. The server — or an authoritative 
system like a domain controller — independently performs the same 
calculation and checks whether the two match. If the client didn't 
actually possess the correct hash, it can't produce the response.

### The Handshake, Step by Step
Client Server
|------ Negotiate Protocol (dialects) -------->|
|<----- Chosen dialect + auth options ---------|
|------ Session Setup (NTLMSSP Negotiate) ----->|
|<----- Session Setup (challenge token) --------|
| [client computes NTLM hash + response] |
|------ Session Setup (NTLMSSP Authenticate) -->|
|<----- Session Setup (success + session ID) ---|
|------ Requests using session ID ------------->|


1. **Negotiate Protocol** — client tells the server which SMB dialects it 
   understands.
2. **Server response** — picks the best mutually supported dialect, and 
   lists which authentication mechanisms it accepts.
3. **Session Setup (Negotiate)** — client initiates NTLMSSP.
4. **Session Setup (Challenge)** — server sends back a random challenge 
   value, signaling more steps are required.
5. **Client-side computation** — the client derives the NTLM hash from 
   the user's domain, username, and password, then combines it with the 
   server's challenge and a client-generated random value to produce a 
   response value.
6. **Session Setup (Authenticate)** — client sends that computed response 
   back to the server.
7. **Server validation** — the server (or a domain controller on its 
   behalf) performs the same computation independently and compares 
   results. A match means successful authentication; the server issues a 
   session identifier.
8. **Authenticated requests** — every subsequent SMB request (browsing a 
   share, reading a file, etc.) includes that session ID so the server 
   knows the client is already authenticated.

### The Encoding Complication: ASN.1 Inside Positional Binary
SMB's general message format is positional binary (fields in a fixed 
order, fixed or length-prefixed sizes). NTLMSSP tokens, however, travel 
wrapped inside a **GSS-API** container that uses **ASN.1** encoding — an 
entirely different binary encoding standard, commonly seen in 
cryptographic and telecom protocols. This means a single SMB message can 
contain two different encoding schemes side by side: the outer SMB 
structure is positional binary, while the embedded authentication token 
is ASN.1. Any implementation has to correctly switch parsing strategy 
mid-message.

## Why This Matters Operationally
- **NTLM authenticates based on a hash, not the plaintext password.** 
  This is the root cause of an entire category of attacks (like 
  pass-the-hash) where possessing the hash alone — without ever knowing 
  the actual password — is enough to authenticate.
- **The challenge-response design has historically had implementation 
  weaknesses** (particularly around challenge randomness and relaying 
  authentication attempts to a third party) that enable NTLM relay 
  attacks.
- **NTLM hashes, once captured, can potentially be cracked offline** — 
  their cryptographic strength depends heavily on password complexity, 
  which is why weak domain passwords remain a major real-world risk in 
  Windows environments.
- **SMB signing** (a feature that cryptographically signs each message) 
  is one of the primary defenses against relay-style attacks — 
  understanding the handshake explains exactly why signing closes that gap.

## Key Specifications (for reference)
- **MS-SMB2** — the SMB2 protocol specification itself
- **MS-SPNG / RFC 4178** — the GSS-API negotiation wrapper (ASN.1 encoded) 
  that carries authentication tokens
- **MS-NLMP** — defines the NTLMSSP token structure and the exact 
  challenge-response calculations (not ASN.1 encoded — this inner layer 
  uses its own binary format)
- **ASN.1 (X.690)** — the general-purpose encoding standard used by the 
  GSS-API layer

  