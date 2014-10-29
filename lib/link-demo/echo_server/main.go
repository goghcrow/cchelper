package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"link"
)

var (
	benchmark  = flag.Bool("bench", false, "is for benchmark, will disable print")
	buffersize = flag.Int("buffer", 1024, "session read buffer size")
)

func log(v ...interface{}) {
	if !*benchmark {
		fmt.Println(v...)
	}
}

// This is an echo server demo work with the echo_client.
// usage:
//     go run echo_server/main.go
func main() {
	flag.Parse()

	link.DefaultReadBufferSize = *buffersize

	protocol := link.PacketN(2, binary.BigEndian)

	server, err := link.Listen("tcp", "127.0.0.1:10010", protocol)
	if err != nil {
		panic(err)
	}

	println("server start:", server.Listener().Addr().String())

	server.AcceptLoop(func(session *link.Session) {
		log("client", session.Conn().RemoteAddr().String(), "in")

		session.ReadLoop(func(msg []byte) {
			log("client", session.Conn().RemoteAddr().String(), "say:", string(msg))
			session.Send(link.Binary(msg))
		})

		log("client", session.Conn().RemoteAddr().String(), "close")
	})
}
