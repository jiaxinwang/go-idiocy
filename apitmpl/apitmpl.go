package apitmpl

import (
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
)

var Doc openapi2.T

func init() {
	Doc = openapi2.T{}
	Doc.Swagger = "2.0"
	Doc.BasePath = `/`
	Doc.Definitions = make(map[string]*openapi3.SchemaRef)
}
