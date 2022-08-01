package config

import (
	"context"
	"fmt"
	"harago/common"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server Server `yaml:"server"`
	Nats   Nats   `yaml:"nats"`
	DB     DB     `yaml:"db"`
	Harbor Harbor `yaml:"harbor"`
	Etcd   Etcd   `yaml:"etcd"`
}

type Server struct {
	Port int `yaml:"port" env-default:"5678"`
}

type Nats struct {
	Servers  []string      `yaml:"servers"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Timeout  time.Duration `yaml:"timeout"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Timezone string `yaml:"timezone"`
}

type Harbor struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Etcd struct {
	Endpoints []string `env:"HARAGO_ETCD_ENDPOINTS" yaml:"endpoints"`
	Username  string   `env:"HARAGO_ETCD_USERNAME" yaml:"username"`
	Password  string   `env:"HARAGO_ETCD_PASSWORD" yaml:"password"`
	ConfigKey string   `env:"HARAGO_ETCD_CONFIG_KEY" yaml:"-"`
}

func NewConfig(opts *Options) (*Config, error) {
	var cfg Config

	if opts.Etcd {
		if err := fromEtcd(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	}

	if err := fromFile(&cfg, opts.ConfigFile); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func fromEtcd(cfg *Config) error {
	log.Println("read config from etcd")
	var etcdCfg Etcd
	if err := cleanenv.ReadEnv(&etcdCfg); err != nil {
		return fmt.Errorf("failed to read env variables; %w", err)
	}

	var cli *clientv3.Client
	var err error
	if etcdCfg.Username != "" {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   etcdCfg.Endpoints,
			DialTimeout: common.DefaultTimeout,
			Username:    etcdCfg.Username,
			Password:    etcdCfg.Password,
		})
	} else {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   etcdCfg.Endpoints,
			DialTimeout: common.DefaultTimeout,
		})
	}
	if err != nil {
		return fmt.Errorf("failed to connect etcd; %w", err)
	}
	defer func() {
		if err = cli.Close(); err != nil {
			log.Printf("failed to close etcd client cleany; %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultTimeout)
	defer cancel()

	resp, err := cli.Get(ctx, etcdCfg.ConfigKey)
	if err != nil {
		return fmt.Errorf("failed to get kv; %w", err)
	}

	if len(resp.Kvs) == 0 {
		return fmt.Errorf("failed to find value from key: %s", etcdCfg.ConfigKey)
	}

	if err := yaml.Unmarshal(resp.Kvs[0].Value, cfg); err != nil {
		return fmt.Errorf("failed to unmarshal value; %w", err)
	}
	return nil
}

func fromFile(cfg *Config, configPath string) error {
	log.Println("read config from file")
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return fmt.Errorf("failed to read config; %w", err)
	}
	return nil
}
