package main

import (
	protobuf "awesomeProject6/transport-service/proto/transport"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

type repository interface {
	Available(*protobuf.Description) (*protobuf.Transport, error)
}

type TransportRepository struct {
	transports []*protobuf.Transport
}

func (repo *TransportRepository) Available(dscr *protobuf.Description) (*protobuf.Transport, error) {
	for _, transport := range repo.transports {
		if dscr.ContainerCount <= transport.ContainerCapacity && dscr.Weight <= transport.Weight {
			return transport, nil
		}
	}
	return nil, errors.New("no transport available")
}

type transportService struct {
	repo repository
}

func (s *transportService) Available(ctx context.Context, req *protobuf.Description)  (res *protobuf.Response, err error) {

	log.Printf("Incoming request: Get available transport")

	transport, err := s.repo.Available(req)
	if err != nil {
		log.Fatalf("Aval error: %v", err)
		return nil, err
	}
	log.Printf("Avalible transport: %v", transport)
	r := &protobuf.Response{Transport: transport, Transports: nil}
	return r, nil
}

func main() {
	log.SetOutput(os.Stdout)
	transports := []*protobuf.Transport{
		{Id: "transport1", Name: "Name1", Weight: 20000, ContainerCapacity: 5},
		{Id: "transport2", Name: "Name2", Weight: 1, ContainerCapacity: 1},
	}
	repo := &TransportRepository{transports}

	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}
	srv := grpc.NewServer()

	protobuf.RegisterTransportServiceServer(srv, &transportService{repo})

	reflection.Register(srv)

	log.Println("Server", ":9090")
	if err := srv.Serve(l); err != nil {
		log.Fatalf("Serve error: %v", err)
	}

}
