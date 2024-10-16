package main

import (
	"math/rand"
	"strconv"
)

func main() {
	x := "literal"
	y := "literal"

	if v := rand.Int(); v == 1 {
		x = strconv.Itoa(v)
		y = "still literal"
	}

	log(x)         // want `The message of log should be constant`
	log(y)         // it's ok
	log("literal") // it's ok
	log(f())       // want `The message of log should be constant`

	logVargs(x)         // want `The message of logVargs should be constant`
	logVargs(y)         // it's ok
	logVargs("literal") // it's ok
	logVargs(f())       // want `The message of logVargs should be constant`

	var slice []string
	slice = append(slice, "literal")
	logVargs(slice...)
}

func log(x string) {
	// do nothing
}

func logVargs(x ...string) {
	// do nothing
}

func f() string {
	return "literal"
}
