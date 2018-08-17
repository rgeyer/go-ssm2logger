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
	"time"

	"github.com/jacobsa/go-serial/serial"
	. "github.com/rgeyer/ssm2logger/ssm2lib"
	"github.com/spf13/cobra"
)

// sniffCmd represents the sniff command
var sniffCmd = &cobra.Command{
	Use:   "sniff",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		options := serial.OpenOptions{
			PortName:              "/dev/tty.usbserial-1420",
			BaudRate:              4800,
			DataBits:              8,
			StopBits:              1,
			InterCharacterTimeout: 100,
			MinimumReadSize:       0,
			ParityMode:            serial.PARITY_NONE,
		}

		f, err := serial.Open(options)

		if err != nil {
			return fmt.Errorf("Error opening serial port: %s", err)
		} else {
			defer f.Close()
		}

		allbytes := []byte{}
		start := time.Now().Second()
		end := start + 10
		duration := start
		for duration < end {
			resp, err := ioutil.ReadAll(f)
			if err != nil {
				if err != io.EOF {
					return fmt.Errorf("Error reading from serial port: %s", err)
				}
			} else {
				fmt.Println("Fetched bytes ", len(resp))
				allbytes = append(allbytes, resp...)
			}
			duration = time.Now().Second()
		}

		fmt.Println("Finished snooping")

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
						fmt.Println(packet.GetBytes())
					}
				} else {
					fmt.Println("Stream ended before remainder of packet arrived")
				}
			} else {
				fmt.Println(fmt.Sprintf("0x%x was not a header byte", allbytes[curbyte]))
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
