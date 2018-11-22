package vcfg_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"log"

	"github.com/RivenZoo/vcfg"
	"fmt"
)

const (
	testConfigFile = "cfg.json"
)

type testConfig struct {
	Address string
	Port    int
}

func initTestConfig() (clean func(), err error) {
	cfg := &testConfig{
		"localhost",
		80,
	}
	var data []byte
	data, err = json.Marshal(cfg)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(testConfigFile, data, 0644)
	if err != nil {
		return
	}
	clean = func() {
		os.Remove(testConfigFile)
	}
	return
}

func ExampleVConfig_ReadAndUnmarshal() {
	clean, err := initTestConfig()
	if err != nil {
		log.Panicf("initTestConfig error: %v", err)
	}
	defer clean()

	cfg := vcfg.NewVConfig()
	err = cfg.ReadConfig(testConfigFile)
	if err != nil {
		log.Panicf("ReadConfig error: %v", err)
	}
	fmt.Println("ReadConfig success")

	testCfg := &testConfig{}
	err = cfg.Unmarshal(testCfg)
	if err != nil {
		log.Panicf("Unmarshal error: %v", err)
	}
	if testCfg.Address == "" {
		log.Panicf("read address fail")
	}
	fmt.Printf("Unmarshal success, address: %s, port: %d", testCfg.Address, testCfg.Port)

	// Output:
	// ReadConfig success
	// Unmarshal success, address: localhost, port: 80
}

func ExampleVConfig_WatchRemoteConfig() {
	clean, err := initTestConfig()
	if err != nil {
		log.Panicf("initTestConfig error: %v", err)
	}
	defer clean()

	testCfg := &testConfig{}
	cfg := vcfg.NewVConfig()
	ch, err := cfg.WatchRemoteConfig(&vcfg.RemoteConfigOption{
		Provider:   vcfg.Consul,
		Endpoint:   "127.0.0.1:8500",
		Path:       "my_config_path",
		ConfigType: "json",
	})
	if err != nil {
		log.Panicf("WatchRemoteConfig error: %v", err)
	}
	fmt.Println("WatchRemoteConfig success")

	go func() {
		for {
			<-ch
			// updated
			cfg.Unmarshal(testCfg)
		}
	}()

	// Output:
	// WatchRemoteConfig success
}
