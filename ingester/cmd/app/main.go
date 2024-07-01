package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	pb "github.com/nnurry/gopds/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func BatchIngest(client pb.BatchIngestClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	stream, err := client.BatchIngest(ctx)

	defer func() {
		log.Println("Returned answer, cancelling context")
		cancel()
	}()

	if err != nil {
		log.Println("Encounter error while initializing client object", err)
		return false
	}

	for i := 0; i < 10; i++ {
		timeNow := time.Now().UTC()
		ingestRequest := &pb.IngestRequest{
			Meta: &pb.MetaField{
				UtcNow: timestamppb.New(timeNow),
				Key:    fmt.Sprint(i),
				Value:  fmt.Sprint(math.Round(float64(i) * math.Pi)),
			},
			Cardinal: &pb.CardinalField{
				Type: pb.CardinalType_STANDARD_HLL,
			},
			Filter: &pb.FilterField{
				Type:           pb.FilterType_STANDARD_BLOOM,
				MaxCardinality: 10000,
				ErrorRate:      0.0081,
			},
		}
		if err = stream.Send(ingestRequest); err != nil {
			log.Fatal("Can't send ingest request", err)
		} else {
			log.Println("Sent ingest request", ingestRequest.String(), ingestRequest.GetCardinal().GetType())
		}
	}
	response := &pb.BatchIngestResponse{}

	stream.CloseSend()

	err = stream.RecvMsg(&response)
	log.Println("got error while receiving the message: ", err)

	return response.Success
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(":50051", opts...)
	if err != nil {
		log.Fatal("Can't create gRPC client", err)
	}

	defer conn.Close()

	client := pb.NewBatchIngestClient(conn)

	log.Println("Batch ingest:", BatchIngest(client))

}
