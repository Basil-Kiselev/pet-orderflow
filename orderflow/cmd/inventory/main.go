package main

import (
	"fmt"
	"log"
	"net"
	"orderflow/internal/gen/inventory"
	"orderflow/internal/transport/inventorygrpc"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("fail to start listener %v", err)
	}
	defer listener.Close()

	// Временный костыль, пока нет бд у сервиса
	stocks := map[string]int32{
		"item1": 28,
		"item2": 36,
	}

	stockServer := inventorygrpc.NewServer(stocks)
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	inventory.RegisterInventoryServiceServer(grpcServer, stockServer)
	fmt.Println("grpc starting on 50051 port")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("fail to start grpc server: %v", err)
	}

}
