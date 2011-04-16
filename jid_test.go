package main

import "fmt"
import "testing"

func Test(t *testing.T) {
}

func Benchmark(bm *testing.B) {
	var j Jid
	for i := 0; i < bm.N; i++ {
		j = makeJid("local@domain/resource")
	}
	fmt.Println(j)
}
