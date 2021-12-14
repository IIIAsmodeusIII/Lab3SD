package main

import (
    "fmt"
	"net"
	"log"
    "math/rand"
    "time"
    "flag"

	"context"

	"google.golang.org/grpc"
	pb "github.com/IIIAsmodeusIII/Tarea3/Proto"
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



// ================================= Server ================================= //
func (s *server) GetServer(ctx context.Context, req *pb.GetServerReq) (*pb.GetServerResp, error) {

    fmt.Printf("[GetServer Request] %v - %v.\n", req.Command, req.Version)
    server_index   := rand.Intn(3)

    fmt.Printf("[GetServer Response] Address: %v.\n", servers_ips[server_index])
	return &pb.GetServerResp{
		Address: servers_ips[server_index],
	}, nil
}

func (s *server) GetRebelds(ctx context.Context, req *pb.GetRebeldsReq) (*pb.GetRebeldsResp, error) {

    // Set random server
    var server_address string
    var server_index int

    var ammount int32
    var version []int32

    i := 0
    for {
        server_index   = rand.Intn(3)

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

        if r.Server == int32(-1) {
            i++
        }

        // If clock its older or server doesnt have answer, try another server
        if r.Version[server_index] < req.Version[server_index]{
            continue
        }else{
            ammount = r.Ammount
            version = r.Version
            break
        }

        if i == 6 {
            ammount = 0
            version = append(version, int32(0))
            version = append(version, int32(0))
            version = append(version, int32(0))
            break
        }
    }

    // If clock its just fine, send response
    fmt.Printf("[GetRebelds Response] Ammount: %v. Version: %v.\n", ammount, version)
    return &pb.GetRebeldsResp{
        Ammount: ammount,
        Server: int32(server_index),
        Version: version,
    }, nil
}

func StartServer(port string){
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

    // Get params
    port              := flag.String("p", ":50049", "Port to listen")
    fulcrum_1_address := flag.String("s1", "localhost:50050", "Address of first server")
    fulcrum_2_address := flag.String("s2", "localhost:50051", "Address of second server")
    fulcrum_3_address := flag.String("s3", "localhost:50052", "Address of third server")

    flag.Parse()

    // Set ip address
    servers_ips[0] = *fulcrum_1_address
    servers_ips[1] = *fulcrum_2_address
    servers_ips[2] = *fulcrum_3_address

	StartServer(*port)
}
