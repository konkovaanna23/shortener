// Package main — это утилита статического анализа кода на языке Go.
//
// Она объединяет несколько сторонних и стандартных анализаторов
// в единый инструмент (multichecker), который можно запускать на проекте.
//
// Подключённые анализаторы:
//
//   - SA* (из staticcheck.io) — обнаружение ошибок, потенциальных багов,
//     проблем с производительностью и плохих паттернов.
//   - S10* (из staticcheck.io) — предложения по упрощению кода.
//   - ST*, Q*, LR* (из stylecheck) — проверка стиля и соглашений по именованию.
//   - Стандартные анализаторы из golang.org/x/tools/go/analysis/passes:
//     - asmdecl: проверка объявлений ассемблерных функций
//     - assign: обнаружение бесполезных присваиваний
//     - atomic: проверка правильного использования sync/atomic
//     - bools: логические ошибки в выражениях
//     - buildtag: синтаксис директив сборки
//     - errorsas: корректное использование errors.As
//     - httpresponse: проверка, что HTTP-ответ закрыт
//     - loopclosure: захват переменных в замыканиях циклов
//     - lostcancel: потеря контекстного отменяющего функционала
//     - nilness: обнаружение возможных nil-разыменований
//     - printf: соответствие формата и аргументов в fmt.Printf
//     - shadow: затенение переменных
//     - shift: сдвиг на небезопасное количество бит
//     - stdmethods: сигнатуры стандартных методов (например, String())
//     - structtag: синтаксис тегов структур
//     - tests: проверка тестов
//     - unmarshal: безопасность при unmarshaling
//     - unreachable: недостижимый код
//     - unsafeptr: использование unsafe.Pointer
//     - unusedresult: игнорирование возвращаемых значений
//   - exitanalyzer.ExitAnalyzer — анализатор, запрещающий
//     использование os.Exit вне функции main.main.
//
// Запуск:
//
//	go run cmd/staticlint/main.go ./...
//
// Или после сборки:
//
//	./staticlint ./...

package main

import (
	"fmt"

	"github.com/konkovaanna23/shortener/cmd/staticlint/exitanalyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"

	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	var checks []*analysis.Analyzer

	for _, v := range staticcheck.Analyzers {
		if v != nil && len(v.Analyzer.Name) >= 2 && v.Analyzer.Name[:2] == "SA" {
			checks = append(checks, v.Analyzer)
		}
	}

	for _, v := range staticcheck.Analyzers {
		if v != nil && len(v.Analyzer.Name) >= 3 && v.Analyzer.Name[:3] == "S10" {
			checks = append(checks, v.Analyzer)
		}
	}

	for _, v := range stylecheck.Analyzers {
		if v != nil {
			checks = append(checks, v.Analyzer)
		}
	}

	checks = append(checks,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilness.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	)

	checks = append(checks, exitanalyzer.ExitAnalyzer)

	fmt.Printf("Запуск %d анализаторов..\n", len(checks))
	multichecker.Main(checks...)
}
