package gossa

import (
	"fmt"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var Analyzer = &analysis.Analyzer{
	Name: "gossaanalysis",
	Doc:  "go-ssa-analysis",
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs

	for _, fn := range funcs {
		// デバッグ出力
		fn.WriteTo(os.Stdout)
		for _, b := range fn.Blocks {
			// fmt.Printf("funcs[i].Blocks -> %#v\n", fn.Blocks)
			// すべての命令を調べる
			for _, instr := range b.Instrs {
				call, ok := instr.(ssa.CallInstruction)
				if !ok {
					fmt.Printf("not call instruction %T %+v\n", instr, instr)
					continue
				}

				// TODO この辺見直し
				callee := call.Common().StaticCallee()
				if callee == nil {
					// is not static
					continue
				}

				if callee.Pkg.Pkg.Path() != "a" { // TODO customizable
					continue
				}

				fmt.Println("params ->", callee.Name(), callee.Params)
				args := call.Operands(nil)
				if len(args) == 0 {
					continue
				}

				for i, arg := range args[1:] { // skip receiver
					if arg == nil {
						continue
					}

					v := *arg
					fmt.Printf("param[%d] -> %s %#v\n", i, callee.Name(), v)
					if !isDeterministic(v) {
						pass.Reportf(call.Pos(), "The message of %s should be constant", callee.Name())
					}
				}
			}
		}
	}
	return nil, nil
}

func isDeterministic(v ssa.Value) bool {
	switch v := v.(type) {
	case *ssa.Phi:
		for _, e := range v.Edges {
			if !isDeterministic(e) {
				return false
			}
		}
		return true
	case *ssa.Const:
		// ok
		return true
	case *ssa.Slice:
		refs := v.X.Referrers()
		if refs == nil {
			return true
		}
		var addr *ssa.IndexAddr
		for _, e := range *refs {
			a, ok := e.(*ssa.IndexAddr)
			if !ok {
				continue
			}
			addr = a
		}
		fmt.Println("addr pos ->", addr.Pos())
		addrRefs := addr.Referrers()
		if addrRefs == nil {
			return true
		}
		for _, e := range *addrRefs {
			s, ok := e.(*ssa.Store)
			if !ok {
				continue
			}
			if !isDeterministic(s.Val) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
