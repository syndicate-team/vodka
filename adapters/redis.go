package adapters

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// Redis - redis DB mapper
type Redis struct {
	config Config
	client *redis.Client
}

/*
NewRedis - adapter constructor
*/
func NewRedis(config Config) *Redis {
	return &Redis{
		config: config,
	}
}

// Connect - connecting to repository
func (r *Redis) Connect() error {
	return r.connect()
}

/*
Get - getting data by key
*/
func (r *Redis) Get(key string) ([]byte, error) {
	err := r.checkConnection()
	if err != nil {
		return nil, err
	}
	// TODO: check for redis.Nil: if err == redis.Nil {
	return r.client.Get(key).Bytes()
}

/*
Del - deleting data by key
*/
func (r *Redis) Del(key string) error {
	err := r.checkConnection()
	if err != nil {
		return err
	}
	return r.client.Del(key).Err()
}

/*
Set - setting key with value and expiration time
*/
func (r *Redis) Set(key string, value interface{}, exp time.Duration) error {
	err := r.checkConnection()
	if err != nil {
		return err
	}
	return r.client.Set(key, value, exp).Err()
}

/*
SetJSON - setting key with value and expiration time that will be saved as JSON
*/
func (r *Redis) SetJSON(key string, value interface{}, exp time.Duration) error {
	str, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(key, str, exp)
}

func (r *Redis) connect() error {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.getAddr(),
		Password: r.config.Password, // no password set
		DB:       0,                 // use default DB
	})
	_, err := r.client.Ping().Result()
	return err
}

func (r *Redis) checkConnection() error {
	if r.client == nil {
		return r.connect()
	}

	fmt.Printf("Redis connection: %+v\n", r.client)
	return nil
}

func (r *Redis) getAddr() string {
	return r.config.Host + ":" + strconv.Itoa(r.config.Port)
}
