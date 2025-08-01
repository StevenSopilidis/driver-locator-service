package handlers

import (
	"context"
	"encoding/binary"
	"log"
	"math"
	"net"
	"time"

	"github.com/StevenSopilidis/driver-locator-service/internal/domain"
	workerpool "github.com/StevenSopilidis/driver-locator-service/internal/infrastructure/worker_pool"
	"github.com/google/uuid"
)

type UDPServer struct {
	conn               *net.UDPConn
	addr               string
	port               int
	max_concurent_reqs int
	dataCh             chan *domain.Driver
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

	dataCh := make(chan *domain.Driver, max_concurent_reqs)

	return &UDPServer{
		conn:               conn,
		addr:               addr,
		port:               port,
		max_concurent_reqs: max_concurent_reqs,
		dataCh:             dataCh,
	}, nil
}

func (s *UDPServer) ListenAndServe(ctx context.Context) {
	log.Printf("Server starting at: %s:%d", s.addr, s.port)

	sem := make(chan struct{}, s.max_concurent_reqs)

	pool := workerpool.NewWorkerPool()
	go pool.Run(s.dataCh)

	for {
		select {
		case <-ctx.Done():
			log.Println("Server shutting down")
			return
		default:
			// 16 bytes for uuid
			// 16 bytes for long + lat
			sem <- struct{}{} // if channel full wait to release some
			buffer := make([]byte, 32)
			n, clientAddr, err := s.conn.ReadFromUDP(buffer)

			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					<-sem
					continue
				}

				log.Printf("Read error: %v", err)
				<-sem
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
}

func (s *UDPServer) handleRequest(data []byte, addr *net.UDPAddr) {
	log.Println("Received address location for driver: ", addr.IP.String())

	id, err := uuid.FromBytes(data[:16])
	if err != nil {
		return
	}

	lat := math.Float64frombits(binary.BigEndian.Uint64(data[16:24]))
	lng := math.Float64frombits(binary.BigEndian.Uint64(data[24:32]))

	s.dataCh <- &domain.Driver{
		Id:        id,
		Latitude:  lat,
		Longitude: lng,
		Last_seen: time.Now(),
	}
}

func (s *UDPServer) Shutdown() {
	close(s.dataCh)

	if err := s.conn.Close(); err != nil {
		log.Fatalf("Could not close udp listener: %s", err.Error())
	}
}
