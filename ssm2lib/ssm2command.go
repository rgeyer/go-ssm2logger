package ssm2lib

import "fmt"

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

func (c Ssm2Command) String() string {
	switch command := c; command {
	case Ssm2CommandNone:
		return "None"
	case Ssm2CommandReadBlockRequestA0:
		return "Read Block Request"
	case Ssm2CommandReadBlockResponseE0:
		return "Read Block Response"
	case Ssm2CommandReadAddressesRequestA8:
		return "Read Address Request"
	case Ssm2CommandReadAddressesResponseE8:
		return "Read Address Response"
	case Ssm2CommandWriteBlockRequestB0:
		return "Write Block Request"
	case Ssm2CommandWriteBlockResponseF0:
		return "Write Block Response"
	case Ssm2CommandWriteAddressRequestB8:
		return "Write Address Request"
	case Ssm2CommandWriteAddressResponseF8:
		return "Write Address Response"
	case Ssm2CommandInitRequestBF:
		return "Init Request"
	case Ssm2CommandInitResponseFF:
		return "Init Response"
	default:
		return fmt.Sprintf("0x%.2x", byte(c))
	}
}

func (c Ssm2Command) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", c.String())), nil
}
