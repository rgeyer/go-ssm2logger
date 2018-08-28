// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacobsa/go-serial/serial"
	. "github.com/rgeyer/ssm2logger/ssm2lib"
	"github.com/spf13/cobra"
)

// sniffCmd represents the sniff command
var sniffCmd = &cobra.Command{
	Use:   "sniff",
	Short: "Passively listens to traffic on the OBDII bus and logs it to a binary file.",
	Long: `Passively listens to traffic on the OBDII bus and logs it to a binary file.

	It will sniff in a loop, as quickly as possible, until a SIGINT or SIGTERM is encountered.

	Useful for reverse engineering other scantools connected to the ECU or TCU`,
	RunE: func(cmd *cobra.Command, args []string) error {
		options := serial.OpenOptions{
			PortName:              port,
			BaudRate:              4800,
			DataBits:              8,
			StopBits:              1,
			InterCharacterTimeout: 0,
			MinimumReadSize:       1,
			ParityMode:            serial.PARITY_NONE,
		}

		f, err := serial.Open(options)
		if err != nil {
			return fmt.Errorf("Error opening serial port: %s", err)
		} else {
			defer f.Close()
		}

		sigs := make(chan os.Signal, 1)
		loop := true
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for sig := range sigs {
				fmt.Println(sig)
				if sig == syscall.SIGINT || sig == syscall.SIGTERM {
					loop = false
					// TODO: Nasty dirty hack because I need to understand channels and
					// timeouts better, or use better options for opening my serial port.
					// I throw a byte on the wire for my loop to read, allowing it to
					// stop blocking and properly exit.
					// f.Write([]byte{0x00})
				}
			}
		}()

		allbytes := []byte{}
		for loop {
			resp, err := ioutil.ReadAll(f)
			if err != nil {
				if err != io.EOF {
					return fmt.Errorf("Error reading from serial port: %s", err)
				}
			} else {
				allbytes = append(allbytes, resp...)
			}
		}

		fmt.Println("Finished snooping")

		ioutil.WriteFile("snooped.bin", allbytes, 0644)

		packets := []Ssm2Packet{}

		curbyte := 0
		for curbyte < len(allbytes) {
			if allbytes[curbyte] == Ssm2PacketFirstByte {
				data_size := int(allbytes[curbyte+int(Ssm2PacketIndexDataSize)])
				packet_end := curbyte + int(Ssm2PacketMinSize) + data_size
				if len(allbytes) >= packet_end {
					packet_bytes := allbytes[curbyte:packet_end]
					packet := NewPacketFromBytes(packet_bytes)
					curbyte = curbyte + len(packet_bytes)
					packets = append(packets, *packet)
					js, err := json.Marshal(packet)
					if err != nil {
						fmt.Println("Couldn't marshal the packet", err)
					} else {
						fmt.Println(string(js))
						fmt.Println(packet.Bytes())
					}
				} else {
					fmt.Println("Stream ended before remainder of packet arrived")
					curbyte += 1
				}
			} else {
				fmt.Println(fmt.Sprintf("0x%.2x was not a header byte", allbytes[curbyte]))
				curbyte += 1
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sniffCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sniffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sniffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
