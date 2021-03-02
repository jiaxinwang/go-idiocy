package parse

import (
	"idiocy/logger"
	"os"
	"path/filepath"
)

// Gin ...
func Gin(dir string) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				logger.S.Infow(path, "size", info.Size())
			}
			return err
		})
	if err != nil {
		return err
	}
	return nil
}
