package AppReplication

import (
	redisPool "clove/internals/data/redispool"
	"clove/internals/repository"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// SaveApp saves the app info comming from the global replication, into the local redis instance
func (c *AppReplication) SaveApp(ctx context.Context, app repository.App) error {
	bytes, err := json.Marshal(app)
	if err != nil {
		return err
	}
	cmd := c.conn.Set(ctx, c.FormatAppKey(app.ID.Bytes), string(bytes), 0)
	_, err = cmd.Result()
	if err != nil {
		return err
	}
	return nil
}

// FetchApp fetches the app info from the local redis instance
func (c *AppReplication) FetchApp(ctx context.Context, appid uuid.UUID) (*repository.App, error) {

	result := c.conn.Get(ctx, c.FormatAppKey(appid))
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
func (c *AppReplication) FormatAppKey(appId uuid.UUID) string {
	return fmt.Sprintf("app:%s", appId.String())
}
