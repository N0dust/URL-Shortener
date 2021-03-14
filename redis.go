package main

import (
	"app/tools"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/mattheath/base62"
)

const (
	// URLIDKEY is global counter
	URLIDKEY = "next.url.id"
	// ShortlinkKey mapping the shortlink to the url
	ShortlinkKey = "shortlink:%s:url"
	// URLHashKey mapping the hash of the url to the shortlink
	URLHashKey = "urlhash:%s:url"
	// ShortlinkDetailKey mapping the shortlink to the detail of url
	ShortlinkDetailKey = "shortlink:%s:detail"
)

// RedisCli is ..
type RedisCli struct {
	Cli *redis.Client
}

// URLDetail is ..
type URLDetail struct {
	URL                 string        `json:"url"`
	CreatedAt           string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

// NewRedisCli is ..
func NewRedisCli(addr, passwd string, db int) *RedisCli {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}

	return &RedisCli{Cli: client}
}

// Shorten convert url to shortlink
func (r *RedisCli) Shorten(url string, exp int64) (string, error) {
	// convert url to sha1 hash
	h := tools.ToSha1(url)

	// fetch it if the url is cached
	d, err := r.Cli.Get(fmt.Sprintf(URLHashKey, h)).Result()
	if err == redis.Nil {
		// not exists，nothing to do
	} else if err != nil {
		return "", err
	} else {
		return d, nil
	}

	// increase the global counter
	err = r.Cli.Incr(URLIDKEY).Err()
	if err != nil {
		return "", nil
	}

	// encode global counter to base62
	id, err := r.Cli.Get(URLIDKEY).Int64()
	if err != nil {
		return "", err
	}
	eid := base62.EncodeInt64(id)

	// store the url against this encoded id
	err = r.Cli.Set(fmt.Sprintf(ShortlinkKey, eid), url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	// store the url against the hash of it
	err = r.Cli.Set(fmt.Sprintf(URLHashKey, h), eid, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	details, err := json.Marshal(&URLDetail{
		URL:                 url,
		CreatedAt:           time.Now().String(),
		ExpirationInMinutes: time.Duration(exp),
	})
	if err != nil {
		return "", err
	}

	// store the url detail against this encoded id
	err = r.Cli.Set(fmt.Sprintf(ShortlinkDetailKey, eid), details, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", nil
	}

	return eid, nil
}

// ShortlinkInfo returns the detail of the shortlink
func (r *RedisCli) ShortlinkInfo(eid string) (interface{}, error) {
	d, err := r.Cli.Get(fmt.Sprintf(ShortlinkDetailKey, eid)).Result()
	if err == redis.Nil {
		return "", errors.New("unknown short URL")
	} else if err != nil {
		return "", err
	} else {
		return d, nil
	}
}

// Unshorten convert shortlink to url
func (r *RedisCli) Unshorten(eid string) (string, error) {
	url, err := r.Cli.Get(fmt.Sprintf(ShortlinkKey, eid)).Result()
	if err == redis.Nil {
		return "", errors.New("unknown short URL")
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}
