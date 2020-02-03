/**
* @Author: zhangjian@mioji.com
* @Date: 2019/9/10 下午6:40
 */
package utils

import (
	"net"
	"errors"
)

// GLocalHostIp 获取本地ip
func GetLocalHostIp() (string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), nil
					}
				}
			}
		}
	}
	return "", errors.New("获取ip失败")
}
