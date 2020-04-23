/**
* @Author: zhangjian@mioji.com
* @Date: 2019/9/10 下午6:40
 */
package utils

import (
	"encoding/json"
	"math/rand"
	"net"
	"errors"
	"time"
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

// json化数据
func JsonEncode(v interface{}) (string, error) {
	jsonByte, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonByte), nil
}

//产生随机字符串，不长于32位
func GetNonceStr() string {
	//随机种子
	rand.Seed(time.Now().UnixNano())
	//随机选取的字符串
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	str := ""
	//生成字符串
	for i := 0; i < 32; i++ {
		s := rand.Intn(len(chars) - 1)
		temp := string(chars[s])
		str += temp
	}
	//返回值字符串
	return str
}