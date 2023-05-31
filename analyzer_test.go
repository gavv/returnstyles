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
		name    string
		options map[string]any
	}{
		{
			name:    "default",
			options: nil,
		},
		{
			name: "default",
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
			name: "strict",
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
			name: "purism",
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
			name: "disable_named",
			options: map[string]any{
				AllowUnnamedFlag:        true,
				AllowNamedFlag:          false,
				AllowPartiallyNamedFlag: false,
			},
		},
		{
			name: "disable_unnamed",
			options: map[string]any{
				AllowUnnamedFlag:        false,
				AllowNamedFlag:          true,
				AllowPartiallyNamedFlag: true,
			},
		},
		{
			name: "yaml",
			options: map[string]any{
				ConfigFlag: filepath.Join(testdata, "src/yaml/config.yaml"),
			},
		},
		{
			name:    "generated_excluded",
			options: nil,
		},
		{
			name: "generated_included",
			options: map[string]any{
				AllowUnnamedFlag:        true,
				AllowNamedFlag:          false,
				AllowPartiallyNamedFlag: false,
				IncludeGeneratedFlag:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Analyzer.Flags = flags()
			for k, v := range tt.options {
				Analyzer.Flags.Set(k, fmt.Sprintf("%v", v))
			}

			analysistest.Run(t, testdata, Analyzer, tt.name)
		})
	}
}
