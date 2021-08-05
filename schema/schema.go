package schema

import (
	"idiocy/logger"
	"io/ioutil"
	"sync"

	"github.com/jiaxinwang/err2"
	"golang.org/x/mod/modfile"
)

type Schema struct {
	ModulePath string
	cacheStore *sync.Map
}

func NewSchema(modFilePath string) *Schema {
	var err error
	defer err2.Handle(&err, func() {
		logger.S.Errorw("can't create Schema", "err", err)
	})

	ret := new(Schema)
	err2.Try(ret.ParseModules(modFilePath))
	return ret
}

func (schema *Schema) ParseModules(filename string) (err error) {
	defer err2.Return(&err)
	content := err2.Bytes.Try(ioutil.ReadFile(filename))
	schema.ModulePath = modfile.ModulePath(content)
	return nil
}
