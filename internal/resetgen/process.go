// Package resetgen - пакет для генерации файлов с методами reset для структур.
package resetgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"path/filepath"
	"strings"

	"github.com/konkovaanna23/shortener/internal/file"
	"golang.org/x/tools/go/packages"
)

// ProcessPackage - основна функция для генерации файлов
func ProcessPackage(pkg *packages.Package) {
	var structs []StructInfo
	fmt.Println("Генерируем")
	for i, file := range pkg.Syntax {
		f := pkg.CompiledGoFiles[i]

		if strings.HasSuffix(f, ".gen.go") {
			continue
		}

		currentPkgPath := pkg.PkgPath

		fmt.Println(currentPkgPath)
		ast.Inspect(file, func(n ast.Node) bool {
			gen, ok := n.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				return true
			}

			if gen.Doc == nil {
				return true
			}

			for _, comment := range gen.Doc.List {
				if strings.TrimSpace(comment.Text) == "// generate:reset" {
					for _, spec := range gen.Specs {
						ts, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}

						if _, ok := ts.Type.(*ast.StructType); ok {
							obj := pkg.Types.Scope().Lookup(ts.Name.Name)
							if obj != nil {
								info := ExtractStructInfo(ts.Name.Name, pkg.Types.Scope().Lookup(ts.Name.Name), currentPkgPath)
								structs = append(structs, info)
							}
						}
					}
					return true

				}
			}

			return true
		})
	}
	if len(structs) > 0 {
		if err := generateResetFile(pkg, structs); err != nil {
			log.Printf("ошибка генерации reset.gen.go in %s: %v", pkg.PkgPath, err)
		}
	}
}

func generateResetFile(pkg *packages.Package, structs []StructInfo) error {

	if len(pkg.CompiledGoFiles) == 0 {
		return nil
	}

	dir := filepath.Dir(pkg.CompiledGoFiles[0])

	filename := filepath.Join(dir, "reset.gen.go")

	data, err := generateResetData(pkg.Name, structs)
	if err != nil {
		return err
	}

	err = file.SaveToFile(filename, data)
	if err != nil {
		return err
	}
	return nil
}

func generateResetData(packageName string, structs []StructInfo) ([]byte, error) {

	allImports := make(map[string]string)
	for _, s := range structs {
		for imp := range s.Imports {
			allImports[imp] = ""
		}
	}

	var buf bytes.Buffer
	data := struct {
		Package string
		Imports map[string]string
		Structs []StructInfo
	}{
		Package: packageName,
		Structs: structs,
		Imports: allImports,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("ошибка форматирования: %w\n%s", err, buf.String())
	}

	return formatted, nil
}
