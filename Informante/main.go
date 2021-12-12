// Menu
// GetServerIp
//  Connect to Broker
//  Wait Broker
// Send Action
//  Connect to Fulcrum
//  Wait Fulcrum

package main

import (
	"fmt"
	"net"
	"log"
	"time"
    "strings"
	"context"

	pb "github.com/IIIAsmodeusIII/Tarea3/Proto"
    "google.golang.org/grpc"
)

const (
	broker_address = "localhost:50049"
    debug          = true
)

type server struct {
	pb.UnimplementedBrokerServer
}

type server struct {
	pb.UnimplementedFulcrumServer
}


func DEBUG(data string){
    if debug {
        fmt.Printf("[DEBUG] %v\n", data)
    }
}

func Menu(){

    // Get Command
	var choice string
	fmt.Scanln(&choice)

    // Connect to Broker
	conn, err := grpc.Dial(broker_address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Problema al conectar al servidor: %v", err)
	}
	defer conn.Close()
	c := pb.NewBrokerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Get fulcrum server address
	r, err := c.GetServer(ctx, &pb.command{
        command: choice,
        VS1: 0,
        VS2: 0,
        VS3: 0
    })

	if err != nil {
		log.Fatalf("No se pudede acceder al servicio: %v", err)
	}

    DEBUG(r.GetAddress())

}
