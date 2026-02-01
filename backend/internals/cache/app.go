package cache

import (
	"clove/internals/data/valkeyPool"
	repository "clove/internals/services/generatedRepo"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

type AppCache struct{}

var appRef = AppCache{}

func Apps() *AppCache {
	return &AppCache{}
}

func (a *AppCache) formatAppKey(appID uuid.UUID) string {
	return fmt.Sprintf("apps:%s", appID.String())
}

func (a *AppCache) Get(ctx context.Context, appID uuid.UUID) (*repository.App, error) {
	valkeyClient := valkeyPool.Client(valkeyPool.ValkeyStore)

	bytes, err := valkeyClient.DoCache(ctx, valkeyClient.B().Get().Key(a.formatAppKey(appID)).Cache(), time.Second*30).AsBytes()
	if err != nil {

		return nil, err
	}

	fetchedApp := repository.App{}
	if err = json.Unmarshal(bytes, &fetchedApp); err != nil {
		return nil, err
	}
	return &fetchedApp, nil
}
func (a *AppCache) Set(ctx context.Context, app repository.App) error {
	appBytes, err := json.Marshal(app)
	if err != nil {
		return err
	}
	valkeyClient := valkeyPool.Client(valkeyPool.ValkeyStore)

	err = valkeyClient.Do(ctx, valkeyClient.B().Set().Key(a.formatAppKey(app.ID.Bytes)).Value(valkey.BinaryString(appBytes)).Ex(time.Second*30).Build()).Error()
	if err != nil {
		return err
	}
	return nil
}
