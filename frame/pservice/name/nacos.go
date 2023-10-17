package name

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func NewNacosNaming(cCfg *constant.ClientConfig, sCfg []constant.ServerConfig) (*NacosNaming, error) {
	nacosNaming := &NacosNaming{}
	if err := nacosNaming.init(cCfg, sCfg); err != nil {
		return nil, err
	}
	return nacosNaming, nil
}

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
