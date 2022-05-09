// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func handleConn(c net.Conn, loc *time.Location) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, fmt.Sprintf("%s\t : %s", loc, time.Now().In(loc).Format("15:04:05\n")))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	p := flag.Int("port", 0, "Port in whic the server will run")
	flag.Parse()
	if *p <= 0 {
		log.Fatal("Invalid port")
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *p))
	if err != nil {
		log.Fatal(err)
	}
	loc, err := time.LoadLocation(os.Getenv("TZ"))
	if err != nil {
		println("Invalid time zone '%s' using UTC instead.")
		loc = time.UTC
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn, loc) // handle connections concurrently
	}
}
