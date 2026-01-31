package exitanalyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestCheckAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ExitAnalyzer, "a", "b")
}
