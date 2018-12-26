package cli

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/9d77v/go-lib/clients/config"
	"github.com/9d77v/go-lib/clients/etcd"
)

//NewConfig get configs from directory
func NewConfig(configPath string) (map[string]*config.AppConfig, *config.DefaultConfig) {
	appConfigs := make(map[string]*config.AppConfig)
	defaultConfig := new(config.DefaultConfig)
	b, err := ioutil.ReadFile(configPath + "/conf.yml")
	if err != nil {
		log.Fatalln("error:reading configuration file")
	}
	err = yaml.Unmarshal(b, &defaultConfig)
	if err != nil {
		log.Fatalln("error:configuration file yml format parsing")
	}
	err = filepath.Walk(configPath+"/app", func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		appFile, appErr := ioutil.ReadFile(path)
		if appErr != nil {
			log.Fatalln("error:reading configuration file")
		}
		appConfig := new(config.AppConfig)
		err = yaml.Unmarshal(appFile, appConfig)
		if err != nil {
			log.Fatalln("error:configuration file yml format parsing")
		}
		serviceName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		appConfigs[serviceName] = config.NewAppConfig(defaultConfig, appConfig)
		return nil
	})
	if err != nil {
		log.Printf("filepath.Walk() returned %v\n", err)
	}
	return appConfigs, defaultConfig
}

//SyncConfig write changed config to etcd
func SyncConfig(dialTimeout time.Duration,
	requestTimeout time.Duration, configPath string) error {
	profile := filepath.Base(configPath)
	appConfigMap, defaultConfig := NewConfig(configPath)
	cli, err := etcd.NewClient(dialTimeout)
	defer cli.Close()
	if err != nil {
		log.Panicln("etcd connect failed,error:", err)
	}
	for k, v := range appConfigMap {
		updates := []struct {
			config interface{}
			key    string
		}{
			{v.DB, "db"},
			{v.Redis, "redis"},
			{v.Jaeger, "jaeger"},
			{v.Rabbitmq, "rabbitmq"},
		}
		for _, update := range updates {
			if update.config != nil {
				key := cli.GetEtcdKey(profile, k, update.key)
				value, err := json.Marshal(update.config)
				if err != nil {
					continue
				}
				err = cli.SyncKey(requestTimeout, key, value)
				if err != nil {
					return err
				}
			}
		}
	}
	globals := []struct {
		config interface{}
		key    string
	}{
		{defaultConfig.ExpressConfig, "express"},
		{defaultConfig.SMSConfig, "sms"},
	}
	for _, global := range globals {
		if global.config != nil {
			key := cli.GetEtcdKey(profile, "global", global.key)
			value, err := json.Marshal(global.config)
			if err != nil {
				continue
			}
			err = cli.SyncKey(requestTimeout, key, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
