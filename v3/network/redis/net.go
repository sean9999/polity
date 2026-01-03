package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const redisPort = 6379
const redisHost = "localhost"

type Network struct {
	rdb *redis.Client
}

func (n *Network) Up(ctx context.Context) error {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})
	n.rdb = rdb
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("can't find a redis server on %s:%d. %w", redisHost, redisPort, err)
	}
	return nil
}

func (n *Network) Spawn() *Node {
	return &Node{
		rdb: n.rdb,
	}
}
