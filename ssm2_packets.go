package main

import "fmt"

type Ssm2Packet struct {
	Destination byte
	Source      byte
	Command     byte
	Size        int
	Capacity    int
	buffer      []byte
}

func NewInitPacket(src byte, dest byte) *Ssm2Packet {
	return &Ssm2Packet{
		Destination: dest,
		Source:      src,
		Command:     0xBF,
		Size:        1,
		buffer:      make([]byte, 6),
		Capacity:    6,
	}
}

func NewReadPacket(src byte, dest byte, pids []byte) *Ssm2Packet {
	buffer_size := 5 + 2 + (3 * len(pids))
	data_size := 2 + (3 * len(pids))

	fmt.Println("Buffer Size: ", buffer_size)
	fmt.Println("Data Size: ", data_size)
	p := &Ssm2Packet{
		Destination: dest,
		Source:      src,
		Command:     0xA8,
		Size:        data_size,
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
	p.buffer[0] = 0x80
	p.buffer[1] = p.Destination
	p.buffer[2] = p.Source
	p.buffer[3] = byte(p.Size)
	p.buffer[4] = p.Command
	p.buffer[len(p.buffer)-1] = p.checksumCalculated()
	return p.buffer
}
