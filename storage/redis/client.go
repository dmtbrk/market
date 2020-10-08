package redis

import "github.com/go-redis/redis/v8"

func NewClientFromURL(url string) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}
