package net

import (
	"fmt"
	"log"
	"net"
	"strconv"
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

//zhaojianyun 接口最大网速与网路接口状态
//type SMIORunStatus struct {
//	RunStatus uint32
//	Speed     uint64
//}

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

	//网络接口map
	interfacesByName := map[string]net.Interface{}
	for _, iface := range interfaces {
		interfacesByName[iface.Name] = iface
	}

	//获取网关信息
	gateways := ReadGateways()

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

		//解析mask地址
		ip, mask, _ := ParseIPMask(tiface)

		//接口配置状态
		//var adminStatus uint32
		//flags := strings.Split(tiface.Flags.String(), "|")
		//if len(flags) > 0 {
		//	if strings.ReplaceAll(flags[0], " ", "") == "up" {
		//		adminStatus = 1
		//	}
		//}
		//接口运行状态和网速
		//gateway, ok := gateways[io.Name]
		//if !ok {
		//	gateway = "---"
		//}
		//instates := ReadRunStatus(io.Name)

		fields := map[string]interface{}{
			//"sdd" : iface.
			"index": tiface.Index,
			"name":  tiface.Name,
			"mtu":   tiface.MTU,
			//"speed":        instates.Speed,
			"ip":       ip,
			"net_mask": mask,
			//"gateway":      gateway,
			"mac": tiface.HardwareAddr.String(),
			//"admin_status": adminStatus,
			//"run_status":   instates.RunStatus,
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
 * 函数名： ParseIPMask(iface net.Interface)
 * 作用：根据IP解析Mask地址
 * 返回值：IP地址，MASK地址
 */
func ParseIPMask(iface net.Interface) (string, string, error) {
	ipv4 := "--"
	mask := "--"
	adds, err := iface.Addrs()
	if err != nil {
		log.Fatal("get network addr failed: ", err)
		return ipv4, mask, nil
	}
	for _, ip := range adds {
		if strings.Contains(ip.String(), ".") {
			_, ipNet, err := net.ParseCIDR(ip.String())
			if err != nil {
				return ipv4, mask, nil
			}
			val := make([]byte, len(ipNet.Mask))
			copy(val, ipNet.Mask)
			var s []string
			for _, i := range val[:] {
				s = append(s, strconv.Itoa(int(i)))
			}
			ipv4 = ip.String()[:strings.Index(ip.String(), "/")]
			mask = strings.Join(s, ".")
			break
		}
	}
	return ipv4, mask, nil
}
