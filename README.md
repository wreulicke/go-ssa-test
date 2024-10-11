# go-ssa-test

This is my repository template for cli tool in go.

## Installation

```
go install github.com/wreulicke/go-ssa-test/cmd/go-ssa-test@latest
```

## Usage

```bash
go-ssa-test -h
gossaanalysis: go-ssa-analysis

Usage: gossaanalysis [-flag] [package]


Flags:
  -V    print version and exit
  -all
        no effect (deprecated)
  -c int
        display offending line with this many lines of context (default -1)
  -cpuprofile string
        write CPU profile to this file
  -debug string
        debug flags, any subset of "fpstv"
  -fix
        apply all suggested fixes
  -flags
        print analyzer flags in JSON
  -json
        emit JSON output
  -memprofile string
        write memory profile to this file
  -source
        no effect (deprecated)
  -tags string
        no effect (deprecated)
  -test
        indicates whether test files should be analyzed, too (default true)
  -trace string
        write trace log to this file
  -v    no effect (deprecated)
```

## License

MIT License