package ssm2lib

import "fmt"

type Ssm2InitResponsePacket struct {
	Packet Ssm2PacketBytes
}

func NewSsm2InitResponsePacketFromBytes(packet Ssm2PacketBytes) (*Ssm2InitResponsePacket, error) {
	if packet.GetCommand() != Ssm2CommandInitResponseFF {
		return nil, fmt.Errorf("Can not construct an Ssm2InitResponsePacket from supplied bytes. The command in the packet should be %s, but received %s", Ssm2CommandInitResponseFF.String(), packet.GetCommand().String())
	}

	return &Ssm2InitResponsePacket{Packet: packet}, nil
}

func (p *Ssm2InitResponsePacket) GetSsmId() []byte {
	return p.Packet[Ssm2PacketHeaderSize : Ssm2PacketHeaderSize+3]
}

func (p *Ssm2InitResponsePacket) GetRomId() []byte {
	romIdIndex := Ssm2PacketHeaderSize + 3
	return p.Packet[romIdIndex : romIdIndex+5]
}

func (p *Ssm2InitResponsePacket) GetCapabilityBytes() []byte {
	capabilitiesIndex := Ssm2PacketHeaderSize + 3 + 5
	return p.Packet[capabilitiesIndex : len(p.Packet)-1]
}
