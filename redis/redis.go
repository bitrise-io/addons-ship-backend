package redis

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/bitrise-io/api-utils/utils"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

// Interface ...
type Interface interface {
	GetString(string) (string, error)
	GetInt64(key string) (int64, error)
	Set(string, interface{}, int) error
}

// Client ...
type Client struct {
	conn redis.Conn
}

// New ...
func New() *Client {
	return &Client{
		conn: NewPool(
			os.Getenv("REDIS_URL"),
			int(utils.GetInt64EnvWithDefault("REDIS_MAX_IDLE_CONNECTION", 50)),
			int(utils.GetInt64EnvWithDefault("REDIS_MAX_ACTIVE_CONNECTION", 1000)),
		).Get(),
	}
}

// Close ...
func (c *Client) Close() error {
	return c.conn.Close()
}

// NewPool ...
func NewPool(urlStr string, maxIdle, maxActive int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
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
func (c *Client) Set(key string, value interface{}, ttl int) error {
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

// GetString ...
func (c *Client) GetString(key string) (string, error) {
	value, err := redis.String(c.conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetInt64 ...
func (c *Client) GetInt64(key string) (int64, error) {
	value, err := redis.Int64(c.conn.Do("GET", key))
	if err != nil {
		return 0, err
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
