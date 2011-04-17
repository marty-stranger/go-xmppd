package s2s

import "net"
import "fmt"
import "testing"

func _Test(t *testing.T) {
	addr := address("google.com")
	fmt.Println(addr)

	c, e := net.Dial("tcp", "", addr)
	if e != nil { panic(e) }
	fmt.Println(c)
}

func Test(t *testing.T) {
	go Run()
	Connect("gmail.com")
}
