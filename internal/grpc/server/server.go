// Package server implements the server side of grpc
package server

import (
	"github.com/chucky-1/pricer/internal/model"
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"context"
	"errors"
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
func (s *Server) Send(stream protocol.Pricer_SendServer) error {
	grpcID := uuid.New().String()
	internalChan := make(chan *protocol.StockID)
	externalChan := make(chan *model.Stock)

	recv, err := stream.Recv()
	if err != nil {
		return err
	}
	if recv.Act == "INIT" {
		s.rep.Send(recv.List, grpcID, externalChan)
		go func() {
			for {
				st, ok := <-externalChan
				if !ok {
					break
				}
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
		}()
	} else {
		return errors.New("init not made")
	}

	go func() {
		for {
			select {
			case <-stream.Context().Done():
				return
			default:
				recv, err = stream.Recv()
				if err != nil {
					log.Error(err)
					continue
				}
				internalChan <- recv
			}
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			grpc := protocol.GrpcID{Id: grpcID}
			_, err = s.Close(context.Background(), &grpc)
			if err != nil {
				log.Error(err)
			}
			return stream.Context().Err()
		case recv = <-internalChan:
			switch {
			case recv.Act == "ADD":
				s.rep.Add(recv.List, grpcID, externalChan)
			case recv.Act == "DEL":
				s.rep.Del(recv.List, grpcID)
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
