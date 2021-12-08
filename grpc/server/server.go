// Package server implements the server side of grpc
package server

import (
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"

	"io"
)

// Server contains methods of application on service side
type Server struct {
	protocol.UnimplementedPricerServer
	rep *repository.Repository
}

// NewServer is constructor
func NewServer(rep *repository.Repository) *Server {
	return &Server{rep: rep}
}

// Get func takes an element from database
func (s *Server) Get(stream protocol.Pricer_GetServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		get, err := s.rep.Get(int(in.Id))
		if err != nil {
			return err
		}

		stock := protocol.Stock{
			Id:    int32(get.ID),
			Title: get.Title,
			Price: get.Price,
		}
		err = stream.Send(&stock)
		if err != nil {
			return err
		}
	}
}
