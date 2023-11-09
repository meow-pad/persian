package network

import (
	"fmt"
	"net"
)

func GetFirstNonLoopbackIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iFace := range interfaces {
		addrArr, aErr := iFace.Addrs()
		if aErr != nil {
			return nil, aErr
		}

		for _, addr := range addrArr {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP, nil
				}
			}
		} // end of for
	} // end of for

	return nil, fmt.Errorf("no non-loopback IP address found")
}
