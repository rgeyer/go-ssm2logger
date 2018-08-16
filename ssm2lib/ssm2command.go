package ssm2lib

type Ssm2Command byte

const (
	Ssm2CommandNone                    Ssm2Command = 0
	Ssm2CommandReadBlockRequestA0      Ssm2Command = 0xa0
	Ssm2CommandReadBlockResponseE0     Ssm2Command = 0xe0
	Ssm2CommandReadAddressesRequestA8  Ssm2Command = 0xa8
	Ssm2CommandReadAddressesResponseE8 Ssm2Command = 0xe8
	Ssm2CommandWriteBlockRequestB0     Ssm2Command = 0xb0
	Ssm2CommandWriteBlockResponseF0    Ssm2Command = 0xf0
	Ssm2CommandWriteAddressRequestB8   Ssm2Command = 0xb8
	Ssm2CommandWriteAddressResponseF8  Ssm2Command = 0xf8
	Ssm2CommandInitRequestBF           Ssm2Command = 0xbf
	Ssm2CommandInitResponseFF          Ssm2Command = 0xff
)
