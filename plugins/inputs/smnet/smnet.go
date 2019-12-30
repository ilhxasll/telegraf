package net

import (
	"fmt"
	"net"
	"strings"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/filter"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/inputs/system"
)

//zhaojyun smnet
type SMNetIOStats struct {
	filter filter.Filter
	ps     system.PS

	skipChecks          bool
	IgnoreProtocolStats bool
	Interfaces          []string
}

func (_ *SMNetIOStats) Description() string {
	return "Read metrics about network interface usage"
}

var smNetSampleConfig = `
  ## By default, telegraf gathers stats from any up interface (excluding loopback)
  ## Setting interfaces will tell it to gather these explicit interfaces,
  ## regardless of status.
  ##
  # interfaces = ["eth0"]
  ##
  ## On linux systems telegraf also collects protocol stats.
  ## Setting ignore_protocol_stats to true will skip reporting of protocol metrics.
  ##
  # ignore_protocol_stats = false
  ##
`

func (_ *SMNetIOStats) SampleConfig() string {
	return smNetSampleConfig
}

func (s *SMNetIOStats) Gather(acc telegraf.Accumulator) error {
	netio, err := s.ps.NetIO()
	if err != nil {
		return fmt.Errorf("error getting net io info: %s", err)
	}

	if s.filter == nil {
		if s.filter, err = filter.Compile(s.Interfaces); err != nil {
			return fmt.Errorf("error compiling filter: %s", err)
		}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("error getting list of interfaces: %s", err)
	}
	interfacesByName := map[string]net.Interface{}
	for _, iface := range interfaces {
		interfacesByName[iface.Name] = iface
	}

	for _, io := range netio {
		if len(s.Interfaces) != 0 {
			var found bool

			if s.filter.Match(io.Name) {
				found = true
			}

			if !found {
				continue
			}
		} else if !s.skipChecks {
			iface, ok := interfacesByName[io.Name]
			if !ok {
				continue
			}

			if iface.Flags&net.FlagLoopback == net.FlagLoopback {
				continue
			}

			if iface.Flags&net.FlagUp == 0 {
				continue
			}
		}

		tags := map[string]string{
			"interface": io.Name,
		}

		tiface, _ := interfacesByName[io.Name]
		fields := map[string]interface{}{
			//"sdd" : iface.
			"index":        tiface.Index,
			"name":         tiface.Name,
			"mtu":          tiface.MTU,
			"speed":        0,
			"ip":           0,
			"net_mask":     0,
			"gateway":      0,
			"mac":          tiface.HardwareAddr.String(),
			"admin_Status": 0,
			"run_state":    0,
			"bytes_sent":   io.BytesSent,
			"bytes_recv":   io.BytesRecv,
			"packets_sent": io.PacketsSent,
			"packets_recv": io.PacketsRecv,
			"err_in":       io.Errin,
			"err_out":      io.Errout,
			"drop_in":      io.Dropin,
			"drop_out":     io.Dropout,
		}
		acc.AddCounter("smnet", fields, tags)
	}

	// Get system wide stats for different network protocols
	// (ignore these stats if the call fails)
	if !s.IgnoreProtocolStats {
		netprotos, _ := s.ps.NetProto()
		fields := make(map[string]interface{})
		for _, proto := range netprotos {
			for stat, value := range proto.Stats {
				name := fmt.Sprintf("%s_%s", strings.ToLower(proto.Protocol),
					strings.ToLower(stat))
				fields[name] = value
			}
		}
		tags := map[string]string{
			"interface": "all",
		}
		acc.AddFields("smnet", fields, tags)
	}

	return nil
}

func init() {
	inputs.Add("smnet", func() telegraf.Input {
		return &SMNetIOStats{ps: system.NewSystemPS()}
	})
}
