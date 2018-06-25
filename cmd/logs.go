// Copyright © 2018 Scott Ware <scottdware@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/scottdware/go-easycsv"

	"github.com/scottdware/go-panos"

	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve logs from the device and export them to a CSV file",
	Long: `You can query the device logs via this tool the same way you would on the GUI.
The different log types you can retrieve are:

config, system, traffic, threat, wildfire, url, data

When using the --query flag, be sure to enclose your criteria in quotes "" like so:

--query "(addr.src in 10.0.0.0/8)"

The default search type is traffic. Based on your query, and the device, log retrieval
and export could take a while.`,
	Run: func(cmd *cobra.Command, args []string) {
		pass := passwd()
		creds := &panos.AuthMethod{
			Credentials: []string{user, pass},
		}

		pan, err := panos.NewSession(device, creds)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		params := &panos.LogParameters{
			Query: query,
			NLogs: nlogs,
		}

		jobID, err := pan.QueryLogs(ltype, params)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		time.Sleep(time.Duration(lwait) * time.Second)

		logs, err := pan.RetrieveLogs(jobID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !strings.Contains(fh, ".csv") {
			fh = fmt.Sprintf("%s.csv", fh)
		}

		csv, err := easycsv.NewCSV(fh)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		csv.Write("DeviceName,Serial,Rule,TimeGenerated,TimeReceived,Type,Subtype,From,To,Source,SourceUser")
		csv.Write(",SourcePort,SourceCountry,Destination,DestinationPort,DestinationCountry,Application")
		csv.Write(",Action,NAT,NATSourceIP,NATSourcePort,NATDestinationIP,NATDestinationPort,Packets")
		csv.Write(",PacketsSent,PacketsReceived,Bytes,BytesSent,BytesReceived,Protocol,SessionID")
		csv.Write(",ParentSessionID,SessionEndReason,RepeatCount,Start,Elapsed,Category,ThreatCategory")
		csv.Write(",ThreatName,ThreatID,Misc,Severity,Direction,InboundInterface,OutboundInterface,ID")
		csv.Write(",Domain,ReceiveTime,SequenceNumber,ActionFlags,ConfigVersion,Vsys,Logset,Flags")
		csv.Write(",Pcap,PcapID,Flagged,Proxy,URLDenied,CaptivePortal,NonStandardDestinationPort,Transaction")
		csv.Write(",PBFClient2Server,PBFServer2Client,TemporaryMatch,SymmetricReturn,SSLDecryptMirror")
		csv.Write(",CredentialDetected,MPTCP,TunnelInspected,ReconExcluded,TunnelType,TPadding,CPadding")
		csv.Write(",TunnelIMSI,VsysID,ReportID,URLIndex,HTTPMethod,XForwardedFor,Referer,UserAgent")
		csv.Write(",SignatureFlags,ContentVersion,FileDigest,Filetype,Sender,Recipient,Subject,Cloud,Padding")
		csv.Write(",ActionSource,TunnelID,IMSI,MonitorTag,IMEI,DeviceGroupHierarchy1,DeviceGroupHierarchy2")
		csv.Write(",DeviceGroupHierarchy3,DeviceGroupHierarchy4,Host,Command,Admin,Client,Result,Path")
		csv.Write(",BeforeChangePreview,AfterChangePreview,FullPath,EventID,Module,Description\n")

		time.Sleep(50 * time.Millisecond)

		for _, l := range logs.Logs {
			csv.Write(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
				l.DeviceName, l.Serial, l.Rule, l.TimeGenerated, l.TimeReceived, l.Type,
				l.Subtype, l.From, l.To, l.Source, l.SourceUser))
			csv.Write(fmt.Sprintf(",%d,%s,%s,%d,%s,%s,%s,%s,%s,%d,%s",
				l.SourcePort, l.SourceCountry, l.Destination, l.DestinationPort, l.DestinationCountry,
				l.Application, l.Action, l.NAT, l.NATSourceIP, l.NATSourcePort, l.NATDestinationIP))
			csv.Write(fmt.Sprintf(",%d,%d,%d,%d,%d,%d,%d,%s,%d,%d,%s,%d",
				l.NATDestinationPort, l.Packets, l.PacketsSent, l.PacketsReceived, l.Bytes,
				l.BytesSent, l.BytesReceived, l.Protocol, l.SessionID, l.ParentSessionID, l.SessionEndReason, l.RepeatCount))
			csv.Write(fmt.Sprintf(",%s,%s,%s,%s,\"%s\",%d,\"%s\",%s,%s,%s,%s",
				l.Start, l.Elapsed, l.Category, l.ThreatCategory, l.ThreatName, l.ThreatID, l.Misc,
				l.Severity, l.Direction, l.InboundInterface, l.OutboundInterface))
			csv.Write(fmt.Sprintf(",%d,%d,%s,%s,%s,%d,%s,%s,%s,%s",
				l.ID, l.Domain, l.ReceiveTime, l.SequenceNumber, l.ActionFlags, l.ConfigVersion, l.Vsys,
				l.Logset, l.Flags, l.Pcap))
			csv.Write(fmt.Sprintf(",%d,%s,%s,%s,%s,%s,%s,%s,%s,%s",
				l.PcapID, l.Flagged, l.Proxy, l.URLDenied, l.CaptivePortal, l.NonStandardDestinationPort,
				l.Transaction, l.PBFClient2Server, l.PBFServer2Client, l.TemporaryMatch))
			csv.Write(fmt.Sprintf(",%s,%s,%s,%s,%s,%s,%s,%d,%d,%d",
				l.SymmetricReturn, l.SSLDecryptMirror, l.CredentialDetected, l.MPTCP, l.TunnelInspected,
				l.ReconExcluded, l.TunnelType, l.TPadding, l.CPadding, l.TunnelIMSI))
			csv.Write(fmt.Sprintf(",%d,%d,%d,%s,%s,%s,\"%s\",%s,%s,%s",
				l.VsysID, l.ReportID, l.URLIndex, l.HTTPMethod, l.XForwardedFor,
				l.Referer, l.UserAgent, l.SignatureFlags, l.ContentVersion, l.FileDigest))
			csv.Write(fmt.Sprintf(",%s,%s,%s,%s,%s,%d,%s,%d,%s,%s,%s,%d,%d,%d",
				l.Filetype, l.Sender, l.Recipient, l.Subject, l.Cloud, l.Padding, l.ActionSource, l.TunnelID,
				l.IMSI, l.MonitorTag, l.IMEI, l.DeviceGroupHierarchy1, l.DeviceGroupHierarchy2, l.DeviceGroupHierarchy3))
			csv.Write(fmt.Sprintf(",%d,%s,%s,%s,%s,%s,\"%s\",\"%s\",\"%s\",\"%s\",%s,%s,\"%s\"\n",
				l.DeviceGroupHierarchy4, l.Host, l.Command, l.Admin, l.Client, l.Result, l.Path,
				l.BeforeChangePreview, l.AfterChangePreview, l.FullPath, l.EventID, l.Module, l.Description))

			time.Sleep(10 * time.Millisecond)
		}

		csv.End()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().StringVarP(&query, "query", "q", "", "Critera to search the logs on")
	logsCmd.Flags().StringVarP(&ltype, "type", "t", "traffic", "Log type to search under")
	logsCmd.Flags().IntVarP(&nlogs, "nlogs", "n", 20, "Number of logs to retrieve")
	logsCmd.Flags().IntVarP(&lwait, "wait", "w", 5, "Wait time in seconds to delay retrieving logs - helpful for large queries")
	logsCmd.Flags().StringVarP(&user, "user", "u", "", "User to connect to the device as")
	logsCmd.Flags().StringVarP(&device, "device", "d", "", "Firewall or Panorama device to connect to")
	logsCmd.Flags().StringVarP(&fh, "file", "o", "", "Name of the CSV file to export the logs to")
	logsCmd.MarkFlagRequired("user")
	logsCmd.MarkFlagRequired("device")
	logsCmd.MarkFlagRequired("export")
}
