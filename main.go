package main

import (
	"fmt"
	"go/ast"
	"idiocy/apitmpl"
	"idiocy/logger"
	"idiocy/schema"
	"os"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jiaxinwang/common/fs"
	jsoniter "github.com/json-iterator/go"
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

		// work
		v.EnumerateStructAndGinVars()
		// v.EnumerateGinHandles()

	}

	// for _, v := range schema.Structs {
	// 	logger.S.Infof("%#v", v)
	// }

	logger.S.Info("---> ", len(schema.ProjSchema.GinIdentifiers))
	for _, v := range schema.ProjSchema.GinIdentifiers {
		logger.S.Infof("%#v", v.Node)
	}
	logger.S.Info(len(schema.ProjSchema.GinIdentifiers[0].Calls))

	for _, v := range schema.ProjSchema.GinAPIs {
		logger.S.Infof("(%s)%s --> %v", v.Method, v.Path, v.Node)
	}

	for _, v := range schema.ProjSchema.SourceFile {
		v.ParseFile()
		v.EnumerateGinHandles()
	}

	doc := apitmpl.Doc
	// doc.Swagger = "2.0"
	// doc.BasePath = `/`

	doc.Paths = make(map[string]*openapi2.PathItem)

	for _, v := range schema.ProjSchema.GinAPIs {
		logger.S.Infof("(%s)%s --> %v", v.Method, v.Path, v.Node)
		// doc.Paths[v.Path] = &openapi2.PathItem{}
		pathItem := &openapi2.PathItem{}
		doc.Paths[v.Path] = pathItem
		opt := &openapi2.Operation{}
		switch v.Method {
		case "GET":
			pathItem.Get = opt
		case "POST":
			pathItem.Post = opt
		case "PUT":
			pathItem.Put = opt
		case "PATCH":
			pathItem.Patch = opt
		case "DELETE":
			pathItem.Delete = opt
		}

		if v.APIParam != nil && !strings.EqualFold(v.APIParam.StructName, "") {
			parts := strings.Split(v.APIParam.StructName, ".")
			base := parts[len(parts)-1]
			opt.Parameters = openapi2.Parameters{
				&openapi2.Parameter{
					In:          "body",
					Name:        "body",
					Required:    true,
					Description: "TODO:",
					Schema:      openapi3.NewSchemaRef(fmt.Sprintf("#/definitions/%s", base), nil),
				},
			}
		}

		opt.Responses = make(map[string]*openapi2.Response)

		for _, v := range v.APIResopnse {
			parts := strings.Split(v.StructName, ".")
			base := parts[len(parts)-1]
			opt.Responses[v.Code] = &openapi2.Response{
				Description: fmt.Sprintf("TODO: 缺少 %s 的描述", base),
				Schema:      openapi3.NewSchemaRef(fmt.Sprintf("#/definitions/%s", base), nil),
				// Ref:         fmt.Sprintf("#/definitions/%s", base),
			}
		}

	}

	// doc.Paths[`/health`] = &openapi2.PathItem{
	// 	Get: &openapi2.Operation{
	// 		Parameters: openapi2.Parameters{
	// 			&openapi2.Parameter{
	// 				Name:     "foo",
	// 				In:       "query",
	// 				Type:     "string",
	// 				Required: true,
	// 			},
	// 		},
	// 		Responses: map[string]*openapi2.Response{
	// 			"200": {Description: "health"},
	// 		},
	// 	},
	// 	Parameters: openapi2.Parameters{
	// 		&openapi2.Parameter{
	// 			Name:     "hahaha",
	// 			Type:     "string",
	// 			In:       "query",
	// 			Required: true,
	// 		},
	// 	},
	// }

	// contentJSON, _ := json.Marshal(doc)
	contentJSON, _ := jsoniter.MarshalIndent(doc, "", "    ")
	// logger.S.Info(string(contentJSON))
	swaggerDocFilepath := path.Join(path.Dir(c.String("dir")), `idiocy`, `swagger.json`)
	fs.Save(contentJSON, swaggerDocFilepath)
	return nil
}
