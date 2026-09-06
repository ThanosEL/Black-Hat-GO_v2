# Exploiting DNS

## Lab Setup: BIND9 Primary/Secondary
To make the DNS exercises in this chapter realistic, two BIND9 DNS 
servers were set up as a primary/secondary pair, hosting a fake zone: 
`blackhatlab.local`.

- ubuntu1 (192.168.177.132) — Primary/Master
- ubuntu2 (192.168.177.133) — Secondary/Slave

### Zone contents
| Record | Value |
|---|---|
| `ns1.blackhatlab.local` | 192.168.177.132 |
| `ns2.blackhatlab.local` | 192.168.177.133 |
| `www.blackhatlab.local` | 192.168.177.132 |
| `mail.blackhatlab.local` | 192.168.177.132 |
| `vpn.blackhatlab.local` | 192.168.177.132 |
| `dev.blackhatlab.local` | 192.168.177.132 |

### Primary config (`/etc/bind/named.conf.local` on ubuntu1)
```
zone "blackhatlab.local" {
type master;
file "/etc/bind/db.blackhatlab.local";
allow-transfer { 192.168.177.133; };
};
```

### Secondary config (`/etc/bind/named.conf.local` on ubuntu2)
```
zone "blackhatlab.local" {
type slave;
file "/var/cache/bind/db.blackhatlab.local";
masters { 192.168.177.132; };
};
```

### Verified
- Both servers answer authoritative queries independently (`dig @<ip> blackhatlab.local NS`).
- Zone transfer (AXFR) from primary → secondary confirmed via 
  `journalctl -u bind9` (`Transfer status: success`, 11 records).
- `allow-transfer` restricts AXFR to the secondary's IP only — tested from 
  a third host (Windows host, via WSL `dig`) to confirm unauthorized zone 
  transfers are refused. This mirrors the classic DNS misconfiguration 
  pentesters check for (`allow-transfer { any; }` would let anyone dump 
  the entire zone).


## DNS
#### get-a-record.go
Sends a raw DNS query for an A record (`stacktitan.com` → `blackhatlab.local` 
in this lab) using `github.com/miekg/dns`, and processes the response's 
`Answer` slice via type assertion (`answer.(*dns.A)`) to pull out the actual 
IP address. A records are the basic DNS record type mapping a hostname to 
an IPv4 address — the foundation every other lookup in this chapter builds 
on. Also demonstrates following CNAME chains: if a hostname doesn't have 
a direct A record but instead points to another hostname via CNAME, you 
follow that chain until you reach one that does.

#### subdomain-guesser/
A concurrent subdomain brute-forcer. Reads a wordlist, builds candidate 
FQDNs (`word.domain.com`), and resolves each one against a DNS server 
using a worker-pool pattern (goroutines + channels) — the same pattern 
used for the concurrent port scanner in Chapter 2. Follows CNAME chains 
until it reaches a real A record, so aliased subdomains still resolve 
correctly. Tested against a local BIND9 lab (`blackhatlab.local`) instead 
of a real target.

Usage:
```bash
go run main.go -domain blackhatlab.local -wordlist subdomains.txt -server 192.168.177.132:53 -c 100
```

Scan with 5000 words:
```
-c 1	5.895s
-c 100	0.663s
```