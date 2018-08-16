package ssm2lib

import "fmt"

type Ssm2Device byte

const (
	Ssm2DeviceNone             Ssm2Device = 0
	Ssm2DeviceEngine10         Ssm2Device = 0x10
	Ssm2DeviceTransmission18   Ssm2Device = 0x18
	Ssm2DeviceDiagnosticToolF0 Ssm2Device = 0xf0
)

func (d Ssm2Device) String() string {
	switch device := d; device {
	case Ssm2DeviceNone:
		return "None"
	case Ssm2DeviceEngine10:
		return "Engine"
	case Ssm2DeviceTransmission18:
		return "Transmission"
	case Ssm2DeviceDiagnosticToolF0:
		return "DiagnosticTool"
	default:
		return "Unknown"
	}
}

func (d Ssm2Device) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}
