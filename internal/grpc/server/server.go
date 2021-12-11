// Package server implements the server side of grpc
package server

import (
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"errors"
)

// Server contains methods of application on service side
type Server struct {
	protocol.UnimplementedPricerServer
	rep    *repository.Repository
}

// NewServer is constructor
func NewServer(rep *repository.Repository) *Server {
	return &Server{rep: rep}
}

// Send listens on the channel and sends data to the client
func (s *Server) Send(id *protocol.Id, stream protocol.Pricer_SendServer) error {
	userID := uuid.New().String()
	if userID == "" {
		return errors.New("UserID didn't generate")
	}
	ch, err := s.rep.Send(int(id.Id), userID)
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
		err = stream.Send(&stock)
		if err != nil {
			switch {
			case err.Error() == "rpc error: code = Unavailable desc = transport is closing":
				err = s.rep.Close(int(id.Id), userID)
				if err != nil {
					log.Error(err)
				} else {
					return errors.New("client exit")
				}
			default:
				log.Error(err)
			}
		}
	}
}
