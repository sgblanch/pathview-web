package db

import (
	"time"

	rs "github.com/gin-contrib/sessions/redis"
	"github.com/gomodule/redigo/redis"
	"github.com/sgblanch/pathview-web/internal/util"
	"github.com/spf13/cobra"
)

type RedisPool struct {
	*redis.Pool
	Address string `mapstructure:"address"`
}

func (p *RedisPool) Open() {
	p.Pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", p.Address)
		},
	}
}

type RedisStore struct {
	rs.Store
	AuthenticationKey string `mapstructure:"auth-key"`
	EncryptionKey     string `mapstructure:"enc-key"`
}

func (p *RedisStore) Open(pool *RedisPool) error {
	authkey, err := util.ParseKey(p.AuthenticationKey, 64)
	cobra.CheckErr(err)
	enckey, err := util.ParseKey(p.EncryptionKey, 32)
	cobra.CheckErr(err)

	p.Store, err = rs.NewStoreWithPool(pool.Pool, authkey, enckey)

	return err
}
