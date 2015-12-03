package provider

import (
	"github.com/Dataman-Cloud/seckilling/seckill-proxy/types"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/zookeeper"
)

// Zookepper holds configurations of the Zookepper provider.
type Zookepper struct {
	Kv
}

// Provide allows the provider to provide configurations to traefik
// using the given configuration channel.
func (provider *Zookepper) Provide(configurationChan chan<- types.ConfigMessage) error {
	provider.StoreType = store.ZK
	zookeeper.Register()
	return provider.provide(configurationChan)
}
