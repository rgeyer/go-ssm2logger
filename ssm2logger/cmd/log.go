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
	"encoding/csv"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "github.com/rgeyer/ssm2logger/ssm2lib"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logfile_path string

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Logs SSM2 data, along with any plugin data to a csv file",
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
					//fmt.Println(resp_bytes[8] & (1 << 6)) A test of looking for a specific parameter using bitwise operators. Gotta move this elsewhere.
					if param.EcuByteIndex < uint(len(capBytes)) {
						if (capBytes[param.EcuByteIndex] & (1 << param.EcuBit)) > 0 {
							supportedParams = append(supportedParams, param)
						}
					}
				}
			}
		}

		logger.WithFields(log.Fields{
			"SsmId":                  hex.EncodeToString(initResponse.GetSsmId()),
			"RomId":                  hex.EncodeToString(initResponse.GetRomId()),
			"Supported Capabilities": len(supportedParams),
		}).Info("Initialized ECM")

		// Cooldown between writes?!
		time.Sleep(200 * time.Millisecond)

		// readResponse, err := ssm2_conn.ReadAddresses([]byte{0x46, 0x3c, 0x3D, 0x1C, 0x20, 0x22, 0x29, 0x32, 0x3B, 0xA, 0xD, 0x11, 0xE, 0x9, 0x13})
		// if err != nil {
		// 	return err
		// }

		sigs := make(chan os.Signal, 1)
		loop := true
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for sig := range sigs {
				if sig == syscall.SIGINT || sig == syscall.SIGTERM {
					loop = false
				}
			}
		}()

		timestamp := time.Now()
		logfilename := fmt.Sprintf("%s/%s-%d-log.csv", logfile_path, hex.EncodeToString(initResponse.GetRomId()), timestamp.Unix())

		csvfile, err := os.Create(logfilename)
		if err != nil {
			return err
		}

		defer csvfile.Close()

		writer := csv.NewWriter(csvfile)
		defer writer.Flush()

		header := []string{"timestamp"}
		var actual_ecu_params []Ssm2Parameter
		for _, param := range supportedParams {
			logger.Info(param)
			if param.EcuByteIndex > 0 {
				header = append(header, fmt.Sprintf("%s (%s)", param.Name, param.Conversions[0].Units))
				actual_ecu_params = append(actual_ecu_params, param)
			}
		}

		writer.Write(header)

		_, err = ssm2_conn.ReadParameters(actual_ecu_params)
		if err != nil {
			return err
		}

		for loop {
			readPacket, err := ssm2_conn.GetNextPacketInStream()
			if err != nil {
				return err
			}

			readResponseBytes := readPacket.GetData()
			logger.WithFields(log.Fields{"length": len(readResponseBytes), "data": hex.EncodeToString(readResponseBytes)}).Debug("Read Response")

			row := []string{fmt.Sprintf("%d", time.Now().Unix())}
			for idx, param := range actual_ecu_params {
				length := 1
				if param.Address.Length > 1 {
					length = param.Address.Length
				}
				value := readResponseBytes[idx : idx+length]
				convertedValue, err := param.Convert(param.Conversions[0].Units, value)
				if err != nil {
					return err
				}
				row = append(row, fmt.Sprintf("%f", convertedValue))
				logger.WithFields(log.Fields{
					"param":           param.Name,
					"length":          length,
					"byteval":         hex.EncodeToString(value),
					"expr":            param.Conversions[0].Expr,
					"converted_value": convertedValue,
				}).Debug("Parameter Conversion")
			}

			writer.Write(row)
		}

		logger.Info("Received Stop Signal and discontinued logging")

		return nil
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

	logCmd.Flags().StringVar(&logfile_path, "logfile-path", "", "Path where the logfile will be generated. The actual file will be <logfile-path>/<ecu romid>-<timestamp>-log.csv. Default is the current directory.")
	viper.BindPFlag("logfile-path", logCmd.Flags().Lookup("logfile-path"))
	viper.SetDefault("logfile-path", ".")
}
