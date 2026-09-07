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

