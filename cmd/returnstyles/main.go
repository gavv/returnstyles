package main

import (
	"github.com/gavv/returnstyles"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(returnstyles.Analyzer)
}
