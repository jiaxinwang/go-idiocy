package parse

import (
	"os"
	"path/filepath"

	"github.com/jiaxinwang/err2"
)

// Gin ...
func Gin(dir string) error {
	var err error
	defer err2.Return(&err)
	err = filepath.Walk(dir,
		func(curDir string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				// fset := token.NewFileSet()
				// path, _ := filepath.Abs(curDir)
				// f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
				// if err != nil {
				// 	return err
				// }

				// b := err2.Bytes.Try(fs.ReadBytes(curDir))

				// code := &schema.SourceCode{
				// 	fset:     fset,
				// 	astFile:  f,
				// 	src:      b,
				// 	filename: curDir,
				// 	decls:    make(map[string]string),
				// }
				// code.find()

			}
			return err
		})
	if err != nil {
		return err
	}
	return nil
}
