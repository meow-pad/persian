package name

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosNaming struct {
	naming_client.INamingClient
}

func (naming *NacosNaming) init(cCfg *constant.ClientConfig, sCfg []constant.ServerConfig) (err error) {
	naming.INamingClient, err = clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  cCfg,
		ServerConfigs: sCfg,
	})
	return
}
