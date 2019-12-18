package systeminfo

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/inputs/system"
)

type SysInfoStats struct {
	ps system.PS
}

func (_ *SysInfoStats) Description() string {
	return "Read metrics about /etc/.systeminfo"
}

func (_ *SysInfoStats) SampleConfig() string { return "" }

func (_ *SysInfoStats) Gather(acc telegraf.Accumulator) error {

	f, err := os.Open("/etc/.systeminfo1")
	if err != nil {
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	l := 0
	var fields map[string]interface{}
	fields = make(map[string]interface{})
	for {
		line, err := bfRd.ReadString('\n')
		a := strings.Split(line, "：")
		if len(a) > 1 {
			b := a[0]
			c := a[1]
			if b == "测试" {
				fields["测试"] = c
			}
			if b == "产品名称" {
				fields["pro_name"] = "测试产品"
			}
			if b == "标识码（产品唯一标识）" {
				fields["pro_code"] = 1
			}
			if b == "电磁泄露发射防护类型" {
				fields["launch_type"] = 2
			}
			if b == "生产者（制造商）" {
				fields["manufacturer"] = "测试产品"
			}
			if b == "系统版本" {
				fields["sys_version"] = c
			}
			if b == "内核版本" {
				fields["kernel"] = c
			}
			if b == "系统位数" {
				fields["sys_number"] = c
			}
			if b == "三合一内核版本" {
				fields["three_kernel"] = c
			}
			if b == "三合一软件版本" {
				fields["three_version"] = c
			}
			if b == "固件版本（BIOS）" {
				fields["bios"] = c
			}
			if b == "处理器信息" {
				fields["cpu_info"] = c
			}
			if b == "内存" {
				fields["memory"] = c
			}
			if b == "硬盘序列号" {
				fields["disk_number"] = c
			}
			if b == "硬盘容量" {
				fields["disk_capacity"] = c
			}
			if b == "主板版本号" {
				fields["mainboard_version"] = c
			}
			if b == "系统更新时间" {
				fields["sys_update_time"] = c
			}
		} else {
			if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
				if err == io.EOF {
					break
				} else {
					l = l + 1
				}
				break
			}
		}
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				break
			}
			return err
		}
	}
	acc.AddGauge("systeminfo", fields, nil)
	return nil
}

func init() {
	ps := system.NewSystemPS()
	inputs.Add("systeminfo", func() telegraf.Input {
		return &SysInfoStats{ps: ps}
	})
}
