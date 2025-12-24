package appConsts

import "clove/internals/repository"

func GetAppBufferSize(app repository.AppType) int {
	switch app {
	case repository.AppTypePro:
		return 32 * 1024
	case repository.AppTypeStandard:
		return 16 * 1024
	default:
		return 4 * 1024
	}
}
