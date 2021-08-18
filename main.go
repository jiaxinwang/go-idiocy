package main

import (
	"encoding/json"
	"go/ast"
	"idiocy/logger"
	"idiocy/schema"
	"os"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/go-openapi/spec"
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
	projSchema := schema.NewSchema(c.String("dir"))
	projSchema.LoadSourceFiles()

	for _, v := range projSchema.SourceFile {
		v.GinIdents = make([]*ast.Ident, 0)
		v.ParseFile()
		v.BuildStacks()
		v.EnumerateStructAndGinVars()
		v.EnumerateGinHandles()
	}

	for _, v := range projSchema.SourceFile {
		v.ParseFile()
		v.EnumerateGinHandles()
	}

	for _, v := range schema.APIs {
		logger.S.Infof("%#v", v)
	}

	schema, _ := spec.Swagger20Schema()
	schema.Title = "major-tom API"
	schema.ID = `http://swagger.io/v2/schema.json#`
	logger.S.Infof("%#v", schema.SwaggerSchemaProps)

	doc := openapi2.T{}
	doc.BasePath = "http://127.0.0.1:51414"

	doc.Paths = make(map[string]*openapi2.PathItem)

	doc.Paths[`/health`] = &openapi2.PathItem{
		Get: &openapi2.Operation{
			Parameters: openapi2.Parameters{
				&openapi2.Parameter{
					Name:     "foo",
					Required: true,
				},
			},
		},
		Parameters: openapi2.Parameters{
			&openapi2.Parameter{
				Name:     "hahaha",
				Required: true,
			},
		},
	}

	contentJSON, _ := json.Marshal(doc)
	logger.S.Info(string(contentJSON))

	return nil
}
