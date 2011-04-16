package main

import (
	"fmt"
	"io"
	"os"
)

type LogReadWriter struct { io.ReadWriter }

func (rw *LogReadWriter) Read(b []byte) (n int, e os.Error) {
	n, e = rw.ReadWriter.Read(b)
	fmt.Println("read", string(b))
	return
}

func (rw *LogReadWriter) Write(b []byte) (n int, e os.Error) {
	fmt.Println("write", string(b))
	n, e = rw.ReadWriter.Write(b)
	return
}
