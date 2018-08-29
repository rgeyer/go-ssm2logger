package ssm2lib

import (
	"encoding/json"
	"fmt"
)

type Ssm2PacketIndex int

const (
	/// <summary>
	/// This first byte is always supposed to be 128 = 0x80.
	/// </summary>
	Ssm2PacketIndexHeader Ssm2PacketIndex = 0
	/// <summary>
	/// Destination device.
	/// Typically 0x10 = engine, 0x18 = transmission, 0xF0 = diagnostic tool.
	/// </summary>
	Ssm2PacketIndexDestination Ssm2PacketIndex = 1
	/// <summary>
	/// Source device. Same principle as Destination.
	/// </summary>
	Ssm2PacketIndexSource Ssm2PacketIndex = 2
	/// <summary>
	/// Inline payload length, counting all following bytes except checksum.
	/// </summary>
	Ssm2PacketIndexDataSize Ssm2PacketIndex = 3
	Ssm2PacketIndexCommand  Ssm2PacketIndex = 4
	/// <summary>
	/// Generic data, length varies.
	/// Some packet types use a padding byte as first data byte.
	/// </summary>
	Ssm2PacketIndexData Ssm2PacketIndex = 5
)

const (
	Ssm2PacketMaxSize    int  = 4 + 255 + 1
	Ssm2PacketHeaderSize int  = 5
	Ssm2PacketMinSize    int  = Ssm2PacketHeaderSize + 1
	Ssm2PacketFirstByte  byte = 0x80
)

type Ssm2Packet struct {
	Packet Ssm2PacketBytes
}

type Ssm2PacketBytes []byte

func NewPacketBytes(dest Ssm2Device, src Ssm2Device, command Ssm2Command, data []byte) Ssm2PacketBytes {
	retval := make([]byte, Ssm2PacketMinSize+len(data))
	retval[0] = Ssm2PacketFirstByte
	retval[Ssm2PacketIndexDestination] = byte(dest)
	retval[Ssm2PacketIndexSource] = byte(src)
	retval[Ssm2PacketIndexDataSize] = byte(len(data) + 1)
	retval[Ssm2PacketIndexCommand] = byte(command)
	retval[len(retval)-1] = CalculateChecksum(retval)
	return Ssm2PacketBytes(retval)
}

func (b Ssm2PacketBytes) GetFirstByte() byte {
	return b[0]
}

func (b Ssm2PacketBytes) GetDataSize() int {
	return int(b[Ssm2PacketIndexDataSize])
}

func (b Ssm2PacketBytes) GetData() []byte {
	return b[Ssm2PacketIndexData : int(Ssm2PacketIndexData)+b.GetDataSize()]
}

func (b Ssm2PacketBytes) GetCommand() Ssm2Command {
	return Ssm2Command(b[Ssm2PacketIndexCommand])
}

func (b Ssm2PacketBytes) Validate() error {
	if b.GetFirstByte() != Ssm2PacketFirstByte {
		return fmt.Errorf("First byte of packet is wrong. Expected 0x80, got 0x%.2x", b.GetFirstByte())
	}
	return nil
}

func (p *Ssm2Packet) ToJson() (string, error) {
	js, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(js), nil
}

func NewPacketFromBytes(bytes []byte) *Ssm2Packet {
	p := &Ssm2Packet{
		Packet: bytes,
	}
	return p
}

func NewInitRequestPacket(src Ssm2Device, dest Ssm2Device) *Ssm2Packet {
	buffer := NewPacketBytes(dest, src, Ssm2CommandInitRequestBF, nil)
	return &Ssm2Packet{
		Packet: buffer,
	}
}

func NewReadAddressRequestPacket(src Ssm2Device, dest Ssm2Device, pids []byte, push_mode bool) *Ssm2Packet {
	// TODO:
	// As a result of the maximum data size of 255 bytes, there can be a total of
	// 83 pids in a single read request. (255-5) / 3 = 83.
	// Need to check for this and return an error or something if more than 83
	// are requested

	// Header (5 bytes)
	// + Request Type (0x00 for single request 0x01 for continuous)
	// + 3 bytes for each address + 1 checksum byte
	packet_size := Ssm2PacketHeaderSize + 1 + (3 * len(pids)) + 1
	buffer := make([]byte, packet_size)
	buffer[0] = Ssm2PacketFirstByte
	buffer[Ssm2PacketIndexDestination] = byte(dest)
	buffer[Ssm2PacketIndexSource] = byte(src)
	// This is a hack, since the size must include the command byte, and exclude
	// the checksum byte. TODO: Not sure if this is right. Have to wrap my head
	// around it better.
	buffer[Ssm2PacketIndexDataSize] = byte(packet_size - Ssm2PacketHeaderSize)
	buffer[Ssm2PacketIndexCommand] = byte(Ssm2CommandReadAddressesRequestA8)
	// 0x00 for single request 0x01 for continuous
	if push_mode {
		buffer[Ssm2PacketIndexData] = 0x01
	} else {
		buffer[Ssm2PacketIndexData] = 0x00
	}

	pids_idx := Ssm2PacketIndexData + 1
	for _, pid := range pids {
		buffer[pids_idx] = 0x00 // A blank value for PID1
		pids_idx += 1
		buffer[pids_idx] = 0x00 // A blank value for PID1
		pids_idx += 1
		buffer[pids_idx] = pid // PID1
		pids_idx += 1
	}

	buffer[len(buffer)-1] = CalculateChecksum(buffer)
	p := &Ssm2Packet{
		Packet: buffer,
	}
	return p
}

func CalculateChecksum(buffer []byte) byte {
	var sum int
	sum = 0
	length := len(buffer) - 1
	for i := 0; i < length; i++ {
		sum = sum + int(buffer[i])
	}
	return byte(sum)
}
