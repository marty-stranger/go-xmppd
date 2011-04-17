package main

import "xmppd/s2s"
import "testing"
import __regexp__ "regexp"

var tests = []testing.InternalTest{
	{"s2s.Test", s2s.Test},
}
var benchmarks = []testing.InternalBenchmark{ //
}

func main() {
	testing.Main(__regexp__.MatchString, tests)
	testing.RunBenchmarks(__regexp__.MatchString, benchmarks)
}
