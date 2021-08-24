package schema

import (
	"go/ast"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jiaxinwang/go-idiocy/logger"

	"github.com/jiaxinwang/err2"
	"golang.org/x/mod/modfile"
)

func init() {
}

var ProjSchema *Schema

type Schema struct {
	ModuleFilePath string
	ModulePath     string
	cacheStore     *sync.Map
	SourceFile     []SourceFile
	GinIdentifiers []*GinIdentifier
	GinAPIs        []*GinAPI
}

func NewSchema(modFilePath string) *Schema {
	var err error
	defer err2.Handle(&err, func() {
		logger.S.Errorw("can't create Schema", "err", err)
	})

	ret := new(Schema)
	ret.GinIdentifiers = make([]*GinIdentifier, 0)
	ret.GinAPIs = make([]*GinAPI, 0)
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
				schema.SourceFile = append(schema.SourceFile, SourceFile{FullPath: filename, Path: filepath.Dir(filename)})
			}
		}
		return nil
	})
	return nil
}

func (s *Schema) GinIdentifierWithFileIdent(source *SourceFile, callExpr *ast.CallExpr) *GinIdentifier {
	another := GinIdentifier{Source: source, InstancingCall: callExpr}
	for _, v := range s.GinIdentifiers {
		if v.Equal(&another) {
			return v
		}
	}
	return nil
}

func (s *Schema) AddGinIdentifier(source *SourceFile, ident *GinIdentifier) bool {
	for _, v := range s.GinIdentifiers {
		if ident.Equal(v) {
			return false
		}
	}
	s.GinIdentifiers = append(s.GinIdentifiers, ident)
	return true
}
