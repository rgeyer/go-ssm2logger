package ssm2lib

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

// TODO: Add some goodies here for showing which parameters are supported. Here's
// an example of decoding the init response.
// fmt.Println(resp_bytes[8] & (1 << 6)) A test of looking for a specific parameter using bitwise operators. Gotta move this elsewhere.
type Ssm2Connection struct {
	serial     io.ReadWriteCloser
	buf_serial *bufio.ReadWriter
	logger     *log.Entry
	buffer     []byte
}

// I wasn't smart enough to figure out the timing myself, I got that answer here
// https://www.microchip.com/forums/m405110.aspx
func MicrosecondsOnTheWireBytes(buffer []byte) int {
	// bit_time = 1 / baud * databytes * (1 start bit, 8 bit word, 1 stop bit = 10)
	return int(math.Round(1.0 / 4800.00 * 1000000.0 * float64(len(buffer)) * 10.0))
}

func MicrosecondsOnTheWireByteCount(count int) int {
	// bit_time = 1 / baud * databytes * (1 start bit, 8 bit word, 1 stop bit = 10)
	return int(math.Round(1.0 / 4800.00 * 1000000.0 * float64(count) * 10.0))
}

func (c *Ssm2Connection) SetLogger(logger *log.Logger) {
	c.logger = logger.WithFields(log.Fields{"logger": "Ssm2Connection"})
}

func (c *Ssm2Connection) Open(port string) error {
	config := &serial.Config{
		Name:        port,
		Baud:        4800,
		StopBits:    serial.Stop1,
		Parity:      serial.ParityNone,
		ReadTimeout: 1 * time.Second,
	}

	serial_port, err := serial.OpenPort(config)
	if err != nil {
		return fmt.Errorf("Error opening serial port: %s", err)
	}
	c.serial = serial_port
	c.buf_serial = bufio.NewReadWriter(bufio.NewReader(serial_port), bufio.NewWriter(serial_port))

	return nil
}

func (c *Ssm2Connection) Close() {
	c.serial.Close()
}

func (c *Ssm2Connection) InitEngine() (*Ssm2InitResponsePacket, error) {
	initPacket := NewInitRequestPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10)
	packetBytes, err := c.sendPacketAndFetchResponsePacket(initPacket.Packet)
	if err != nil {
		return nil, err
	}
	return NewSsm2InitResponsePacketFromBytes(packetBytes)
}

/// <summary>
/// Maximum number of addresses to be used for a single packet.
/// Defaults to 36 which most control units from year 2002+ should support.
/// Possible theoretical range: 1 ≤ x ≤ 84.
/// Control unit will not answer when using to many addresses!
/// (Tested ok on modern cars (2008+): ≤ 82, packet size is 253 bytes for 82)
/// (2005 cars might support ≤ 45.)
/// (84 is theoretical limit because of packet length byte)
/// </summary>
func (c *Ssm2Connection) ReadAddresses(addresses [][]byte) (Ssm2PacketBytes, error) {
	readPacket := NewReadAddressRequestPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10, addresses, false)
	return c.sendPacketAndFetchResponsePacket(readPacket.Packet)
}

func (c *Ssm2Connection) ReadParameters(params []Ssm2Parameter) (Ssm2PacketBytes, error) {
	packet_size := Ssm2PacketHeaderSize + 1 + (3 * len(params)) + 1
	buffer := make([]byte, packet_size)
	buffer[0] = Ssm2PacketFirstByte
	buffer[Ssm2PacketIndexDestination] = byte(Ssm2DeviceEngine10)
	buffer[Ssm2PacketIndexSource] = byte(Ssm2DeviceDiagnosticToolF0)
	// This is a hack, since the size must include the command byte, and exclude
	// the checksum byte. TODO: Not sure if this is right. Have to wrap my head
	// around it better.
	buffer[Ssm2PacketIndexDataSize] = byte(packet_size - Ssm2PacketHeaderSize)
	buffer[Ssm2PacketIndexCommand] = byte(Ssm2CommandReadAddressesRequestA8)
	// 0x00 for single request 0x01 for continuous
	buffer[Ssm2PacketIndexData] = 0x01

	param_idx := Ssm2PacketIndexData + 1
	for _, param := range params {
		address, err := param.Address.GetAddressBytes()
		if err != nil {
			return nil, err
		}
		buffer[param_idx] = address[0]
		param_idx += 1
		buffer[param_idx] = address[1]
		param_idx += 1
		buffer[param_idx] = address[2]
		param_idx += 1
	}

	buffer[len(buffer)-1] = CalculateChecksum(buffer)
	return c.sendPacketAndFetchResponsePacket(buffer)
}

func (c *Ssm2Connection) ReadAddressesContinous(addresses [][]byte) (Ssm2PacketBytes, error) {
	readPacket := NewReadAddressRequestPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10, addresses, true)
	return c.sendPacketAndFetchResponsePacket(readPacket.Packet)
}

func (c *Ssm2Connection) sendPacketAndFetchResponsePacket(packet Ssm2PacketBytes) (Ssm2PacketBytes, error) {
	if c.logger != nil {
		c.logger.WithFields(log.Fields{"command": packet.GetCommand(), "bytes": hex.EncodeToString(packet)}).Debug("Sending SSM2 Command")
	}

	wrotebytes, err := c.serial.Write(packet)
	if err != nil {
		return nil, fmt.Errorf("Failed to send serial: %s", err)
	}

	// TODO: Check if we successfully wrote all the bytes!

	if c.logger != nil {
		c.logger.WithFields(log.Fields{"wrote_bytes": wrotebytes}).Debug("Packet sent, time to try to fetch it")
	}

	// Should be able to grab the sent packet right away, so no wait.
	// TODO: Turns out, I'm not getting all the bytes on this read, and probably
	// subsequent reads!
	// DEBU[0000] Sending SSM2 Command                          bytes=8010f005a80000004673 command="Read Address Request" logger=Ssm2Connection
	// DEBU[0000] Packet sent, time to try to fetch it          logger=Ssm2Connection wrote_bytes=10
	// DEBU[0000] Read the write                                count=5 data=8010f005a80000000000 error="<nil>" logger=Ssm2Connection
	// DEBU[0000] Read the header, to the command byte          count=5 data=0000004673 error="<nil>" logger=Ssm2Connection
	// DEBU[0000] Read the remaining bytes, to the checksum     count=7 data=0000004673 error="<nil>" logger=Ssm2Connection
	// {"Destination":"None","Source":"None","Command":"0x73","DataSize":70,"Data":"c4DwEALogOoAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="} <nil>
	// request_bytes := make([]byte, wrotebytes)
	// count, err := c.serial.Read(request_bytes)
	// if c.logger != nil {
	// 	c.logger.WithFields(log.Fields{"count": count, "error": err, "data": hex.EncodeToString(request_bytes)}).Debug("Read the write")
	// }

	// Is this totally unecessary? Sleeping between putting bits on the wire, and
	// reading them. This should be "instantaneous" since the file is local in
	// both cases?
	time.Sleep(time.Duration(MicrosecondsOnTheWireBytes(packet)) * time.Microsecond)

	_, err = c.GetNextPacketInStream()
	if err != nil {
		return nil, err
	}

	responsePacket, err := c.GetNextPacketInStream()

	return responsePacket, err
}

func (c *Ssm2Connection) GetNextPacketInStream() (Ssm2PacketBytes, error) {
	var retval Ssm2PacketBytes
	// Always assume that we have a fully formed packet waiting for us to fetch.
	// Also assume we're trying to fetch a response packet immediately after
	// sending a request and wait for the ECU/TCU to put a packet header on the wire
	time.Sleep(time.Duration(MicrosecondsOnTheWireByteCount(Ssm2PacketHeaderSize)) * time.Microsecond)
	// Grab the header (up to the command byte)
	header_bytes := make([]byte, Ssm2PacketHeaderSize)
	err := c.ensureSerialRead(&header_bytes)
	//count, err := c.serial.Read(header_bytes)
	if c.logger != nil {
		c.logger.WithFields(log.Fields{
			//"count": count,
			"error": err,
			"data":  hex.EncodeToString(header_bytes),
		}).Debug("Read the header, to the command byte")
	}

	// Grab the remainder of the packet. Using the DataSize byte value, which will
	// include all data, plus the checksum byte
	remaining_bytes := make([]byte, header_bytes[Ssm2PacketIndexDataSize])
	// Wait for the remaining packets to be put on the wire
	time.Sleep(time.Duration(MicrosecondsOnTheWireBytes(remaining_bytes)) * time.Microsecond)
	//count, err = c.serial.Read(remaining_bytes)
	err = c.ensureSerialRead(&remaining_bytes)
	if c.logger != nil {
		c.logger.WithFields(log.Fields{
			//"count": count,
			"error": err,
			"data":  hex.EncodeToString(remaining_bytes),
		}).Debug("Read the remaining bytes, to the checksum")
	}
	packet_bytes := append(header_bytes, remaining_bytes...)
	if c.logger != nil {
		c.logger.WithFields(log.Fields{
			// "count": count,
			"error": err,
			"data":  hex.EncodeToString(packet_bytes),
		}).Debug("Here's the entire packet")
	}
	retval = Ssm2PacketBytes(packet_bytes)
	return retval, nil
}

func (c *Ssm2Connection) ensureSerialRead(desiredBuffer *[]byte) error {
	desired_buf_len := len(*desiredBuffer)
	first_read := make([]byte, desired_buf_len)
	first_count, err := c.serial.Read(first_read)
	// Check to see if we're out pacing the protocol
	if first_count < len(*desiredBuffer) {
		remaining_bytes_to_read := desired_buf_len - first_count
		second_read := make([]byte, remaining_bytes_to_read)
		time_in_microseconds := MicrosecondsOnTheWireByteCount(remaining_bytes_to_read)
		if c.logger != nil {
			c.logger.WithFields(log.Fields{"wait": time_in_microseconds, "expected_count": desired_buf_len, "read_count": first_count, "error": err}).Debug("Didn't fill the read buffer, throttling and retrying precisely once")
		}
		time.Sleep(time.Duration(time_in_microseconds) * time.Microsecond)
		count, err := c.serial.Read(second_read)
		if err != nil {
			return err
		}
		if first_count+count < len(*desiredBuffer) {
			return fmt.Errorf("Couldn't fill supplied read buffer. Expected %d bytes, but could only get %d after one retry with a %d microsecond cooldown", desired_buf_len, first_count+count, time_in_microseconds)
		}
		*desiredBuffer = append(first_read[:first_count], second_read...)
	} else {
		*desiredBuffer = first_read
	}
	return nil
}

/*
func (c *Ssm2Connection) safeRead(buf []byte) (int, error) {
	buf_buf := make([]byte, len(buf))
	// TODO: This is nasty, and I need to have a context which I can use to interrupt
	for {
		count, err := c.serial.Read(buf_buf)
		if count == 0 {
			if c.logger != nil {
				c.logger.WithFields(log.Fields{"expected length": len(buf), "temp buf length": len(buf_buf), "read length": count, "error": err}).Debug("Looping on 0 byte read")
			}
			continue
		} else {
			buf = buf_buf
		}
		if count < len(buf) {
			if c.logger != nil {
				c.logger.WithFields(log.Fields{"expected length": len(buf), "read length": count}).Warn("Buffer empty before all expected bytes were read.")
			}
		}
		if err != nil && err.Error() == "EOF" {
			// Try to wait exactly the requisite time for the remaining bytes.
			remaining_bytes_to_read := len(buf) - count
			time_in_microseconds := 208 * remaining_bytes_to_read
			time.Sleep(time.Duration(time_in_microseconds) * time.Microsecond)

			if c.logger != nil {
				c.logger.WithFields(log.Fields{"remaining bytes": remaining_bytes_to_read, "microseconds": time_in_microseconds}).Debug("Hit EOF reading serial, trying again, once after a sane wait duration.")
			}

			newbuf := make([]byte, remaining_bytes_to_read)
			newcount, err := c.serial.Read(newbuf)
			if err != nil {
				return count + newcount, err
			}
			buf = append(buf_buf, newbuf...)
		}
		return count, nil
	}
}*/
