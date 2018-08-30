package ssm2lib

import (
	"encoding/hex"
	"fmt"

	"github.com/Knetic/govaluate"
)

// TODO: These are all prefixed as "SSM2", but they're really RomRaider definitions
// Need to move this into an appropriate namespace and tweak things a bit.

type Ssm2Logger struct {
	Version   string         `xml:"version,attr"`
	Protocols []Ssm2Protocol `xml:"protocols>protocol"`
}

type Ssm2Protocol struct {
	Id             string          `xml:"id,attr"`
	Baud           int             `xml:"baud,attr"`
	DataBits       int             `xml:"databits,attr"`
	StopBits       int             `xml:"stopbits,attr"`
	Parity         int             `xml:"parity,attr"`
	ConnectTimeout int             `xml:"connect_timeout,attr"`
	SendTimeout    int             `xml:"send_timeout,attr"`
	Parameters     []Ssm2Parameter `xml:"parameters>parameter"`
}

// TODO: Some of these have dependencies on other params, and are actually
// derived values
type Ssm2Parameter struct {
	Id           string                    `xml:"id,attr"`
	Name         string                    `xml:"name,attr"`
	Description  string                    `xml:"desc,attr"`
	EcuByteIndex uint                      `xml:"ecubyteindex,attr"`
	EcuBit       uint                      `xml:"ecubit,attr"`
	Target       uint                      `xml:"target,attr"`
	Address      Ssm2ParameterAddress      `xml:"address"`
	Conversions  []Ssm2ParameterConversion `xml:"conversions>conversion"`
}

func (p Ssm2Parameter) Convert(unit string, value []byte) (float64, error) {
	var intval int
	// TODO: I'm making several assumptions here. I've only tested with 1 byte
	// responses so far, and I'm not 100% sure what the 2+ byte responses are or
	// how they work.
	if len(value) > 1 {
		intval = int(uint(value[1]) | uint(value[0])<<8)
	} else if len(value) == 1 {
		intval = int(value[0])
	} else {
		intval = 0
	}
	for _, conversion := range p.Conversions {
		if conversion.Units == unit {
			params := make(map[string]interface{}, 1)
			params["x"] = intval
			expr, err := govaluate.NewEvaluableExpression(conversion.Expr)
			if err != nil {
				return 0, err
			}

			result, err := expr.Evaluate(params)
			if err != nil {
				return 0, err
			}
			return result.(float64), nil
		}
	}
	return 0, fmt.Errorf("Unable to find a converstion with unit (%s)", unit)
}

type Ssm2ParameterAddress struct {
	Address string `xml:",chardata"`
	Length  int    `xml:"length,attr"` // Not sure what this is used for?
	Bit     int    `xml:"bit,attr"`    // This looks like it's useful when we get to switches, but I'm not there yet
}

func (a Ssm2ParameterAddress) GetAddressBytes() ([]byte, error) {
	if len(a.Address) > 2 {
		return hex.DecodeString(a.Address[2:])
	}
	return []byte{}, fmt.Errorf("Parameter Address malformed %s", a.Address)
}

type Ssm2ParameterConversion struct {
	Units     string  `xml:"units,attr"`
	Expr      string  `xml:"expr,attr"`
	Format    string  `xml:"format,attr"`
	GaugeMin  float64 `xml:"gauge_min,attr"`
	GaugeMax  float64 `xml:"gauge_max,attr"`
	GaugeStep float64 `xml:"gauge_step,attr"`
}
