package main

import (
	"bytes"
	"fmt"
	"runtime"
)

const (
	Normal = "\033[0m"
	Bold = "\033[1m"
	Red = "\033[31m"
	Blue = "\033[34m"
	White = "\033[37m"
	OnBlack = "\033[40m"
)

func debugln(v ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)

	b := new(bytes.Buffer)

	name := f.Name()

	fmt.Fprint(b, White, name, Normal, " ")

	switch name {
	case "main.*Stream·Read":
		fmt.Fprint(b, Bold, Red, OnBlack)
		fmt.Fprintln(b, v...)
		fmt.Fprint(b, Normal)
	case "main.*Stream·Write":
		fmt.Fprint(b, Bold, Blue, OnBlack)
		fmt.Fprintln(b, v...)
		fmt.Fprint(b, Normal)
	default:
		fmt.Fprintln(b, v...)
	}

	fmt.Print(b.String())
}
