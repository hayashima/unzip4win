package unzip4win

import (
	"path/filepath"
)

func IsLookLikeZipFile(zipPath string) bool {
	return filepath.Ext(zipPath) == "zip"
}
