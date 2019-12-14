package systeminfo

import (
	"bufio"
	"fmt"
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

	f, err := os.Open("/etc/.systeminfo")
	if err != nil {
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	l := 0
	var fields map[string]string
	fields = make(map[string]string)
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
				fields["pro_name"] = c
			}
			if b == "标识码（产品唯一标识）" {
				fields["pro_code"] = c
			}
			if b == "电磁泄露发射防护类型" {
				fields["launch_type"] = c
			}
			if b == "生产者（制造商）" {
				fields["manufacturer"] = c
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
					return nil
				} else {
					l = l + 1
				}
				return err
			}
		}
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	fields2 := map[string]interface{}{
		"pro_name":          fields["pro_name"],
		"pro_code":          fields["pro_code"],
		"launch_type":       fields["launch_type"],
		"manufacturer":      fields["manufacturer"],
		"sys_version":       fields["sys_version"],
		"kernel":            fields["kernel"],
		"sys_number":        fields["sys_number"],
		"three_kernel":      fields["three_kernel"],
		"three_version":     fields["three_version"],
		"bios":              fields["bios"],
		"cpu_info":          fields["cpu_info"],
		"memory":            fields["memory"],
		"disk_number":       fields["disk_number"],
		"disk_capacity":     fields["disk_capacity"],
		"mainboard_version": fields["mainboard_version"],
		"sys_update_time":   fields["sys_update_time"],
	}

	acc.AddGauge("systeminfo", fields2, nil)
	return nil
}

func init() {
	ps := system.NewSystemPS()
	inputs.Add("systeminfo", func() telegraf.Input {
		return &SysInfoStats{ps: ps}
	})
}
