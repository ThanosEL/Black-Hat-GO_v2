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

