package ssm2lib

type Ssm2Device byte

const (
	Ssm2DeviceNone             Ssm2Device = 0
	Ssm2DeviceEngine10         Ssm2Device = 0x10
	Ssm2DeviceTransmission18   Ssm2Device = 0x18
	Ssm2DeviceDiagnosticToolF0 Ssm2Device = 0xf0
)

func (d Ssm2Device) String() string {
	if d == Ssm2DeviceEngine10 {
		return "Engine"
	}

	return "None"
}
