// Package server implements the server side of grpc
package server

import (
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"context"
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
func (s *Server) Send(listID *protocol.ListID, stream protocol.Pricer_SendServer) error {
	grpcID := uuid.New().String()
	ch, err := s.rep.Send(listID.Id, grpcID)
	if err != nil {
		log.Error(err)
	}
	for {
		select {
		case <-stream.Context().Done():
			go func() {
				grpc := protocol.GrpcID{Id: grpcID}
				_, err = s.Close(context.Background(), &grpc)
				if err != nil {
					log.Error(err)
				}
			}()
			return stream.Context().Err()
		case st := <-ch:
			stock := protocol.Stock{
				Id:    st.ID,
				Title: st.Title,
				Price: st.Price,
			}
			err = stream.Send(&stock)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

// Close func called when the client exits
func (s *Server) Close(ctx context.Context, grpcID *protocol.GrpcID) (*protocol.Response, error) {
	err := s.rep.Close(grpcID.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
