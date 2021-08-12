package schema

import (
	"idiocy/logger"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jiaxinwang/err2"
	"golang.org/x/mod/modfile"
)

var APIs []API

func init() {
	APIs = make([]API, 0)
}

type Schema struct {
	ModuleFilePath string
	ModulePath     string
	cacheStore     *sync.Map
	SourceFile     []SourceFile
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
	schema.ModuleFilePath = filename
	schema.ModulePath = modfile.ModulePath(content)
	return nil
}

var skips = []string{
	".git",
	"vendor",
}

func (schema *Schema) LoadSourceFiles() (err error) {
	sourceRoot := path.Dir(schema.ModuleFilePath)
	err = filepath.Walk(sourceRoot, func(filename string, info os.FileInfo, err error) error {
		if info.IsDir() {
			for _, v := range skips {
				if strings.EqualFold(filename, path.Join(sourceRoot, v)) {
					return filepath.SkipDir
				}
			}
		} else {
			if strings.EqualFold(filepath.Ext(filename), ".go") {
				logger.S.Debug(filename)
				schema.SourceFile = append(schema.SourceFile, SourceFile{FullPath: filename, Path: filepath.Dir(filename)})
			}
		}
		return nil
	})
	return nil
}
