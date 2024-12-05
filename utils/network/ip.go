package network

import (
	"fmt"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"net"
)

func GetLocalIP(ipStr string, subnet string) (net.IP, error) {
	if ipStr == "" || ipStr == "0.0.0.0" {
		return GetFirstActiveIP(subnet)
	}
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addresses {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.String() == ipStr {
				return ipNet.IP, nil
			}
		}
	}
	return nil, fmt.Errorf("no local IP address found")
}

func GetFirstActiveIP(subnet string) (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var network *net.IPNet
	if len(subnet) > 0 {
		_, network, err = net.ParseCIDR(subnet)
		if err != nil {
			return nil, err
		}
	}

	for _, iFace := range interfaces {
		// 过滤掉非活动的接口
		if iFace.Flags&net.FlagUp == 0 || iFace.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrArr, aErr := iFace.Addrs()
		if aErr != nil {
			return nil, aErr
		}

		for _, addr := range addrArr {
			ipNet, ok := addr.(*net.IPNet)
			if ok {
				//if ipNet.IP.To4() != nil {
				//	return ipNet.IP, nil
				//}
				if network != nil {
					if network.Contains(ipNet.IP) {
						return ipNet.IP, nil
					}
				} else {
					return ipNet.IP, nil
				}
			}
		} // end of for
	} // end of for

	return nil, fmt.Errorf("no active IP address found")
}

func GetActiveIps() ([]net.IP, error) {
	// 获取本机所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, iFace := range interfaces {
		// 过滤掉非活动的接口
		if iFace.Flags&net.FlagUp == 0 || iFace.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取该接口的地址
		addresses, aErr := iFace.Addrs()
		if aErr != nil {
			plog.Error("Error getting addresses:", pfield.Error(aErr))
			continue
		}

		for _, addr := range addresses {
			// 检查 addr 是否是 IP 地址
			if ipNet, ok := addr.(*net.IPNet); ok {
				//fmt.Println("Interface:", iFace.Name, "IP Address:", ipNet.IP.String())
				ips = append(ips, ipNet.IP)
			}
		}
	}
	return ips, nil
}
