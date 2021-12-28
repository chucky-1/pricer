// Package server implements the server side of grpc
package server

import (
	"github.com/chucky-1/pricer/internal/model"
	"github.com/chucky-1/pricer/internal/repository"
	"github.com/chucky-1/pricer/protocol"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
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

// Subscribe listens on the channel and sends data to the client
func (s *Server) Subscribe(stream protocol.Prices_SubscribeServer) error {
	grpcID := uuid.New().String()
	ch := make(chan *model.Symbol)

	go func() {
		for {
			select {
			case <- stream.Context().Done():
				return
			default:
				symbol, ok := <-ch
				if !ok {
					continue
				}
				time, err := decodeTime(symbol.Time)
				if err != nil {
					log.Error(err)
				}
				err = stream.Send(&protocol.SubscribeResponse{
					SymbolId: symbol.ID,
					Bid:      symbol.Bid,
					Ask:      symbol.Ask,
					Update:   &timestamppb.Timestamp{Seconds: time},
				})
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			err := s.rep.Close(grpcID)
			if err != nil {
				log.Error(err)
			}
			return stream.Context().Err()
		default:
			recv, err := stream.Recv()
			if err != nil {
				log.Error(err)
				continue
			}
			switch {
			case recv.Action.String() == "ADD":
				symbolIDList := make([]int32, 0, len(recv.SymbolId))
				for _, symbolID := range recv.SymbolId {
					symbolIDList = append(symbolIDList, symbolID)
				}
				s.rep.Add(symbolIDList, grpcID, ch)
			case recv.Action.String() == "DEL":
				symbolIDList := make([]int32, 0, len(recv.SymbolId))
				for _, symbolID := range symbolIDList {
					symbolIDList = append(symbolIDList, symbolID)
				}
				s.rep.Del(symbolIDList, grpcID)
			}
		}
	}
}

func decodeTime(time string) (int64, error) {
	mkr, err := strconv.Atoi(time)
	if err != nil {
		return 0, err
	}
	return int64(mkr)/1000, nil
}
