// Package exitanalyzer - пакет для анализа использования os.Exit().
package exitanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitAnalyzer = &analysis.Analyzer{
	Name: "exitanalyzer",
	Doc:  "check for exit() calls",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {

		isMainPkg := file.Name.Name == "main"

		var currentFunc *ast.FuncDecl

		ast.Inspect(file, func(n ast.Node) bool {
			if decl, ok := n.(*ast.FuncDecl); ok {
				currentFunc = decl
				return true
			}

			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				xIdent, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if xIdent.Name == "os" && sel.Sel.Name == "Exit" {
					if !isMainPkg {
						pass.Reportf(call.Lparen, "использование %s.%s вне main недопустимо", xIdent.Name, sel.Sel.Name)
						return true
					} else {
						if currentFunc != nil && currentFunc.Name.Name != "main" {
							pass.Reportf(call.Lparen, "использование %s.%s вне функции main недопустимо", xIdent.Name, sel.Sel.Name)
							return true
						}
					}
				}

			}
			return true
		})
	}
	return nil, nil
}
