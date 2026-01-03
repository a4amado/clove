package appservice

import "github.com/jackc/pgx/v5/pgtype"

func (as *AppCtx) Delete() error {
	queries := as.BaseCtx.Queries
	if as.BaseCtx.Tx != nil {
		queries = queries.WithTx(*as.BaseCtx.Tx)
	}
	return queries.App_Delete(as.BaseCtx.ReqCtx, pgtype.UUID{
		Bytes: as.AppID,
		Valid: true,
	})
}
