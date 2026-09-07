## Environment Setup
Requires `libpcap-dev` (Linux) for `gopacket` to interface with the OS 
packet capture stack:
```bash
sudo apt-get install libpcap-dev
go get github.com/google/gopacket
go get github.com/google/gopacket/pcap@v1.1.19
```

#### identify-devices/
Enumerates available network interfaces and their IP/netmask info using 
`gopacket/pcap`. No elevated privileges needed — just enumeration, not 
capture. Useful first step before pointing a capture tool at the right 
interface.

#### filter-capture/
Opens a live capture handle on a specific interface and applies a BPF 
filter (`tcp and port 80`) so only matching traffic is processed — the 
same filter syntax used by `tcpdump`/Wireshark. Requires root/sudo since 
it reads raw packets off the wire. Prints each captured packet's decoded 
layers (Ethernet, IPv4, TCP) to stdout.

#### ftp-sniffer/ 
Extends the BPF-filtered packet capture above, but filters specifically 
for port 21 (FTP) traffic and inspects the application-layer payload for 
`USER`/`PASS` commands — capturing cleartext login credentials sent over 
unencrypted FTP sessions. 

#### syn-flood-scanner/
Port scanner to work correctly against targets 
protected by SYN cookies (SYN-flood protection), which otherwise make 
every port look "open" during a standard TCP handshake. 

Runs a packet sniffer concurrently while attempting real TCP connections 
to each target port, filtering for TCP flag combinations (ACK+FIN, ACK, 
or ACK+PSH — `tcp[13] == 0x11/0x10/0x18` in BPF syntax) that only occur 
if an actual service exchanges data beyond the initial handshake. Ports 
that never produce this follow-up traffic are treated as false positives 
from SYN cookie obfuscation, not real open ports.

```bash
sudo go run main.go <interface> <target_ip> <port1,port2,port3>
```