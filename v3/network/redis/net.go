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

func ServerIsRunning(ctx context.Context) bool {
	n := new(Network)
	err := n.Up(ctx)
	if err != nil {
		return false
	}
	n.Down()
	return true
}

func (n *Network) Down() error {
	return n.rdb.Close()
}

func (n *Network) Up(ctx context.Context) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("can't find a redis server on %s:%d. %w", redisHost, redisPort, err)
	}
	n.rdb = rdb
	return nil
}

func (n *Network) Spawn() *Node {
	return &Node{
		rdb: n.rdb,
	}
}
