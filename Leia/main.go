package main

import (
    "fmt"
	"log"
    "strings"
    "bufio"
    "os"

	"context"

	"google.golang.org/grpc"
	pb "github.com/IIIAsmodeusIII/Tarea3/Proto"
)

const (
	broker_address = "localhost:50049"
    debug          = true
)

type server struct {
	pb.UnimplementedBrokerServer
}

type PlanetaryData struct {
    name string
    server int32
    version []int32
}

var server_files []PlanetaryData

// ================================ Aux Func ================================ //
func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func FindFile(file string) int {

    for i, data := range server_files {
       if data.name == file {
           return i
       }
    }

    return -1
}

func DEBUG(data string){
    if debug {
        fmt.Printf("[DEBUG] %v\n", data)
    }
}

func Menu(){

    fmt.Println("[*] Type 'exit()' to close. Insert command: ")
    for true {
        fmt.Print("-> ")

        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        input   := scanner.Text()

        data    := strings.Split(input, " ")

        if input == "exit()" {
            break
        }else if(data[0] == "GetRebeldsNumber"){
            answer := AskRebelds(input)
            fmt.Printf("[->] Rebelds: %v.\n\n", answer)
        }else{
            fmt.Printf("[->] Comando '%v' desconocido.\n", input)
        }
    }
}



// =========================== Registro planetario ========================== //
func AskRebelds(command string) int32{

    // Split data
    data   := strings.Split(command, " ")
    planet := data[1]
    city   := data[2]

    // Find register
    index  := FindFile(planet)
    var register PlanetaryData

    // If not exist, append
    if index == -1 {
        new_ver     := make([]int32, 3)
        new_ver[0]   = int32(0)
        new_ver[1]   = int32(0)
        new_ver[2]   = int32(0)

        new_file    := PlanetaryData{
            name: planet,
            server: int32(-1),
            version: new_ver,
        }
        register = new_file

        server_files = append(server_files, register)
        index        = len(server_files) - 1
    }else{
        register = server_files[index]
    }

    // Connect to Broker
    conn, err := grpc.Dial(broker_address, grpc.WithInsecure(), grpc.WithBlock())
    failOnError(err, "Problema al conectar al servidor.")

    defer conn.Close()
    c := pb.NewBrokerClient(conn)

    ctx := context.Background()
    r, err := c.GetRebelds(ctx, &pb.GetRebeldsReq{
        Planet: planet,
        City: city,
        Server: register.server,
        Version: register.version,
    })
    failOnError(err, "No se pudede acceder al servicio.")

    // Update data
    server_files[index].version[0] = r.Version[0]
    server_files[index].version[1] = r.Version[1]
    server_files[index].version[2] = r.Version[2]
    server_files[index].server     = r.Server

    // Return rebelds
    return r.Ammount
}



// ========================================================================== //
func main(){
    Menu()
}