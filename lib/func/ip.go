package lib

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	confLib "go-gateway/lib/conf"
	"math/rand"
	"net"
	"os"
	"time"
)

func GetLocalIPs() (ips []net.IP) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			}
		}
	}
	return ips
}

func NewSpanId() string {
	// 将其转换为 uint32 类型，确保时间戳的长度不超过 32 位
	timestamp := uint32(time.Now().Unix())
	// 获取本地 IP 地址（confLib.LocalIP）的 IPv4 地址，并将其转换为一个 4 字节的大端序整数。
	ipToLong := binary.BigEndian.Uint32(confLib.LocalIP.To4())
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
	return b.String()
}

func GetTraceId() (traceId string) {
	return calcTraceId(confLib.LocalIP.String())
}

func calcTraceId(ip string) (traceId string) {
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()

	b := bytes.Buffer{}
	netIP := net.ParseIP(ip)
	if netIP == nil {
		b.WriteString("00000000")
	} else {
		// IPv4 地址由 4 个字节（32 位）组成
		// 得到了一个长度为 8 的十六进制字符串，表示整个 IPv4 地址
		// 192.168.1.1 -> c0a80101
		b.WriteString(hex.EncodeToString(netIP.To4()))
	}
	// 取时间戳的前 32位，转化为 8位 16进制
	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
	// 取时间戳的前 16，转化为 4位 16进制
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	// 随机生成一个 24位的随机数 -> int32
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0") // 末两位标记来源,b0为go
	return b.String()
}


func InIPSliceStr(targetIP string, ipSliceStr []string) bool {
	if targetIP == "" || len(ipSliceStr) == 0 {
		return false
	}
	for _, ipSliceNode := range ipSliceStr {
		ip := net.ParseIP(ipSliceNode)
		if ip != nil {
			// ip
			if targetIP == ip.String() {
				return true
			}
		} else {
			// mask
			_, mask, err := net.ParseCIDR(ipSliceNode)
			if err != nil {
				fmt.Println("ParseCIDR error: ", err)
				continue
			}

			if mask.Contains(net.ParseIP(targetIP)) {
				return true
			}
		}
	}
	return false
}



