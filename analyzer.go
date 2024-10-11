package gossa

import (
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

				args := call.Operands(nil)
				if len(args) == 0 {
					continue
				}

				for _, arg := range args[1:] { // skip receiver
					if arg == nil {
						continue
					}

					v := *arg
					switch v := v.(type) {
					case *ssa.Phi:
						isConst := true
						for _, e := range v.Edges {
							if _, ok := e.(*ssa.Const); !ok {
								isConst = false
								break
							}
						}
						if !isConst {
							pass.Reportf(call.Pos(), "The message of %s should be constant", callee.Name())
						}
					case *ssa.Const:
						// ok
					default:
						pass.Reportf(call.Pos(), "The message of %s should be constant", callee.Name())
					}
				}
			}
		}
	}

	return nil, nil
}
