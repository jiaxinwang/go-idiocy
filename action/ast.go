package action

// GenerateDoc ...
func GenerateDoc(filename string) {
	// fset := token.NewFileSet()
	// path, _ := filepath.Abs(filename)
	// f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	// if err != nil {
	// 	logrus.WithError(err).Error()
	// 	return
	// }

	// for _, node := range f.Decls {
	// 	switch node.(type) {
	// 	case *ast.GenDecl:
	// 		genDecl := node.(*ast.GenDecl)
	// 		for _, spec := range genDecl.Specs {
	// 			switch spec.(type) {
	// 			case *ast.TypeSpec:
	// 				typeSpec := spec.(*ast.TypeSpec)
	// 				logrus.WithField("struct", typeSpec.Name.Name).Info()
	// 				switch typeSpec.Type.(type) {
	// 				case *ast.StructType:
	// 					structType := typeSpec.Type.(*ast.StructType)
	// 					for _, field := range structType.Fields.List {
	// 						ident, ok := field.Type.(*ast.Ident)
	// 						if ok {
	// 							fieldType := ident.Name
	// 							for _, name := range field.Names {
	// 								logrus.WithField("type", "*ast.Ident").WithField(name.Name, fieldType).Info()
	// 							}
	// 							continue
	// 						}
	// 						selectorExpr, ok := field.Type.(*ast.SelectorExpr)
	// 						if ok {
	// 							fieldType := selectorExpr.Sel.Name
	// 							for _, name := range field.Names {
	// 								logrus.WithField("type", "*ast.SelectorExpr").WithField(name.Name, fieldType).Info()
	// 							}
	// 							continue

	// 						}

	// 						starExpr, ok := field.Type.(*ast.StarExpr)
	// 						if ok {
	// 							selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
	// 							if !ok {
	// 								continue
	// 							}
	// 							fieldType := selectorExpr.Sel.Name
	// 							for _, name := range field.Names {
	// 								logrus.WithField("type", "*ast.StarExpr").WithField(name.Name, fieldType).Info()
	// 							}
	// 							continue
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}
