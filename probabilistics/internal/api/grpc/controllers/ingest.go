package controllers

import (
	"io"
	"log"
	"math"

	request_schema "github.com/nnurry/gopds/probabilistics/internal/api/rest/schemas/request"
	"github.com/nnurry/gopds/probabilistics/internal/database/postgres"
	"github.com/nnurry/gopds/probabilistics/internal/service"
	"github.com/nnurry/gopds/probabilistics/pkg/models/wrapper"
	pb "github.com/nnurry/gopds/protos"
)

type batchIngestServer struct {
	pb.UnimplementedBatchIngestServer
}

func NewBatchIngestServer() *batchIngestServer {
	return &batchIngestServer{}
}

func (srv *batchIngestServer) BatchIngest(stream pb.BatchIngest_BatchIngestServer) error {

	filterWrapper := wrapper.GetWrapper().FilterWrapper()
	cardinalWrapper := wrapper.GetWrapper().CardinalWrapper()

	tx, _ := postgres.Client.Begin()

	for i := 1; ; i++ {
		// READ MESSAGE FROM THE STREAM
		message, err := stream.Recv()
		if err == io.EOF {
			log.Println("At the end of the stream, committing changes")
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			response := &pb.BatchIngestResponse{
				Success: true,
			}

			err = stream.SendAndClose(response)
			return err
		}
		if err != nil {
			log.Println("Can't process batch ingest request", i, " due to", err)
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			return err
		}

		// PREPARING METADATA
		m := message.GetMeta()
		f := message.GetFilter()
		c := message.GetCardinal()

		filterKey := wrapper.FilterKey{
			Type:           f.GetType().String(),
			Key:            m.GetKey(),
			MaxCardinality: uint(f.GetMaxCardinality()),
			ErrorRate:      math.Trunc(float64(f.ErrorRate)*10000000.0) / 10000000,
		}

		cardinalKey := wrapper.CardinalKey{
			Type: c.GetType().String(),
			Key:  m.GetKey(),
		}

		// FETCHING FILTER/CARDINAL OR CREATING AND SAVE ONE IF NOT EXISTS
		decayFilter := filterWrapper.GetFilter(filterKey, false)

		if decayFilter == nil {
			decayFilter = service.CreateFilter(&request_schema.FilterCreateBody{
				Meta: request_schema.MetaBody{
					Key: filterKey.Key,
				},
				Filter: request_schema.FilterBody{
					Type:           filterKey.Type,
					MaxCardinality: filterKey.MaxCardinality,
					ErrorRate:      filterKey.ErrorRate,
				},
			})

			err := service.SaveFilter(decayFilter, true, true, true, tx)

			tx, _ = postgres.Client.Begin()

			if err != nil {
				log.Println("Got error while saving the filter into database:", err)
			}
		}

		decayCardinal := cardinalWrapper.GetCardinal(cardinalKey, false)
		if decayCardinal == nil {
			decayCardinal = service.CreateCardinal(&request_schema.CardinalCreateBody{
				Meta: request_schema.MetaBody{
					Key: cardinalKey.Key,
				},
				Cardinal: request_schema.CardinalBody{
					Type: cardinalKey.Type,
				},
			})

			err := service.SaveCardinal(decayCardinal, true, true, true, tx)

			tx, _ = postgres.Client.Begin()

			if err != nil {
				log.Println("Got error while saving the cardinal into database:", err)
			}
		}

		// ADD FILTER/CARDINAL INTO DATABASE
		decayFilter.AddString(m.GetValue())
		decayCardinal.AddString(m.GetValue())

		service.SaveFilter(decayFilter, false, false, true, tx)
		service.SaveCardinal(decayCardinal, false, false, true, tx)
	}

}
