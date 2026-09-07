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


#### smb-password-guessing/
Attempts SMB/NTLM authentication against a target using a list of 
candidate usernames paired with a single password, checking which (if any) 
succeed — a technique used to test for weak or default credentials across a domain.

`Note:` Online password guessing can lock accounts out of a domain, effectively resulting in
a denial-of-service attack. Take caution when testing your code and run this against
only systems on which you’re authorized to test

```
go run main.go users.txt "Password123!" LAB.LOCAL 192.168.1.100
```

#### pass-the-hash/ 
Attempts SMB authentication using a captured NTLM password hash directly, 
without ever knowing (or needing) the plaintext password. Works because 
NTLM authentication separates hash calculation from challenge-response 
token calculation — the response computation only needs the hash as 
input, not the domain/username/password that originally produced it. So 
any precomputed hash (e.g. extracted from a compromised machine's memory) 
is sufficient on its own to authenticate elsewhere.

Iterates a list of target hosts, attempting authentication with the same 
username + hash pair against each one — testing for credential/hash 
reuse across a network, a common step in lateral movement during a 
domain compromise.

```
go run main.go targets.txt Administrator LAB.LOCAL aad3b435b51404eeaad3b435b51404ee:31d6cfe0d16ae931b73c59d7e0c089c0
```

*How NTLM Hashes Are Obtained (Overview)*

Understanding where hashes come from in the first place helps explain why 
pass-the-hash is such a persistent problem — hashes end up in a lot of 
places, some more heavily defended than others.

##### 1. Memory extraction (LSASS)
On a running Windows system, credentials used to log in are cached in 
memory by the **LSASS** (Local Security Authority Subsystem Service) 
process, in order to support single sign-on. An attacker with local 
administrator access on a compromised machine can dump LSASS memory and 
extract NTLM hashes (and sometimes even plaintext passwords, depending on 
Windows configuration) for every account that has logged in since boot. 
This is why gaining local admin on even one machine is often treated as a 
serious event — it can expose credentials for other, more privileged 
accounts that happened to log in there (e.g. IT staff doing support).

##### 2. Local SAM database
Every non-domain-controller Windows machine stores local account 
password hashes in the **SAM** (Security Account Manager) database. An 
attacker with sufficient local privileges can extract these hashes, 
though they only cover local accounts, not domain accounts.

##### 3. Domain controller database (NTDS.dit)
Domain controllers store hashes for every domain account in a database 
file called **NTDS.dit**. Compromising a domain controller (or obtaining 
a valid backup of one) exposes the hashes for the entire domain at once — 
this is one of the primary objectives of a full domain compromise.

##### 4. Network capture and relay (no compromise required)
NTLM authentication attempts sometimes occur automatically over a 
network without explicit user action — for example, when a misconfigured 
client tries to resolve a hostname via legacy broadcast protocols 
(NBNS/LLMNR) instead of DNS. An attacker positioned on the same network 
segment can respond to these broadcast requests, tricking the victim 
machine into sending its NTLM challenge-response directly to the 
attacker. Depending on configuration, the attacker can either:
- Capture the exchange and attempt to crack the hash offline, or
- **Relay** the authentication attempt in real time to a different target 
  server, authenticating as the victim without ever needing the hash or 
  password at all.

This category is particularly notable because it doesn't require any 
prior foothold or vulnerability — just network visibility and a 
misconfiguration that's common in default Windows environments.

##### 5. Credential dumping from configuration/scripts
Hashes and even plaintext credentials sometimes end up stored in 
scripts, configuration management tools, scheduled tasks, or Group 
Policy Preferences (a legacy Windows feature with a well-known history 
of storing recoverable credentials). These are typically found through 
general file/system enumeration rather than a dedicated "attack" per se.