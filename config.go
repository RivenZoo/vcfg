package vcfg

import (
	"github.com/spf13/viper"
	"os"
	"time"
	"log"
)

type VConfig struct {
	viper *viper.Viper
}

func NewVConfig() (*VConfig, error) {
	return &VConfig{
		viper: viper.New(),
	}, nil
}

// Unmarshal decode to struct.
// dest: pointer to struct
func (vc *VConfig) Unmarshal(dest interface{}) error {
	return vc.viper.Unmarshal(dest)
}

// ReadConfig read config into VConfig.
// Support config file format: json|yaml
// cfgFile: config file path, support file ext .json | .yaml
func (vc *VConfig) ReadConfig(cfgFile string) error {
	f, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return vc.viper.ReadConfig(f)
}

// UnmarshalConfig read config file into VConfig then decode it to struct.
func (vc *VConfig) UnmarshalConfig(cfgFile string, dest interface{}) (err error) {
	err = vc.ReadConfig(cfgFile)
	if err == nil {
		err = vc.Unmarshal(dest)
	}
	return
}

type RemoteProvider string

const (
	Etcd   RemoteProvider = "etcd"
	Consul RemoteProvider = "consul"
)

type RemoteConfigOption struct {
	// Etcd | Consul
	Provider RemoteProvider `json:"provider"`
	Endpoint string         `json:"endpoint"`
	// Path: config store path.
	Path string `json:"path"`
	// json | yaml
	ConfigType string `json:"config_type"`
}

// WatchRemoteConfig update config from remote provider for every 5 seconds.
// Return update signal channel or error.
// Example:
// ch, err := vc.WatchRemoteConfig(opt)
// for sig := range ch {
// 	 vc.Unmarshal(dest)
// }
func (vc *VConfig) WatchRemoteConfig(opt *RemoteConfigOption) (<-chan struct{}, error) {
	vc.viper.AddRemoteProvider(string(opt.Provider), opt.Endpoint, opt.Path)
	vc.viper.SetConfigType(opt.ConfigType)
	err := vc.viper.ReadRemoteConfig()
	if err != nil {
		return nil, err
	}
	updateSignal := make(chan struct{}, 1)
	go func() {
		for {
			time.Sleep(time.Second * 5)
			err = vc.viper.WatchRemoteConfig()
			if err != nil {
				log.Printf("[Error] WatchRemoteConfig error: %v\n", err)
				continue
			}
			updateSignal <- struct{}{}
		}
	}()

	return updateSignal, nil
}
