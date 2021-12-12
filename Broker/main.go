package main

import (
    "fmt"
	"net"
	"log"
    "math/rand"
    "time"

	"context"

	"google.golang.org/grpc"
	pb "github.com/IIIAsmodeusIII/Tarea3/Proto"
)

const (
	port   = ":50049"
    debug  = true

    fulcrum_1_address = "localhost:50050"
    fulcrum_2_address = "localhost:50051"
    fulcrum_3_address = "localhost:50052"
)

type server struct {
	pb.UnimplementedBrokerServer
}

var servers_ips [3]string

// ================================ Aux Func ================================ //
func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func DEBUG(data string){
    if debug {
        fmt.Printf("[DEBUG] %v\n", data)
    }
}



// ================================= Server ================================= //
func (s *server) GetServer(ctx context.Context, req *pb.GetServerReq) (*pb.GetServerResp, error) {

	return &pb.GetServerResp{
		Address: fulcrum_1_address,
	}, nil
}

func (s *server) GetRebelds(ctx context.Context, req *pb.GetRebeldsReq) (*pb.GetRebeldsResp, error) {

    /*
    Ask fulcrum random server
    Check Clock
    if clock fines, return
    else, check another server
    */

    // Set random server
    var server_address string
    var server_index int

    if req.Server == -1{
        server_index   = rand.Intn(3)
    }else{
        server_index   = int(req.Server)
    }

    server_address = servers_ips[server_index]
    fmt.Printf("[GetRebelds Request] Planet: %v. City: %v. Server: %v.\n", req.Planet, req.City, server_address)

    // Set conection
    conn, err := grpc.Dial(server_address, grpc.WithInsecure(), grpc.WithBlock())
    failOnError(err, "Problema al conectar al servidor.")

    defer conn.Close()
    c := pb.NewFulcrumClient(conn)

    ctx2 := context.Background()
    r, err := c.GetRebelds(ctx2, &pb.GetRebeldsReq{
        Planet: req.Planet,
        City: req.City,
        Server: req.Server,
        Version: req.Version,
    })
    failOnError(err, "No se pudede acceder al servicio")


    fmt.Printf("[GetRebelds Response] Ammount: %v. Version: %v.\n", r.Ammount, r.Version)
	return &pb.GetRebeldsResp{
		Ammount: r.Ammount,
        Server: int32(server_index),
        Version: r.Version,
	}, nil
}

func StartServer(){
    /*
    Function: Start a server with Broker GRPC service.
    */

	// Set server
	lis, err := net.Listen("tcp", port)
    failOnError(err, "Error al crear listener.")

	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, &server{})

	log.Printf("Servidor escuchando en %v\n\n", lis.Addr())
	err = s.Serve(lis)
    failOnError(err, "Fallo al montar servidor.")
}



// ========================================================================== //
func main(){

    // Set seed
    rand.Seed(time.Now().UnixNano())

    // Set ip address
    servers_ips[0] = fulcrum_1_address
    servers_ips[1] = fulcrum_2_address
    servers_ips[2] = fulcrum_3_address

	StartServer()
}
