package durable

import (
	"context"
	"fmt"

	"github.com/crossle/channel-father-mixin-bot/config"
	"github.com/go-pg/pg"
)

func OpenDatabaseClient(ctx context.Context) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.DatabaseHost, config.DatabasePort),
		User:     config.DatabaseUser,
		Password: config.DatabasePassword,
		Database: config.DatabaseName,
	})
	return db
}
