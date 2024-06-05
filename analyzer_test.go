package returnstyles

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	workdir, _ := os.Getwd()
	testdata := filepath.Join(workdir, "testdata")

	tests := []struct {
		testdir string
		options map[string]any
	}{
		{
			testdir: "default",
			options: nil,
		},
		{
			testdir: "default",
			options: map[string]any{
				AllowUnnamedFlag:        true,
				AllowNamedFlag:          true,
				AllowPartiallyNamedFlag: false,
				AllowNormalReturnsFlag:  true,
				AllowNakedReturnsFlag:   true,
				AllowMixingReturnsFlag:  false,
			},
		},
		{
			testdir: "strict",
			options: map[string]any{
				AllowUnnamedFlag:        true,
				AllowNamedFlag:          true,
				AllowPartiallyNamedFlag: false,
				AllowNormalReturnsFlag:  false,
				AllowNakedReturnsFlag:   true,
				AllowMixingReturnsFlag:  false,
			},
		},
		{
			testdir: "purism",
			options: map[string]any{
				AllowUnnamedFlag:        true,
				AllowNamedFlag:          true,
				AllowPartiallyNamedFlag: false,
				AllowNormalReturnsFlag:  true,
				AllowNakedReturnsFlag:   false,
				AllowMixingReturnsFlag:  false,
			},
		},
		{
			testdir: "disable_named",
			options: map[string]any{
				AllowUnnamedFlag:        true,
				AllowNamedFlag:          false,
				AllowPartiallyNamedFlag: false,
			},
		},
		{
			testdir: "disable_unnamed",
			options: map[string]any{
				AllowUnnamedFlag:        false,
				AllowNamedFlag:          true,
				AllowPartiallyNamedFlag: true,
			},
		},
		{
			testdir: "yaml",
			options: map[string]any{
				ConfigFlag: filepath.Join(testdata, "src/yaml/config.yaml"),
			},
		},
		{
			testdir: "generated_excluded",
			options: nil,
		},
		{
			testdir: "generated_included",
			options: map[string]any{
				IncludeGeneratedFlag: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testdir, func(t *testing.T) {
			Analyzer.Flags = flags()
			for k, v := range tt.options {
				Analyzer.Flags.Set(k, fmt.Sprintf("%v", v))
			}

			analysistest.Run(t, testdata, Analyzer, tt.testdir)
		})
	}
}
