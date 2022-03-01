package config

import (
	"log"
	"sync"

	"github.com/sgblanch/pathview-web/internal/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _once sync.Once
var _config *Config

type Config struct {
	Listen     string         `mapstructure:"listen"`
	Kegg       *Kegg          `mapstructure:"kegg"`
	DB         *db.DB         `mapstructure:"postgres"`
	RedisPool  *db.RedisPool  `mapstructure:"redis"`
	RedisStore *db.RedisStore `mapstructure:"session"`
	CSRFKey    string         `mapstructure:"csrf-key"`
	Google     *Auth          `mapstructure:"google"`
}

func Get() *Config {
	_once.Do(func() {
		log.SetFlags(log.LstdFlags | log.Lshortfile)

		_config = &Config{
			Kegg:       &Kegg{},
			DB:         &db.DB{},
			RedisPool:  &db.RedisPool{},
			RedisStore: &db.RedisStore{},
			Google:     &Auth{},
		}

		err := viper.Unmarshal(&_config)
		cobra.CheckErr(err)

		if _config.DB == nil {
			log.Printf("config.DB is null")
		}

		err = _config.DB.Open()
		cobra.CheckErr(err)

		if _config.RedisPool.Address != "" {
			_config.RedisPool.Open()

			err = _config.RedisStore.Open(_config.RedisPool)
			cobra.CheckErr(err)
		} else {
			log.Print("redis not configured, skipping")
		}
	})

	return _config
}

type Auth struct {
	ClientID     string `mapstructure:"client-id"`
	ClientSecret string `mapstructure:"client-secret"`
	Redirect     string `mapstructure:"redirect"`
}

type Kegg struct {
	BaseDir string `mapstructure:"base-dir"`
	KeggDir string `mapstructure:"kegg-dir"`
}
