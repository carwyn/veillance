# Setup Access Point With NetworkManager CLI

```
nmcli c add type wifi ifname '*' con-name veillance autoconnect no ssid Veillance
nmcli c mod veillance 802-11-wireless.mode ap 802-11-wireless.band bg ipv4.method shared
nmcli c mod veillance wifi-sec.key-mgmt wpa-psk
nmcli c mod veillance wifi-sec.psk "veillance"
nmcli c up veillance
```

# Links

* https://vincent.bernat.im/en/blog/2014-intel-7260-access-point
* https://www.hogarthuk.com/?q=node/8
* https://fedoraproject.org/wiki/Networking/CLI
* https://docs.fedoraproject.org/en-US/Fedora/25/html/Networking_Guide/index.html

