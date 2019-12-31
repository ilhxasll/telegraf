package smnet

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"

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

		//接口配置状态
		var adminStatus uint32
		flags := strings.Split(tiface.Flags.String(), "|")
		if len(flags) > 0 {
			if strings.ReplaceAll(flags[0], " ", "") == "up" {
				adminStatus = 1
			}
		}

		//接口运行状态和网速
		gateway, ok := gateways[io.Name]
		if !ok {
			gateway = "000"
		}

		instates := ReadRunStatus(io.Name)

		fields := map[string]interface{}{
			//"sdd" : iface.
			"index":        tiface.Index,
			"name":         tiface.Name,
			"mtu":          tiface.MTU,
			"speed":        instates.Speed,
			"ip":           0,
			"net_mask":     0,
			"gateway":      gateway,
			"mac":          tiface.HardwareAddr.String(),
			"admin_Status": adminStatus,
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
 * 函数名：delete_extra_space(s string) string
 * 功  能:删除字符串中多余的空格(含tab)，有多个空格时，仅保留一个空格，同时将字符串中的tab换为空格
 * 参  数:s string:原始字符串
 * 返回值:string:删除多余空格后的字符串
 */
func deleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)       //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}

/*
 * 函数名：Readgateways() map[string]string
 * 作用：读取网关信息
 * 返回值：map[网络接口名]网关地址
 */
func ReadGateways() map[string]string {
	cmd := exec.Command("route", "-n")
	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
	}
	//执行命令
	if err := cmd.Start(); err != nil {
	}
	//使用带缓冲的读取器
	outputBuf := bufio.NewReader(stdout)
	var gateways map[string]string /*创建集合 */
	gateways = make(map[string]string)
	var i int
	for {
		//一次获取一行,_ 获取当前行是否被读完
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			break
		}
		if i < 2 {
			i++
			continue
		}
		tempgate := strings.Split(deleteExtraSpace(string(output)), " ")

		if len(tempgate) == 8 {

			_, ok := gateways[tempgate[7]]
			if ok {
				continue
			}
			gateways[tempgate[7]] = tempgate[1]
		}
	}

	//wait 方法会一直阻塞到其所属的命令完全运行结束为止
	if err := cmd.Wait(); err != nil {
	}

	return gateways
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
