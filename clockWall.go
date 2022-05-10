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
	defer conn.Close()

	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)

		if n == 0 || err != nil {
			return
		}

		printer <- string(buf)
	}
}

func printer(times chan string) {
	for time := range times {
		fmt.Print(time)
	}

}

func main() {
	times := make(chan string, len(os.Args)-1)
	defer close(times)

	done := make(chan struct{})
	defer close(done)

	goroutines := 0
	for _, arg := range os.Args[1:] {
		v := strings.Split(arg, "=")
		if len(v) != 2 {
			fmt.Printf("Invalid argument `%v`, skipping\n", arg)
		}
		goroutines++
		go func() {
			dialClock(v[1], times)
			done <- struct{}{}
		}()
	}

	go printer(times)

	for goroutines > 0 {
		<-done
		goroutines--
	}

}
