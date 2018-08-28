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
	"io/ioutil"

	. "github.com/rgeyer/ssm2logger/ssm2lib"
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		non_header_bytes := map[string]int{}
		allbytes, err := ioutil.ReadFile("snooped.bin")
		if err != nil {
			return err
		}

		// for b := range f {
		// 	_, ok := non_header_bytes[fmt.Sprintf("0x%x", b)]
		// 	if ok {
		// 		non_header_bytes[fmt.Sprintf("0x%x", b)] = non_header_bytes[fmt.Sprintf("0x%x", b)] + 1
		// 	} else {
		// 		non_header_bytes[fmt.Sprintf("0x%x", b)] = 1
		// 	}
		// }

		curbyte := 0
		for curbyte < len(allbytes) {
			if allbytes[curbyte] == Ssm2PacketFirstByte {
				data_size := int(allbytes[curbyte+int(Ssm2PacketIndexDataSize)])
				packet_end := curbyte + int(Ssm2PacketMinSize) + data_size
				if len(allbytes) >= packet_end {
					packet_bytes := allbytes[curbyte:packet_end]
					packet := NewPacketFromBytes(packet_bytes)
					curbyte = packet_end - 1
					// packets = append(packets, *packet)
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
				_, ok := non_header_bytes[fmt.Sprintf("0x%.2x", allbytes[curbyte])]
				if ok {
					non_header_bytes[fmt.Sprintf("0x%.2x", allbytes[curbyte])] = non_header_bytes[fmt.Sprintf("0x%x", allbytes[curbyte])] + 1
				} else {
					non_header_bytes[fmt.Sprintf("0x%.2x", allbytes[curbyte])] = 1
				}
				curbyte += 1
			}
		}

		js, err := json.MarshalIndent(non_header_bytes, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(js))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// analyzeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// analyzeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
