package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nnurry/gopds/probabilistics/internal/api/grpc/controllers"
	restApi "github.com/nnurry/gopds/probabilistics/internal/api/rest"
	"github.com/nnurry/gopds/probabilistics/internal/database/postgres"
	"github.com/nnurry/gopds/probabilistics/pkg/models/wrapper"
	pb "github.com/nnurry/gopds/protos"
	"google.golang.org/grpc"
)

func main() {
	postgres.Bootstrap()
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGTERM, syscall.SIGINT)

	wrapper.DecayWg.Add(1)

	mux := restApi.SetupMux()
	httpSrv := http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatal("Can't create listener on :50051/tcp", err)
	}

	grpcSrv := grpc.NewServer()

	pb.RegisterBatchIngestServer(grpcSrv, controllers.NewBatchIngestServer())

	go wrapper.Cleanup(osChan, &httpSrv, grpcSrv)

	go func() {
		log.Println("Serving gRPC on :50051/tcp")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal("Can't serve gRPC on :50051/tcp", err)
		}
	}()

	go func() {
		log.Println("Serving HTTP on :5000/tcp")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Can't start server:", err)
			osChan <- syscall.SIGTERM
		}
	}()

	wrapper.DecayWg.Wait()
}
