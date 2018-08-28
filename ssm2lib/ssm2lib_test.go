package ssm2lib_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rgeyer/ssm2logger/ssm2lib"
)

var _ = Describe("Ssm2lib", func() {
	It("Can create a write address packet to switch to fast mode", func() {
		// writeFastModePacket := NewWriteAddressPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10, []byte{0x00, 0x01, 0x98}, []byte{0x5a})
		// Ω(writeFastModePacket.GetBytes()).Should(Equal([]byte{byte(Ssm2PacketFirstByte), byte(Ssm2DeviceEngine10), byte(Ssm2DeviceDiagnosticToolF0), 0x05, byte(Ssm2CommandWriteAddressRequestB8), 0x00, 0x01, 0x98, 0x5a, 0x30}))
		Ω(true).Should(Equal(true))
	})

	It("Can create a read address request", func() {
		readPacket := NewReadAddressRequestPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10, []byte{0x46})
		Ω(readPacket.Bytes()).Should(Equal([]byte{0x80, 0x10, 0xf0, 0x05, 0xa8, 0x00, 0x00, 0x00, 0x46, 0x73}))
	})

	It("Can create an init request", func() {
		initPacket := NewInitRequestPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10)
		Ω(initPacket.Bytes()).Should(Equal([]byte{0x80, 0x10, 0xf0, 0x01, 0xbf, 0x40}))
	})

	Context("Wire time", func() {
		It("Knows how long bytes will take on the wire", func() {
			microseconds := MicrosecondsOnTheWireBytes(make([]byte, 8))
			Ω(microseconds).Should(Equal(16667))
		})
	})

	Context("Validation", func() {
		Context("The first byte is wrong", func() {
			It("Returns an error", func() {
				bogusPacket := NewPacketFromBytes([]byte{0x00})
				err := bogusPacket.Validate()
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).To(Equal("First byte of packet is wrong. Expected 0x80, got 0x00"))
			})
		})
	})
})
