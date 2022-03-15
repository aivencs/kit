package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var once sync.Once
var conf Config

const (
	defaultHost = "localhost:8500"
)

type Config interface {
	Get(path string) interface{}
	GetString(path string) string
	GetInt(path string) int
	GetBool(path string) bool
	GetFloat64(path string) float64
	GetStringMap(path string) map[string]interface{}
	GetStringMapString(path string) map[string]string
	GetStringMapStringSlice(path string) map[string][]string
	GetIntSlice(path string) []int
	GetStringSlice(path string) []string
	WatchRemoteConfig()
}

type ConsulConfig struct {
	Config *viper.Viper
}

type ConfigOptions struct {
	Auth        bool
	Username    string
	Password    string
	Host        string
	Application string
	Env         string
	Watch       bool
}

// config init
func InitConfig(name string, opt ConfigOptions) {
	once.Do(func() {
		switch name {
		case "consul":
			config, err := NewConsulConfig(opt)
			if err != nil {
				log.Fatal(err)
			}
			conf = config
		default:
			config, err := NewConsulConfig(opt)
			if err != nil {
				log.Fatal(err)
			}
			conf = config
		}

	})
}

// new config base consul
func NewConsulConfig(opt ConfigOptions) (Config, error) {
	if utf8.RuneCountInString(opt.Host) == 0 {
		opt.Host = defaultHost
	}
	if opt.Auth {
		os.Setenv("CONSUL_HTTP_AUTH", fmt.Sprintf("%s:%s", opt.Username, opt.Password))
	}
	vip := viper.New()
	name := fmt.Sprintf("%s/%s", opt.Application, opt.Env)
	vip.SetConfigType("yaml")
	vip.AddRemoteProvider("consul", opt.Host, name)
	err := vip.ReadRemoteConfig()
	if err != nil {
		return nil, err
	}
	if opt.Watch {
		go WatchRemoteConfig()
	}
	return &ConsulConfig{
		Config: vip,
	}, err
}

// watch remote config
func WatchRemoteConfig() {
	ticker := time.NewTicker(time.Second * time.Duration(3))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			conf.WatchRemoteConfig()
		default:
			continue
		}
	}
}

/*
function of consul config
*/
func (c *ConsulConfig) WatchRemoteConfig() {
	c.Config.WatchRemoteConfig()
}

func (c *ConsulConfig) Get(path string) interface{} {
	return c.Config.Get(path)
}

func (c *ConsulConfig) GetString(path string) string {
	return c.Config.GetString(path)
}

func (c *ConsulConfig) GetInt(path string) int {
	return c.Config.GetInt(path)
}

func (c *ConsulConfig) GetBool(path string) bool {
	return c.Config.GetBool(path)
}

func (c *ConsulConfig) GetFloat64(path string) float64 {
	return c.Config.GetFloat64(path)
}

func (c *ConsulConfig) GetStringMap(path string) map[string]interface{} {
	return c.Config.GetStringMap(path)
}

func (c *ConsulConfig) GetStringMapString(path string) map[string]string {
	return c.Config.GetStringMapString(path)
}

func (c *ConsulConfig) GetStringMapStringSlice(path string) map[string][]string {
	return c.Config.GetStringMapStringSlice(path)
}

func (c *ConsulConfig) GetIntSlice(path string) []int {
	return c.Config.GetIntSlice(path)
}

func (c *ConsulConfig) GetStringSlice(path string) []string {
	return c.Config.GetStringSlice(path)
}

/*
for caller
*/

func Get(path string) interface{} {
	return conf.Get(path)
}

func GetString(path string) string {
	return conf.GetString(path)
}

func GetInt(path string) int {
	return conf.GetInt(path)
}

func GetBool(path string) bool {
	return conf.GetBool(path)
}

func GetFloat64(path string) float64 {
	return conf.GetFloat64(path)
}

func GetStringMap(path string) map[string]interface{} {
	return conf.GetStringMap(path)
}

func GetStringMapString(path string) map[string]string {
	return conf.GetStringMapString(path)
}

func GetStringMapStringSlice(path string) map[string][]string {
	return conf.GetStringMapStringSlice(path)
}

func GetIntSlice(path string) []int {
	return conf.GetIntSlice(path)
}

func GetStringSlice(path string) []string {
	return conf.GetStringSlice(path)
}
