package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/mattheath/base62"
)

const (
	// URLIDKEY for global counter
	URLIDKEY = "next.url.id"
	//ShortlinkKey shortlink of url key
	ShortlinkKey = "shortlink:%s:url"
	//URLHashKey the hash of the url to the shorlink
	URLHashKey = "urlhash:%s:url"
	//ShortlinkDetailKey the shortlink to the detail of url
	ShortlinkDetailKey = "shortlink:%s:detail"
)

//RedisCli redis client
type RedisCli struct {
	Cli *redis.Client
}

//URLDetail ...
type URLDetail struct {
	URL                 string        `json:"url" `
	CreateAt            string        `json:"create_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

//NewRedisCli new a redis cli
func NewRedisCli(addr string, passwd string, db int) *RedisCli {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	if _, err := c.Ping().Result(); err != nil {
		panic(err)
	}

	return &RedisCli{Cli: c}
}

//Shorten  the method
func (r *RedisCli) Shorten(url string, exp int64) (string, error) {
	h := ToSha1(url)

	//Get URLHashKey
	d, err := r.Cli.Get(fmt.Sprintf(URLHashKey, h)).Result()
	if err == redis.Nil {
		//如果存在，不做任何事情
	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
			//過期了 不需要做任何事情
		} else {
			return d, nil
		}
	}

	//increate the global counter
	err = r.Cli.Incr(URLIDKEY).Err()
	if err != nil {
		return "", err
	}

	//encode global counter to base62
	id, err := r.Cli.Get(URLIDKEY).Int64()
	if err != nil {
		return "", err
	}

	eid := base62.EncodeInt64(id)

	//store the url against this encoded id
	err = r.Cli.Set(fmt.Sprintf(ShortlinkKey, eid), url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	//store the url against the hash of it
	err = r.Cli.Set(fmt.Sprintf(URLHashKey, h), eid, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	detail, err := json.Marshal(
		&URLDetail{
			URL:                 url,
			CreateAt:            time.Now().String(),
			ExpirationInMinutes: time.Duration(exp),
		},
	)
	if err != nil {
		return "", err
	}

	//store the url  detail to  redis
	err = r.Cli.Set(fmt.Sprintf(ShortlinkDetailKey, eid), detail, time.Minute*time.Duration(exp)).Err()

	if err != nil {
		return "", err
	}
	return eid, nil
}

//ShortLinkInfo return the detail of the short link
func (r *RedisCli) ShortLinkInfo(eid string) (interface{}, error) {
	d, err := r.Cli.Get(fmt.Sprintf(ShortlinkDetailKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{Code: 404, Err: fmt.Errorf("Unknow short URL ")}
	} else if err != nil {
		return "", err
	} else {
		return d, nil
	}
}

//Unshorten the short url to url
func (r *RedisCli) Unshorten(eid string) (string, error) {
	url, err := r.Cli.Get(fmt.Sprintf(ShortlinkKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{Code: 404, Err: fmt.Errorf("Unknow short URL ")}
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}
