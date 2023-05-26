package returnstyles

import (
	"flag"
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"gopkg.in/yaml.v3"
)

var Analyzer = &analysis.Analyzer{
	Name:     "returnstyles",
	Doc:      "Return styles linter",
	URL:      "https://github.com/gavv/returnstyles",
	Flags:    flags(),
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

const (
	ConfigFlag              = "config"
	AllowUnnamedFlag        = "allow-unnamed"
	AllowNamedFlag          = "allow-named"
	AllowPartiallyNamedFlag = "allow-partially-named"
	AllowNormalReturnsFlag  = "allow-normal-returns"
	AllowNakedReturnsFlag   = "allow-naked-returns"
	AllowMixingReturnsFlag  = "allow-mixing-returns"
	IncludeCgoFlag          = "include-cgo"
)

type styleConfig struct {
	AllowUnnamed        bool `yaml:"allow-unnamed"`
	AllowNamed          bool `yaml:"allow-named"`
	AllowPartiallyNamed bool `yaml:"allow-partially-named"`
	AllowNormalReturns  bool `yaml:"allow-normal-returns"`
	AllowNakedReturns   bool `yaml:"allow-naked-returns"`
	AllowMixingReturns  bool `yaml:"allow-mixing-returns"`
	IncludeCgo          bool `yaml:"include-cgo"`
}

var config = styleConfig{
	AllowUnnamed:        true,
	AllowNamed:          true,
	AllowPartiallyNamed: false,
	AllowNormalReturns:  true,
	AllowNakedReturns:   true,
	AllowMixingReturns:  false,
	IncludeCgo:          false,
}

var configPath string

type returnStyle int

const (
	normalReturn returnStyle = iota
	nakedReturn
)

type functionStyle int

const (
	noReturnVariables functionStyle = iota
	unnamedReturnVariables
	namedReturnVariables
	partiallyNamedReturnVariables
)

func flags() flag.FlagSet {
	var fs flag.FlagSet

	stringVar := func(ptr *string, flag string, help string) {
		fs.StringVar(ptr, flag, *ptr, help)
	}

	boolVar := func(ptr *bool, flag string, help string) {
		fs.BoolVar(ptr, flag, *ptr, help)
	}

	stringVar(&configPath, ConfigFlag,
		"yaml config file path")

	boolVar(&config.AllowUnnamed, AllowUnnamedFlag,
		"allow functions with unnamed return variables")

	boolVar(&config.AllowNamed, AllowNamedFlag,
		"allow functions with named return variables")

	boolVar(&config.AllowPartiallyNamed, AllowPartiallyNamedFlag,
		"allow functions with partially named return variables")

	boolVar(&config.AllowNormalReturns, AllowNormalReturnsFlag,
		"allow normal (non-naked) returns in functions with named return variables")

	boolVar(&config.AllowNakedReturns, AllowNakedReturnsFlag,
		"allow naked returns in functions with named return variables")

	boolVar(&config.AllowMixingReturns, AllowMixingReturnsFlag,
		"allow mixing normal and naked in functions with named return variables")

	boolVar(&config.IncludeCgo, IncludeCgoFlag,
		"include linting cgo functions")

	return fs
}

func run(pass *analysis.Pass) (interface{}, error) {
	if configPath != "" {
		readConfig(&config, configPath)
	}

	var (
		visitedNodes        = make(map[ast.Node]bool)
		visitedReturnStyles = make(map[ast.Node]map[returnStyle]bool)
	)

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	inspector.WithStack(
		[]ast.Node{
			(*ast.FuncDecl)(nil),
			(*ast.FuncLit)(nil),
			(*ast.ReturnStmt)(nil),
		},
		func(node ast.Node, push bool, stack []ast.Node) (proceed bool) {
			if _, ok := visitedNodes[node]; ok {
				return false
			}

			switch currNode := node.(type) {
			case *ast.FuncDecl, *ast.FuncLit:
				if skip(currNode) {
					return false
				}

				validateFunc(pass, currNode)

			case *ast.ReturnStmt:
				funcNode := findFuncNode(stack)
				retnStyle := detectReturnStyle(funcNode, currNode)

				if skip(funcNode) {
					return false
				}

				if visitedReturnStyles[funcNode] == nil {
					visitedReturnStyles[funcNode] = make(map[returnStyle]bool)
				}
				visitedReturnStyles[funcNode][retnStyle] = true

				validateReturn(
					pass, funcNode, currNode, visitedReturnStyles[funcNode],
				)
			}

			visitedNodes[node] = true
			return true
		})

	return nil, nil
}

func readConfig(result *styleConfig, path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open -config file %q: %v", path, err)
		os.Exit(1)
	}

	var config = struct {
		Styles styleConfig `yaml:"returnstyles"`
	}{
		Styles: *result,
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't parse -config file %q: %v", path, err)
		os.Exit(1)
	}

	*result = config.Styles
}

func skip(funcNode ast.Node) bool {
	switch funcNode := funcNode.(type) {
	case *ast.FuncDecl:
		if strings.HasPrefix(funcNode.Name.Name, "_Cfunc_") && !config.IncludeCgo {
			return true
		}
	}

	return false
}

func validateFunc(pass *analysis.Pass, funcNode ast.Node) {
	funcStyle := detectFuncStyle(funcNode)

	switch funcStyle {
	case noReturnVariables:
		// no-op

	case unnamedReturnVariables:
		if !config.AllowUnnamed {
			pass.Reportf(funcNode.Pos(),
				"functions with unnamed return variables not allowed")
		}

	case namedReturnVariables:
		if !config.AllowNamed {
			pass.Reportf(funcNode.Pos(),
				"functions with named return variables not allowed")
		}

	case partiallyNamedReturnVariables:
		if !config.AllowNamed {
			pass.Reportf(funcNode.Pos(),
				"functions with named return variables not allowed")
		} else if !config.AllowPartiallyNamed {
			pass.Reportf(funcNode.Pos(),
				"functions with partially named return variables not allowed")
		}
	}
}

func validateReturn(
	pass *analysis.Pass,
	funcNode, returnNode ast.Node,
	otherReturnStyles map[returnStyle]bool,
) {
	funcStyle := detectFuncStyle(funcNode)
	returnStyle := detectReturnStyle(funcNode, returnNode)

	switch funcStyle {
	case namedReturnVariables, partiallyNamedReturnVariables:
		if !config.AllowNormalReturns && returnStyle == normalReturn {
			pass.Reportf(returnNode.Pos(),
				"normal (non-naked) returns not allowed in functions with named return variables")
		}

		if !config.AllowNakedReturns && returnStyle == nakedReturn {
			pass.Reportf(returnNode.Pos(),
				"naked returns not allowed")
		}

		if !config.AllowMixingReturns &&
			config.AllowNormalReturns &&
			config.AllowNakedReturns &&
			len(otherReturnStyles) > 1 {
			pass.Reportf(returnNode.Pos(),
				"mixing normal and naked returns not allowed")
		}
	}
}

func detectReturnStyle(funcNode, returnNode ast.Node) returnStyle {
	funcType := findFuncType(funcNode)

	returnStmt, ok := returnNode.(*ast.ReturnStmt)
	if !ok {
		panic("invalid return node")
	}

	if funcType.Results != nil && len(funcType.Results.List) != 0 {
		if len(returnStmt.Results) == 0 {
			return nakedReturn
		}
	}

	return normalReturn
}

func detectFuncStyle(funcNode ast.Node) functionStyle {
	funcType := findFuncType(funcNode)

	hasNamed := false
	hasUnnamed := false

	if funcType.Results != nil {
		numNamed := 0

		for _, field := range funcType.Results.List {
			if len(field.Names) != 0 {
				hasNamed = true
				numNamed++
			}
		}

		if numNamed < funcType.Results.NumFields() {
			hasUnnamed = true
		}
	}

	if hasNamed && hasUnnamed {
		return partiallyNamedReturnVariables
	}

	if hasNamed {
		return namedReturnVariables
	}

	if !hasNamed && !hasUnnamed {
		return noReturnVariables
	}

	return unnamedReturnVariables
}

func findFuncType(funcNode ast.Node) *ast.FuncType {
	switch funcNode := funcNode.(type) {
	case *ast.FuncDecl:
		return funcNode.Type
	case *ast.FuncLit:
		return funcNode.Type
	}
	panic("invalid function node")
}

func findFuncNode(stack []ast.Node) ast.Node {
	for n := len(stack) - 1; n != 0; n-- {
		switch node := stack[n].(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			return node
		}
	}
	panic("invalid return statement node")
}
