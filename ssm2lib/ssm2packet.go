package ssm2lib

import "fmt"

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
	Destination Ssm2Device
	Source      Ssm2Device
	Command     Ssm2Command
	DataSize    int
	Capacity    int
	buffer      []byte
}

func NewPacketFromBytes(bytes []byte) *Ssm2Packet {
	return &Ssm2Packet{
		Destination: Ssm2Device(bytes[Ssm2PacketIndexDestination]),
		Source:      Ssm2Device(bytes[Ssm2PacketIndexSource]),
		Command:     Ssm2Command(bytes[Ssm2PacketIndexCommand]),
		DataSize:    int(bytes[Ssm2PacketIndexDataSize]),
		buffer:      bytes,
	}
}

func NewInitPacket(src Ssm2Device, dest Ssm2Device) *Ssm2Packet {
	return &Ssm2Packet{
		Destination: dest,
		Source:      src,
		Command:     Ssm2CommandInitRequestBF,
		DataSize:    1,
		buffer:      make([]byte, Ssm2PacketHeaderSize+1),
	}
}

func NewReadPacket(src Ssm2Device, dest Ssm2Device, pids []byte) *Ssm2Packet {
	// TODO:
	// As a result of the maximum data size of 255 bytes, there can be a total of
	// 85 pids in a single read request. 255 / 3 = 85.
	// Need to check for this and return an error or something if more than 85
	// are requested
	buffer_size := 5 + 2 + (3 * len(pids))
	data_size := 2 + (3 * len(pids))

	fmt.Println("Buffer Size: ", buffer_size)
	fmt.Println("Data Size: ", data_size)
	p := &Ssm2Packet{
		Destination: dest,
		Source:      src,
		Command:     Ssm2CommandReadAddressesRequestA8,
		DataSize:    data_size,
		buffer:      make([]byte, buffer_size),
		Capacity:    buffer_size,
	}
	p.buffer[5] = 0x00 // Padding.. I guess

	pids_idx := 6
	for _, pid := range pids {
		pids_idx += 1
		p.buffer[pids_idx] = 0x00 // A blank value for PID1
		pids_idx += 1
		p.buffer[pids_idx] = 0x00 // A blank value for PID1
		pids_idx += 1
		p.buffer[pids_idx] = pid // PID1
	}

	return p
}

func (p *Ssm2Packet) checksumCalculated() byte {
	var sum int
	sum = 0
	length := len(p.buffer) - 1
	for i := 0; i < length; i++ {
		sum = sum + int(p.buffer[i])
	}
	return byte(sum)
}

func (p *Ssm2Packet) GetBytes() []byte {
	p.buffer[0] = Ssm2PacketFirstByte
	p.buffer[1] = byte(p.Destination)
	p.buffer[2] = byte(p.Source)
	p.buffer[3] = byte(p.DataSize)
	p.buffer[4] = byte(p.Command)
	p.buffer[len(p.buffer)-1] = p.checksumCalculated()
	return p.buffer
}

func (p *Ssm2Packet) GetData() []byte {
	return p.buffer[Ssm2PacketHeaderSize-1 : Ssm2PacketHeaderSize+p.DataSize-1]
}
