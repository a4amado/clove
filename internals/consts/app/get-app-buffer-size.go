package appConsts

import repository "clove/internals/services/generatedRepo"

// GetAppBufferSize returns the buffer size in bytes for the given application type.
// Pro maps to 32 KB, Standard maps to 16 KB, and all other types map to 4 KB.
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
