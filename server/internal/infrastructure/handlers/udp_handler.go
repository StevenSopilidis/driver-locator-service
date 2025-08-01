package handlers

import (
	"encoding/binary"
	"log"
	"math"
	"net"

	"github.com/google/uuid"
)

type UDPServer struct {
	conn               *net.UDPConn
	addr               string
	port               int
	listening          bool
	max_concurent_reqs int
}

func NewUDPServer(addr string, port int, max_concurent_reqs int) (*UDPServer, error) {
	udpAddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(addr),
	}

	conn, err := net.ListenUDP("udp", &udpAddr)
	if err != nil {
		return nil, err
	}

	return &UDPServer{
		conn:               conn,
		listening:          false,
		addr:               addr,
		port:               port,
		max_concurent_reqs: max_concurent_reqs,
	}, nil
}

func (s *UDPServer) ListenAndServe() {
	s.listening = true
	log.Printf("Server starting at: %s:%d", s.addr, s.port)

	sem := make(chan struct{}, s.max_concurent_reqs)

	for s.listening {
		// 16 bytes for uuid
		// 16 bytes for long + lat

		sem <- struct{}{} // if channel full wait to release some
		buffer := make([]byte, 32)
		n, clientAddr, err := s.conn.ReadFromUDP(buffer)

		if err != nil {
			log.Printf("Error reading from client: %s", clientAddr.IP.String())
			continue
		}

		if n != 32 {
			log.Printf("Received %d bytes expected %d", n, 32)
			continue
		}

		go func(data []byte, addr *net.UDPAddr) {
			defer func() { <-sem }() // release semaphore
			s.handleRequest(buffer, clientAddr)
		}(buffer, clientAddr)
	}
}

func (s *UDPServer) handleRequest(data []byte, addr *net.UDPAddr) {
	log.Println("Received address location for driver: ", addr.IP.String())

	id, err := uuid.FromBytes(data[:16])
	if err != nil {
		log.Fatalf("Failed to parse uuid: %s", err.Error())
	}

	lat := math.Float64frombits(binary.BigEndian.Uint64(data[16:24]))
	lng := math.Float64frombits(binary.BigEndian.Uint64(data[24:32]))

	log.Printf("Received data from %s with long: %f and lat: %f", id.String(), lat, lng)
}

func (s *UDPServer) Shutdown() {
	err := s.conn.Close()
	if err != nil {
		log.Fatalf("Could not close udp listener: %s", err.Error())
	}
}
