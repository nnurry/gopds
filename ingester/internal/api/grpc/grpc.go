package grpc

import (
	"context"
	"log"

	pb "github.com/nnurry/gopds/protos"
)

var IngestChannel = make(chan *pb.IngestRequest, 10000)

func BatchIngest(client pb.BatchIngestClient) bool {
	log.Println("Running batch request")
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := client.BatchIngest(ctx)

	defer func() {
		log.Println("Returned answer, cancelling context")
		cancel()
	}()

	if err != nil {
		log.Println("Encounter error while initializing client object", err)
		return false
	}

	for i := 1; ; i++ {
		ingestRequest := <-IngestChannel
		if err = stream.Send(ingestRequest); err != nil {
			log.Println("Can't send ingest request, cancelling context", err)
			return true
		} else {
			log.Println("Sent ingest request", i, ":", ingestRequest.String(), ingestRequest.GetCardinal().GetType())
		}
	}
}
