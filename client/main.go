package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"

	"github.com/google/uuid"
)

func main() {
	for i := 0; i < 1_000_000; i++ {

		id := uuid.New() // 16 bytes
		lat := 37.7749   // float64 (8 bytes)
		lng := -122.4194

		buff := new(bytes.Buffer)
		buff.Write(id[:])
		binary.Write(buff, binary.BigEndian, lat)
		binary.Write(buff, binary.BigEndian, lng)

		data := buff.Bytes()

		serverAddr := "0.0.0.0:8080"
		conn, err := net.Dial("udp", serverAddr)
		if err != nil {
			log.Fatalf("Failed to connect to server: %s", err.Error())
		}
		defer conn.Close()

		_, err = conn.Write(data)
		if err != nil {
			log.Fatalf("Failed to write to server: %s", err.Error())
		}
	}

}
