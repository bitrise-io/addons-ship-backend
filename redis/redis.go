package redis

import (
	"os"

	"github.com/gomodule/redigo/redis"
)

// Client ...
type Client struct {
	conn redis.Conn
}

// New ...
func New() *Client {
	url := os.Getenv("REDIS_URL")
	return &Client{conn: newPool(url).Get()}
}

// Close ...
func (c *Client) Close() error {
	return c.conn.Close()
}

func newPool(url string) *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 1000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// Set ...
func (c *Client) Set(key, value string, ttl int) error {
	_, err := c.conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	if ttl > 0 {
		_, err := c.conn.Do("EXPIRE", key, ttl)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get ...
func (c *Client) Get(key string) (string, error) {
	value, err := redis.String(c.conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return value, nil
}
