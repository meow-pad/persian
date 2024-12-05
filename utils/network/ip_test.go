package network

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestCIDR(t *testing.T) {
	should := require.New(t)
	ip, network, err := net.ParseCIDR("192.168.1.1/24")
	should.Nil(err)
	t.Logf("%s %s", ip, network)
}

func TestGetFirstActiveIP(t *testing.T) {
	should := require.New(t)
	ip, err := GetFirstActiveIP("192.168.71.1/24")
	should.Nil(err)
	t.Logf("%s", ip)
}
