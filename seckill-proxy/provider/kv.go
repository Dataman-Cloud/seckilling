// Package provider holds the different provider implementation.
package provider

import (
	"strings"
	"text/template"
	"time"

	"github.com/BurntSushi/ty/fun"
	"github.com/Dataman-Cloud/seckilling/seckill-proxy/types"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"strconv"
)

// Kv holds common configurations of key-value providers.
type Kv struct {
	baseProvider
	Endpoint  string
	Prefix    string
	StoreType store.Backend
	kvclient  store.Store
}

func (provider *Kv) provide(configurationChan chan<- types.ConfigMessage) error {
	kv, err := libkv.NewStore(
		provider.StoreType,
		[]string{provider.Endpoint},
		&store.Config{
			ConnectionTimeout: 30 * time.Second,
			Bucket:            "traefik",
		},
	)
	if err != nil {
		return err
	}
	if _, err := kv.List(""); err != nil {
		return err
	}
	provider.kvclient = kv
	if provider.Watch {
		stopCh := make(chan struct{})
		chanKeys, err := kv.WatchTree(provider.Prefix, stopCh)
		if err != nil {
			return err
		}
		go func() {
			for {
				<-chanKeys
				configuration := provider.loadConfig()
				if configuration != nil {
					configurationChan <- types.ConfigMessage{
						ProviderName:  string(provider.StoreType),
						Configuration: configuration,
					}
				}
				defer close(stopCh)
			}
		}()
	}
	configuration := provider.loadConfig()
	configurationChan <- types.ConfigMessage{
		ProviderName:  string(provider.StoreType),
		Configuration: configuration,
	}
	return nil
}

func (provider *Kv) loadConfig() *types.Configuration {
	templateObjects := struct {
		Prefix string
	}{
		provider.Prefix,
	}
	var KvFuncMap = template.FuncMap{
		"List":    provider.list,
		"Get":     provider.get,
		"GetBool": provider.getBool,
		"Last":    provider.last,
	}

	configuration, err := provider.getConfiguration("templates/kv.tmpl", KvFuncMap, templateObjects)
	if err != nil {
		log.Error(err)
	}
	return configuration
}

func (provider *Kv) list(keys ...string) []string {
	joinedKeys := strings.Join(keys, "")
	keysPairs, err := provider.kvclient.List(joinedKeys)
	if err != nil {
		log.Error("Error getting keys: ", joinedKeys, err)
		return nil
	}
	directoryKeys := make(map[string]string)
	for _, key := range keysPairs {
		directory := strings.Split(strings.TrimPrefix(key.Key, strings.TrimPrefix(joinedKeys, "/")), "/")[0]
		directoryKeys[directory] = joinedKeys + directory
	}
	return fun.Values(directoryKeys).([]string)
}

func (provider *Kv) get(keys ...string) string {
	joinedKeys := strings.Join(keys, "")
	keyPair, err := provider.kvclient.Get(joinedKeys)
	if err != nil {
		log.Error("Error getting key: ", joinedKeys, err)
		return ""
	} else if keyPair == nil {
		return ""
	}
	return string(keyPair.Value)
}

func (provider *Kv) getBool(keys ...string) bool {
	value := provider.get(keys...)
	b, err := strconv.ParseBool(string(value))
	if err != nil {
		log.Error("Error getting key: ", strings.Join(keys, ""), err)
		return false
	}
	return b
}

func (provider *Kv) last(key string) string {
	splittedKey := strings.Split(key, "/")
	return splittedKey[len(splittedKey)-1]
}
