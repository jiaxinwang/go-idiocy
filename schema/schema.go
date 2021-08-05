package schema

import (
	"idiocy/logger"
	"io/ioutil"
	"sync"

	"golang.org/x/mod/modfile"
)

type Schema struct {
	ModulePath string
	cacheStore *sync.Map
}

func (schema *Schema) ParseModules(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.S.Info(err)
		return err
	}
	schema.ModulePath = modfile.ModulePath(content)
	return nil
}
