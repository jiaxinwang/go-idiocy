package doc

import (
	"fmt"
	"go/ast"
	"idiocy/apitmpl"
	"idiocy/logger"
	"idiocy/schema"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jiaxinwang/common/fs"
	jsoniter "github.com/json-iterator/go"
)

func Analyse(dir string) {
	schema.ProjSchema = schema.NewSchema(dir)
	schema.ProjSchema.LoadSourceFiles()

	for _, v := range schema.ProjSchema.SourceFile {
		v.GinIdents = make([]*ast.Ident, 0)
		v.ParseFile()
		v.BuildStacks()
		// work
		v.EnumerateStructAndGinVars()
		// v.EnumerateGinHandles()
	}

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
	doc.Paths = make(map[string]*openapi2.PathItem)

	for _, v := range schema.ProjSchema.GinAPIs {
		logger.S.Infof("(%s)%s --> %v", v.Method, v.Path, v.Node)
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

		if len(v.APIResopnse) == 0 {
			opt.Responses["200"] = &openapi2.Response{
				Description: fmt.Sprintf("无返回 body"),
			}
		}

		for _, v := range v.APIResopnse {
			parts := strings.Split(v.StructName, ".")
			base := parts[len(parts)-1]
			opt.Responses[v.Code] = &openapi2.Response{
				Description: fmt.Sprintf("TODO: 缺少 %s 的描述", base),
				Schema:      openapi3.NewSchemaRef(fmt.Sprintf("#/definitions/%s", base), nil),
			}
		}

	}

	contentJSON, _ := jsoniter.MarshalIndent(doc, "", "    ")
	swaggerDocFilepath := path.Join(path.Dir(dir), `idiocy`, `swagger.json`)
	fs.Save(contentJSON, swaggerDocFilepath)
}
