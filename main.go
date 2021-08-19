package main

import (
	"encoding/json"
	"go/ast"
	"idiocy/logger"
	"idiocy/schema"
	"os"
	"path"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/jiaxinwang/common/fs"
	"github.com/urfave/cli/v2"
)

func init() {
}

func main() {
	app := &cli.App{
		Name: "idiocy",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{`c`},
				Value:   "./config.toml",
			},
		},
		Commands: []*cli.Command{
			&genDoc,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

}

var genDoc = cli.Command{
	Name:  "gen",
	Usage: "generate doc",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "dir",
			Aliases: []string{`d`},
		},
		&cli.StringSliceFlag{
			Name:    "file",
			Aliases: []string{`f`},
		},
	},
	Action: run,
}

func run(c *cli.Context) error {
	schema.ProjSchema = schema.NewSchema(c.String("dir"))
	schema.ProjSchema.LoadSourceFiles()

	for _, v := range schema.ProjSchema.SourceFile {
		v.GinIdents = make([]*ast.Ident, 0)
		v.ParseFile()
		v.BuildStacks()
		v.EnumerateStructAndGinVars()

		// v.EnumerateGinHandles()
	}

	logger.S.Info("---> ", len(schema.ProjSchema.GinIdentifiers))

	for _, v := range schema.ProjSchema.SourceFile {
		v.ParseFile()
		v.EnumerateGinHandles()
	}

	for _, v := range schema.APIs {
		logger.S.Infof("%#v", v)
	}

	doc := openapi2.T{}
	doc.Swagger = "2.0"
	doc.BasePath = `/`

	doc.Paths = make(map[string]*openapi2.PathItem)

	doc.Paths[`/health`] = &openapi2.PathItem{
		Get: &openapi2.Operation{
			Parameters: openapi2.Parameters{
				&openapi2.Parameter{
					Name:     "foo",
					In:       "query",
					Required: true,
				},
			},
			Responses: map[string]*openapi2.Response{
				"200": {Description: "health"},
			},
		},
		Parameters: openapi2.Parameters{
			&openapi2.Parameter{
				Name:     "hahaha",
				In:       "query",
				Required: true,
			},
		},
	}

	contentJSON, _ := json.Marshal(doc)
	// logger.S.Info(string(contentJSON))
	swaggerDocFilepath := path.Join(path.Dir(c.String("dir")), `idiocy`, `swagger.json`)
	fs.Save(contentJSON, swaggerDocFilepath)
	return nil
}
