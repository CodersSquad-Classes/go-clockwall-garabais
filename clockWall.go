package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func dialClock(address string, printer chan string) {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Unable to connect to %v, skipping\n", address)
		return
	}

	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)

		if n == 0 || err != nil {
			return
		}

		printer <- string(buf)
	}
}

func main() {
	times := make(chan string, len(os.Args)-1)
	defer close(times)

	for _, arg := range os.Args[1:] {
		v := strings.Split(arg, "=")
		if len(v) != 2 {
			fmt.Printf("Invalid argument `%v`, skipping\n", arg)
		}
		go dialClock(v[1], times)
	}

	for time := range times {
		fmt.Print(time)
	}
}
