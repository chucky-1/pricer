// Package server implements the server side of grpc
package server

import (
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"
	log "github.com/sirupsen/logrus"
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

// Send listens on the channel and sends data to the client
func (s *Server) Send(id *protocol.Id, stream protocol.Pricer_SendServer) error {
	ch, err := s.rep.Send(int(id.Id))
	if err != nil {
		return err
	}
	for {
		st := <-ch
		stock := protocol.Stock{
			Id:    int32(st.ID),
			Title: st.Title,
			Price: st.Price,
		}
		err := stream.Send(&stock)
		if err != nil {
			log.Error(err)
		}
	}
}
