package utils

import (
	"net"
)

func GetIpv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// 检查地址是否为IPv4
		ipNet, ok := address.(*net.IPNet)
		if !ok {
			continue
		}
		if ipNet.IP.IsLoopback() {
			// 跳过回环地址
			continue
		}
		ip := ipNet.IP.To4()
		return ip.String(), nil
	}
	return "", nil
}
