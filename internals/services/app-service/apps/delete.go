package appservice

import "github.com/jackc/pgx/v5/pgtype"

func (as *AppServiceCtx) Delete() error {
	queries := as.App.Queries
	if as.App.Tx != nil {
		queries = queries.WithTx(*as.App.Tx)
	}
	return queries.DeleteApp(as.App.ReqCtx, pgtype.UUID{
		Bytes: as.AppId,
		Valid: true,
	})
}
