package smnet

import (
	"fmt"
	"net"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/filter"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/inputs/system"
	"github.com/safchain/ethtool"
)

//zhaojyun smnet
type SMNetIOStats struct {
	filter filter.Filter
	ps     system.PS

	skipChecks          bool
	IgnoreProtocolStats bool
	Interfaces          []string
}

//zhaojianyun 接口最大网速与网路接口状态
type SMIORunStatus struct {
	RunStatus uint32
	Speed     uint64
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

		instates := ReadRunStatus(io.Name)
		tiface, _ := interfacesByName[io.Name]
		fields := map[string]interface{}{
			//"sdd" : iface.
			"index":        tiface.Index,
			"name":         tiface.Name,
			"mtu":          tiface.MTU,
			"speed":        instates.Speed,
			"ip":           0,
			"net_mask":     0,
			"gateway":      0,
			"mac":          tiface.HardwareAddr.String(),
			"admin_Status": 0,
			"run_state":    instates.RunStatus,
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

	return nil
}

func init() {
	inputs.Add("smnet", func() telegraf.Input {
		return &SMNetIOStats{ps: system.NewSystemPS()}
	})
}

/*
 * 函数名： ReadRunStatus(ifacename string)
 * 作用：获取接口运行状态和网速
 * 返回值：SMIORunStatus
 */
func ReadRunStatus(ifacename string) SMIORunStatus {

	//获取ethtool命令句柄
	ethHandle, err := ethtool.NewEthtool()

	defer ethHandle.Close()

	var instates SMIORunStatus

	//获取接口运行状态
	stats, err := ethHandle.LinkState(ifacename)
	if err == nil {
		instates.RunStatus = stats
	}

	//获取网速
	result, err := ethHandle.CmdGetMapped(ifacename)
	if err == nil {
		speed, ok := result["speed"]
		if ok && speed != 4294967295 {
			instates.Speed = speed
		}
	}
	return instates
}
