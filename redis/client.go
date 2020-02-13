package redis

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// NewRedis creates a new redis connection
func NewRedis(c Config) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", c.Address, redis.DialPassword(c.Password))
	if err != nil {
		log.Errorf("Failed to dial %v for redis: %v", c.Address, err)
		return nil, err
	}
	//      if _, err := conn.Do("SELECT", db); err != nil {
	//        conn.Close()
	//        return nil, err
	//      }
	return conn, nil
}

// NewRedis creates a new redis connection
func NewRedisPool(c Config) *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) { return NewRedis(c) }, 3)
}
