package redis

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
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

func newPool(urlStr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   50,
		MaxActive: 1000,
		Dial: func() (redis.Conn, error) {
			url, err := DialURL(urlStr)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			pass, err := DialPassword(urlStr)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			c, err := redis.Dial("tcp", url, redis.DialPassword(pass))
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return c, nil
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

// DialURL ...
func DialURL(urlToParse string) (string, error) {
	if !strings.HasPrefix(urlToParse, "redis://") {
		urlToParse = "redis://" + urlToParse
	}
	url, err := url.Parse(urlToParse)
	if err != nil {
		return "", err
	}
	if url.Hostname() == "" {
		return "", errors.New("Invalid hostname")
	}
	if url.Port() == "" {
		return "", errors.New("Invalid port")
	}
	return fmt.Sprintf("%s:%s", url.Hostname(), url.Port()), nil
}

// DialPassword ...
func DialPassword(urlToParse string) (string, error) {
	if !strings.HasPrefix(urlToParse, "redis://") {
		urlToParse = "redis://" + urlToParse
	}
	url, err := url.Parse(urlToParse)
	if err != nil {
		return "", err
	}
	pass, _ := url.User.Password()
	return pass, nil
}
