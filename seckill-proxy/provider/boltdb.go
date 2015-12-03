package provider

import (
	"github.com/Dataman-Cloud/seckilling/seckill-proxy/types"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
)

// BoltDb holds configurations of the BoltDb provider.
type BoltDb struct {
	Kv
}

// Provide allows the provider to provide configurations to traefik
// using the given configuration channel.
func (provider *BoltDb) Provide(configurationChan chan<- types.ConfigMessage) error {
	provider.StoreType = store.BOLTDB
	boltdb.Register()
	return provider.provide(configurationChan)
}
