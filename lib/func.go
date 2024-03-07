package lib

import (
	"net"
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


func IsInArrayString(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}

	return false
}

