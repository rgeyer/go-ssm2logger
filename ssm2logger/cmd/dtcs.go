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
	"encoding/xml"
	"io/ioutil"
	"os"

	. "github.com/rgeyer/ssm2logger/ssm2lib"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// dtcsCmd represents the dtcs command
var dtcsCmd = &cobra.Command{
	Use:   "dtcs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		xmlfile, err := os.Open("logger_STD_EN_v336.xml")
		if err != nil {
			return err
		}
		defer xmlfile.Close()

		xmlbytes, err := ioutil.ReadAll(xmlfile)
		if err != nil {
			return err
		}

		logDefs := &Ssm2Logger{}
		err = xml.Unmarshal(xmlbytes, &logDefs)
		if err != nil {
			return err
		}

		ssm2_conn := &Ssm2Connection{}
		ssm2_conn.SetLogger(logger)

		ssm2_conn.Open(port)
		defer ssm2_conn.Close()

		initResponse, err := ssm2_conn.InitEngine()
		if err != nil {
			return err
		}

		capBytes := initResponse.GetCapabilityBytes()
		var supportedParams []Ssm2Parameter

		for _, proto := range logDefs.Protocols {
			if proto.Id == "SSM" {
				for _, param := range proto.Parameters {
					if param.EcuByteIndex < uint(len(capBytes)) {
						if (capBytes[param.EcuByteIndex] & (1 << param.EcuBit)) > 0 {
							supportedParams = append(supportedParams, param)
						}
					}
				}

				var dtcChunks [][]Ssm2Dtc
				chunkSize := 20

				for i := 0; i < len(proto.Dtcs); i += chunkSize {
					end := i + chunkSize

					if end > len(proto.Dtcs) {
						end = len(proto.Dtcs)
					}

					dtcChunks = append(dtcChunks, proto.Dtcs[i:end])
				}

				logger.WithFields(log.Fields{"dtc_total": len(proto.Dtcs), "chunks": len(dtcChunks)}).Info("Split all possible DTCs into a few read address requests")

				dtcCount := 0
				for _, chunk := range dtcChunks {
					var addresses [][]byte
					for _, dtc := range chunk {
						tmpAddr, err := dtc.GetTmpAddressBytes()
						if err != nil {
							logger.WithFields(log.Fields{"error": err, "dtc": dtc.Name}).Error("Unable to get temporary address location for DTC")
							continue
						}
						memAddr, err := dtc.GetMemAddressBytes()
						if err != nil {
							logger.WithFields(log.Fields{"error": err, "dtc": dtc.Name}).Error("Unable to get stored address location for DTC")
							continue
						}
						addresses = append(addresses, tmpAddr)
						addresses = append(addresses, memAddr)
					}
					response, err := ssm2_conn.ReadAddresses(addresses)
					if err != nil {
						logger.WithFields(log.Fields{"error": err}).Error("Unable to query ECM for DTCs")
						continue
					}
					for idx, dtc := range chunk {
						responseBytes := response.GetData()
						dtc.Set = responseBytes[idx*2]&1<<dtc.Bit > 0
						dtc.Stored = responseBytes[idx*2+1]&1<<dtc.Bit > 0
						if dtc.Set || dtc.Stored {
							dtcCount += 1
							logger.WithFields(log.Fields{"set": dtc.Set, "stored": dtc.Stored}).Info(dtc.Name)
						}
					}
				}
				logger.WithFields(log.Fields{"count": dtcCount}).Info("DTCs found")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dtcsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dtcsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dtcsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
