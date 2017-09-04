package sixlowpan

import (
	"bytes"
	"encoding/binary"
	"net"

	"golang.org/x/net/ipv6"
)

//UDPHeader Defines the header of UDP packet
type UDPHeader struct {
	SrcPort uint16
	DstPort uint16
	Length  uint16
	Chksum  uint16
	Payload []byte
}

type udppseudoHeader struct {
	SrcAddress net.IP
	DstAddress net.IP
	UDPLength  uint32
	Zeroes     uint8
	NxtHdr     uint8
}

//Marschal Compile the UDPHeader into byte slice
func (h UDPHeader) Marschal() ([]byte, error) {
	//We know the length of the entire UDPheader + payload, h.Lenght should be the comprise lenght of hdr & payload
	//Size the internal write buffer to h.Length via buf & setting the capacity (=h.Length) but not the length (=0)
	buf := make([]byte, 0, h.Length)
	b := bytes.NewBuffer(buf)

	//Checking extensively just because maybe reason
	var err error
	err = binary.Write(b, binary.BigEndian, h.SrcPort)
	if err != nil {
		return nil, err
	}
	err = binary.Write(b, binary.BigEndian, h.DstPort)
	if err != nil {
		return nil, err
	}
	err = binary.Write(b, binary.BigEndian, h.Length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(b, binary.BigEndian, h.Chksum)
	if err != nil {
		return nil, err
	}
	err = binary.Write(b, binary.BigEndian, h.Payload)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

//UnmarshalUDP Unmarshal the byte slice into UDP header
func UnmarshalUDP(buf []byte) (h UDPHeader, err error) {
	b := bytes.NewBuffer(buf)
	err = binary.Read(b, binary.BigEndian, h.SrcPort)
	if err != nil {
		return h, err
	}
	err = binary.Read(b, binary.BigEndian, h.DstPort)
	if err != nil {
		return h, err
	}
	err = binary.Read(b, binary.BigEndian, h.Length)
	if err != nil {
		return h, err
	}
	err = binary.Read(b, binary.BigEndian, h.Chksum)
	if err != nil {
		return h, err
	}
	err = binary.Read(b, binary.BigEndian, h.Payload)
	if err != nil {
		return h, err
	}
	return h, err
}

//CalcChecksum Calculate the checksum of UDP header, providing the ip header information for pseudoheader
func (h *UDPHeader) CalcChecksum(ip *ipv6.Header) error {
	h.Chksum = 0
	phdr := udppseudoHeader{
		SrcAddress: ip.Src.To16(),
		DstAddress: ip.Dst.To16(),
		UDPLength:  uint32(h.Length),
		Zeroes:     0,
		NxtHdr:     17,
	}
	var b bytes.Buffer
	var err error
	//Writing first part of pseudo header
	err = binary.Write(&b, binary.BigEndian, phdr.SrcAddress)
	if err != nil {
		return err
	}
	err = binary.Write(&b, binary.BigEndian, phdr.DstAddress)
	if err != nil {
		return err
	}
	err = binary.Write(&b, binary.BigEndian, phdr.UDPLength)
	if err != nil {
		return err
	}
	for i := 0; i < 3; i++ { //Total of 24bits of data
		err = binary.Write(&b, binary.BigEndian, phdr.Zeroes)
		if err != nil {
			return err
		}
	}
	err = binary.Write(&b, binary.BigEndian, phdr.NxtHdr)
	if err != nil {
		return err
	}

	//Writing second part, also marschal of UDP header
	bytes, err := h.Marschal()
	if err != nil {
		return err
	}
	err = binary.Write(&b, binary.BigEndian, bytes)
	if err != nil {
		return err
	}

	//Overwrite checksum in provided UDP header
	h.Chksum = checksum(b.Bytes())
	return nil
}

//Checksum Implementation from: https://gist.github.com/chrisnc/0ff3d1c20cb6687454b0 -- rawudp.go
func checksum(buf []byte) uint16 {
	sum := uint32(0)

	for ; len(buf) >= 2; buf = buf[2:] {
		sum += uint32(buf[0])<<8 | uint32(buf[1])
	}
	if len(buf) > 0 {
		sum += uint32(buf[0]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	csum := ^uint16(sum)
	/*
	 * From RFC 768:
	 * If the computed checksum is zero, it is transmitted as all ones (the
	 * equivalent in one's complement arithmetic). An all zero transmitted
	 * checksum value means that the transmitter generated no checksum (for
	 * debugging or for higher level protocols that don't care).
	 */
	if csum == 0 {
		csum = 0xffff
	}
	return csum
}
