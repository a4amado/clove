package appservice

import repository "clove/internals/services/generatedRepo"

func (as *RegionsCtx) Update(arg repository.UpdateAppRegionsParams) error {
	queries := as.App.Queries
	if as.App.Tx != nil {
		queries = queries.WithTx(*as.App.Tx)
	}

	if as.App.Cache {
		// log in the cache
	}

	return queries.UpdateAppRegions(as.App.ReqCtx, arg)

}
