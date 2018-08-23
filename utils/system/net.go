package system

import (
	"fmt"
	"net"
	"errors"
)

func GetIntranetIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.New(fmt.Sprint("get worker ip has wrong: %s", err.Error()))
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("can't find the ipv4 address")
}
