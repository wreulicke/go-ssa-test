package main

import (
	gossa "github.com/wreulicke/go-ssa-test"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(gossa.Analyzer)
}
