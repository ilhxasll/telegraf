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
	var tags map[string]string
	tags = make(map[string]string)
	for {
		line, err := bfRd.ReadString('\n')
		a := strings.Split(line, "：")
		if len(a) > 1 {
			b := a[0]
			c := a[1]
			if b == "测试" {
				tags["测试"] = c
			}
			if b == "产品名称" {
				tags["pro_name"] = c
			}
			if b == "标识码（产品唯一标识）" {
				tags["pro_code"] = c
			}
			if b == "电磁泄露发射防护类型" {
				tags["launch_type"] = c
			}
			if b == "生产者（制造商）" {
				tags["manufacturer"] = c
			}
			if b == "系统版本" {
				tags["sys_version"] = c
			}
			if b == "内核版本" {
				tags["kernel"] = c
			}
			if b == "系统位数" {
				tags["sys_number"] = c
			}
			if b == "三合一内核版本" {
				tags["three_kernel"] = c
			}
			if b == "三合一软件版本" {
				tags["three_version"] = c
			}
			if b == "固件版本（BIOS）" {
				tags["bios"] = c
			}
			if b == "处理器信息" {
				tags["cpu_info"] = c
			}
			if b == "内存" {
				tags["memory"] = c
			}
			if b == "硬盘序列号" {
				tags["disk_number"] = c
			}
			if b == "硬盘容量" {
				tags["disk_capacity"] = c
			}
			if b == "主板版本号" {
				tags["mainboard_version"] = c
			}
			if b == "系统更新时间" {
				tags["sys_update_time"] = c
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
	fmt.Println("fields", fields)
	acc.AddGauge("systeminfo", nil, tags)
	return nil
}

func init() {
	ps := system.NewSystemPS()
	inputs.Add("systeminfo", func() telegraf.Input {
		return &SysInfoStats{ps: ps}
	})
}
