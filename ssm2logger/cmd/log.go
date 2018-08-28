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
	"encoding/hex"
	"fmt"
	"time"

	. "github.com/rgeyer/ssm2logger/ssm2lib"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		ssm2_conn := &Ssm2Connection{}
		ssm2_conn.SetLogger(logger)

		ssm2_conn.Open(port)

		initResponse, err := ssm2_conn.InitEngine()
		if err != nil {
			return err
		}

		fmt.Println(initResponse.ToJson())
		fmt.Println(hex.EncodeToString(initResponse.Bytes()))

		// Cooldown between writes?!
		time.Sleep(200 * time.Millisecond)

		readResponse, err := ssm2_conn.ReadAddresses([]byte{0x46})
		if err != nil {
			return err
		}

		fmt.Println(readResponse.ToJson())

		ssm2_conn.Close()

		return nil

		/*

				// writeFastModePacket := NewWriteAddressPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10, []byte{0x00, 0x01, 0x98}, []byte{0x5a})
				writeFastModePacket := NewReadAddressPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10, []byte{0x00, 0x01, 0x98})

				fmt.Println(fmt.Sprintf("Sending %s as bytes %s", hex.EncodeToString(writeFastModePacket.GetBytes())))

				count, err = f.Write(writeFastModePacket.GetBytes())
				if err != nil {
					return fmt.Errorf("Failed to send serial: %s", err)
				}
				fmt.Println("Wrote bytes", count)

				time.Sleep(500 * time.Millisecond)

				resp, err = ioutil.ReadAll(f)
				if err != nil {
					if err != io.EOF {
						fmt.Println("Error reading from serial port: ", err)
					}
				} else {
					fmt.Println(hex.EncodeToString(resp))
					fmt.Println("Tx: ", hex.EncodeToString(resp[:count]))
					if len(resp) > count {
						fmt.Println("Rx: ", hex.EncodeToString(resp[count:]))
					} else {
						return fmt.Errorf("Unable to read fast mode address")
					}

					resp_bytes := resp[count:]
					response_packet := NewPacketFromBytes(resp_bytes)
					fmt.Println(hex.EncodeToString(response_packet.GetData()))
				}

				f.Close()

				options.BaudRate = 10400

				f, err = serial.Open(options)

				if err != nil {
					return fmt.Errorf("Error (re)opening serial port in fast mode: %s", err)
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
					}
				}
			}()

			csvfile, err := os.Create("log.csv")
			if err != nil {
				return err
			}

			defer csvfile.Close()

			writer := csv.NewWriter(csvfile)
			defer writer.Flush()

			timestamp := time.Now()
			readPid := NewReadAddressContinuousPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10, []byte{0x46, 0x3c, 0x3D, 0x1C, 0x20, 0x22, 0x29, 0x32, 0x3B, 0xA, 0xD, 0x11, 0xE, 0x9, 0x13})
			// readPid := NewReadAddressPacket(Ssm2DeviceDiagnosticToolF0, Ssm2DeviceEngine10, []byte{0x46, 0x3c, 0x3D, 0x1C, 0x20, 0x22, 0x29, 0x32, 0x3B, 0xA, 0xD, 0x11, 0xE, 0x9, 0x13})
			count, err = f.Write(readPid.GetBytes())
			if err != nil {
				return fmt.Errorf("Failed to send serial: %s", err)
			}
			fmt.Println("Wrote bytes", count)

			time.Sleep(50 * time.Millisecond)

			var readPacket *Ssm2Packet
			resp, err = ioutil.ReadAll(f)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading from serial port: ", err)
				}
			} else {
				fmt.Println("Read bytes", len(resp))
				if len(resp) > count {
					readPacket = NewPacketFromBytes(resp[count:])
				}
				// js, err := json.Marshal(readPacket)
				// if err != nil {
				// 	fmt.Println("Couldn't marshal the readPacket", err)
				// } else {
				// 	fmt.Println(string(js))
				// }
				// fmt.Println("Tx: ", hex.EncodeToString(resp[:count]))
				// fmt.Println("Rx: ", hex.EncodeToString(resp[count:]))

				// for _, val := range resp {
				// 	fmt.Println(val)
				// }

				// writer.Write([]string{timestamp.String(), string(readPacket.GetData())})
			}

			for loop {
				time.Sleep(50 * time.Millisecond)
				timestamp = time.Now()
				resp, err = ioutil.ReadAll(f)
				if err != nil {
					if err != io.EOF {
						fmt.Println("Error reading from serial port: ", err)
					}
				} else {
					fmt.Println("Read bytes", len(resp))
					if len(resp) > 0 {
						readPacket = NewPacketFromBytes(resp)

						writer.Write([]string{timestamp.String(), string(readPacket.GetData())})
					}
				}
			}
			return nil*/
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
