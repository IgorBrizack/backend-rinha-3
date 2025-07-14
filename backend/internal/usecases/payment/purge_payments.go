package commands

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type PurgePaymentsCommand struct {
	cacheClient *redis.Client
}

func NewPurgePaymentsCommand(cacheClient *redis.Client) *PurgePaymentsCommand {
	return &PurgePaymentsCommand{
		cacheClient: cacheClient,
	}
}

func (c *PurgePaymentsCommand) Execute(ctx context.Context) error {
	return c.cacheClient.FlushDB(ctx).Err()
}
