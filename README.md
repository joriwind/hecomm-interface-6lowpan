# The fog interface to the 6LoWPAN network
The 6lowpan network is setup using rpl-border-router as gateway, udp-slip as fog slip and node-cose who will use the hecomm communication system. The udp-slip interacts with the 6lowpan network like a regular node from the 6LoWPAN network and this udp-slip node will transport all UDP packets via SLIP protocol to connected device. 

This golang implementation is a interface for go to communicate with udp-slip. The go interface will distinct between debug information passed via SLIP and actual UDP packages. The debug information is always marked by first '\r' character.