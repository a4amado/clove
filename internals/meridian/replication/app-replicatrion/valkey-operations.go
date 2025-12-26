package AppReplication

import (
	"clove/internals/data/valkeyPool"
	"clove/internals/repository"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

// SaveApp saves the app info comming from the global replication, into the local valkey instance
func (c *AppReplication) SaveApp(ctx context.Context, app repository.App) error {
	bytes, err := json.Marshal(app)
	if err != nil {
		return err
	}
	err = c.conn.Do(ctx, c.conn.B().Set().Key(c.FormatAppKey(app.ID.Bytes)).Value(valkey.BinaryString(bytes)).Build()).Error()
	if err != nil {
		return err
	}
	return nil
}

// FetchApp fetches the app info from the local valkey instance
func (c *AppReplication) FetchApp(ctx context.Context, appid uuid.UUID) (*repository.App, error) {

	byts, err := c.conn.Do(ctx, c.conn.B().Get().Key(c.FormatAppKey(appid)).Build()).AsBytes()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return nil, valkeyPool.ErrCacheMiss
		}
		return nil, err
	}

	if len(byts) == 0 {
		return nil, valkeyPool.ErrCacheMiss
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
