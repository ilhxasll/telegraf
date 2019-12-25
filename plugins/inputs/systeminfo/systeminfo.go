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

	f, err := os.Open("/etc/.systeminfo")
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
		if strings.Contains(line, "=") {
			a := strings.Split(line, "=")
			if len(a) > 1 {
				b := a[0]
				c := strings.Trim(a[1], "\n")
				if b == "产品名称" {
					fields["pro_name"] = c
				}
				if b == "ProductName" {
					fields["pro_name"] = c
				}
				if b == "产品型号" {
					fields["pro_number"] = c
				}
				if b == "ProductModel" {
					fields["pro_number"] = c
				}
				if b == "标识码（产品唯一标识）" {
					fields["pro_code"] = c
				}
				if b == "ID" { //安全卡标识码（产品唯一标识）
					fields["pro_code"] = c
				}
				if b == "电磁泄漏发射防护类型" {
					fields["launch_type"] = c
				}
				if b == "ShelterModel" {
					fields["launch_type"] = c
				}
				if b == "生产者（制造商）" {
					fields["manufacturer"] = c
				}
				if b == "Producter" {
					fields["manufacturer"] = c
				}
				if b == "操作系统名称" {
					fields["os_name"] = c
				}
				if b == "Name" {
					fields["os_name"] = c
				}
				if b == "系统版本" {
					fields["sys_version"] = c
				}
				if b == "Release" {
					fields["sys_version"] = c
				}
				if b == "内核版本" {
					fields["kernel"] = c
				}
				if b == "Kernel" {
					fields["kernel"] = c
				}
				if b == "系统位数" {
					fields["sys_number"] = c
				}
				if b == "Bit" {
					fields["sys_number"] = c
				}
				if b == "I/O保密管理模块" {
					fields["io_sec_model"] = c
				}
				if b == "安全卡版本" {
					fields["safe_number"] = c
				}
				if b == "Version" {
					fields["safe_number"] = c
				}
				if b == "固件版本（BIOS）" {
					fields["bios"] = c
				}
				if b == "BiosVersion" {
					fields["bios"] = c
				}
				if b == "处理器信息" {
					fields["cpu_info"] = c
				}
				if b == "CPU" {
					fields["cpu_info"] = c
				}
				if b == "内存" {
					fields["memory"] = c
				}
				if b == "Memory" {
					fields["memory"] = c
				}
				if b == "硬盘序列号" {
					fields["disk_number"] = c
				}
				if b == "HDSerial" {
					fields["disk_number"] = c
				}
				if b == "硬盘容量" {
					fields["disk_capacity"] = c
				}
				if b == "HDCapacity" {
					fields["disk_capacity"] = c
				}
				if b == "主板版本号" {
					fields["mainboard_version"] = c
				}
				if b == "系统安装时间" {
					fields["sys_begin_time"] = c
				}
				if b == "系统更新时间" {
					fields["sys_update_time"] = c
				}
				if b == "UpdateTime" {
					fields["sys_update_time"] = c
				}
				if b == "KernelVersion" { //三合一内核版本
					fields["three_kernel"] = c
				}
				if b == "SoftWareVersion" { //三合一软件版本
					fields["three_version"] = c
				}
				if b == "Product" { //操作系统名称
					fields["sys_product"] = c
				}
				if b == "HDSerial_1" { //HOME 硬盘序列号
					fields["home_disk_number"] = c
				}
				if b == "HDCapacity_1" { //HOME 硬盘容量
					fields["home_disk_capacity"] = c
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
		} else {
			a := strings.Split(line, "：")
			if len(a) > 1 {
				b := a[0]
				c := strings.Trim(a[1], "\n")
				if b == "产品名称" {
					fields["pro_name"] = c
				}
				if b == "产品型号" {
					fields["pro_number"] = c
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
				if b == "操作系统名称" {
					fields["sys_product"] = c
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
				if b == "安全卡版本" {
					fields["safe_number"] = c
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
