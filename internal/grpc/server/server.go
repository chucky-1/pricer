// Package server implements the server side of grpc
package server

import (
	"github.com/chucky-1/pricer/internal/model"
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"errors"
)

// Server contains methods of application on service side
type Server struct {
	protocol.UnimplementedPricesServer
	rep *repository.Repository
}

// NewServer is constructor
func NewServer(rep *repository.Repository) *Server {
	return &Server{rep: rep}
}

// Sub listens on the channel and sends data to the client
func (s *Server) Sub(stream protocol.Prices_SubServer) error {
	grpcID := uuid.New().String()
	internalChan := make(chan *protocol.StockID)
	externalChan := make(chan *model.Stock)

	recv, err := stream.Recv()
	if err != nil {
		return err
	}
	if recv.Act == "INIT" {
		s.rep.Sub(recv.List, grpcID, externalChan)
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
					Update: st.Update,
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
			err = s.rep.Close(grpcID)
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

func (s *Server) SubAll(r *protocol.Request, stream protocol.Prices_SubAllServer) error {
	grpcID := uuid.New().String()
	ch := make(chan *model.Stock)
	s.rep.SubAll(grpcID, ch)
	for {
		select {
		case <-stream.Context().Done():
			err := s.rep.Close(grpcID)
			if err != nil {
				log.Error(err)
			}
			return stream.Context().Err()
		case st := <-ch:
			stock := protocol.Stock{
				Id:     st.ID,
				Title:  st.Title,
				Price:  st.Price,
				Update: st.Update,
			}
			err := stream.Send(&stock)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
