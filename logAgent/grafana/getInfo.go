package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"
)

// cpu info
func getCpuInfo() {
	var cpuInfo = new(CpuInfo)
	percent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("cpu percent:%v\n", percent)
	cpuInfo.CpuPercent = percent[0]
	writesCpuPoints(cpuInfo)
}

func getMemInfo() {
	var memInfo = new(MemInfo)
	info, _ := mem.VirtualMemory()
	//fmt.Printf("mem info:%v\n", info)
	memInfo.Total = info.Total
	memInfo.Available = info.Available
	memInfo.Used = info.Used
	memInfo.UsedPercent = info.UsedPercent
	memInfo.Buffers = info.Buffers
	memInfo.Cached = info.Cached
	writesMemPoints(memInfo)
}

// disk info
func getDiskInfo() {
	var diskInfo = &DiskInfo{PartitionUsageStat: make(map[string]*disk.UsageStat, 16)}
	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Printf("get Partitions failed, err:%v\n", err)
		return
	}
	for _, part := range parts {
		usageStat, _ := disk.Usage(part.Mountpoint)
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageStat
	}
	disk.IOCounters()
	//ioStat, _ := disk.IOCounters()
	//for _, _ = range ioStat {
	//	fmt.Printf("%v:%v\n", k, v)
	//}
	writesDiskPoints(diskInfo)
}

func getNetInfo() {
	var netInfo = &NetInfo{NetIOCountersStat: make(map[string]*IOStat, 8)}
	currentTimeStamp := time.Now().Unix()
	info, _ := net.IOCounters(true)
	for _, v := range info {
		var ioStat = new(IOStat)
		ioStat.BytesSent = v.BytesSent
		ioStat.BytesRecv = v.BytesRecv
		ioStat.PacketsSent = v.PacketsSent
		ioStat.PacketsRecv = v.PacketsRecv
		//将具体网卡数据的ioStat变量添加到map中
		netInfo.NetIOCountersStat[v.Name] = ioStat
		if lastNetIOSTimeStamp == 0 || lastNetInfo == nil {
			continue
		}
		//时间间隔
		interval := currentTimeStamp - lastNetIOSTimeStamp
		//计算速率
		ioStat.BytesSentRate = (float64(ioStat.BytesSent) - float64(lastNetInfo.NetIOCountersStat[v.Name].BytesSent)) / float64(interval)
		ioStat.BytesRecvRate = (float64(ioStat.BytesRecv) - float64(lastNetInfo.NetIOCountersStat[v.Name].BytesRecv)) / float64(interval)
		ioStat.PacketsSentRate = (float64(ioStat.PacketsSent) - float64(lastNetInfo.NetIOCountersStat[v.Name].PacketsSent)) / float64(interval)
		ioStat.PacketsRecvRate = (float64(ioStat.PacketsRecv) - float64(lastNetInfo.NetIOCountersStat[v.Name].PacketsRecv)) / float64(interval)

	}
	//更新全局记录的上一次采集网卡的时间点和网卡数据
	lastNetIOSTimeStamp = currentTimeStamp
	lastNetInfo = netInfo
	//发送influxDB
	writesNetPoints(netInfo)
}
