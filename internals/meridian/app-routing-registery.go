package meridian

import (
	redisPool "clove/internals/data/redispool"
	"clove/internals/repository"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

// SaveApp saves the app info comming from the global replication, into the local redis instance
func (c *Meridian) SaveApp(ctx context.Context, app repository.App) error {
	bytes, err := json.Marshal(app)
	if err != nil {
		return err
	}
	cmd := c.RedisStoreConn.Set(ctx, c.FormatAppKey(app.ID), string(bytes), 0)
	_, err = cmd.Result()
	if err != nil {
		return err
	}
	return nil
}

// FetchApp fetches the app info from the local redis instance
func (c *Meridian) FetchApp(ctx context.Context, appid pgtype.UUID) (*repository.App, error) {

	result := c.RedisStoreConn.Get(ctx, c.FormatAppKey(appid))
	byts, err := result.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, redisPool.ErrCacheMiss
		}
		return nil, err
	}

	if len(byts) == 0 {
		return nil, redisPool.ErrCacheMiss
	}
	fetchedApp := repository.App{}
	if err = json.Unmarshal(byts, &fetchedApp); err != nil {
		return nil, err
	}
	return &fetchedApp, nil

}
func (c *Meridian) FormatAppKey(appId pgtype.UUID) string {
	return fmt.Sprintf("app:%s", appId.String())
}
func (c *Meridian) FormatChannelKey(app uuid.UUID, channel string) string {
	return fmt.Sprintf("%s:%s", app.String(), channel)
}
