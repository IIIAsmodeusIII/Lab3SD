package main

import (
    "fmt"
	"log"
    "strings"
    "bufio"
    "os"
    "flag"

	"context"

	"google.golang.org/grpc"
	pb "github.com/IIIAsmodeusIII/Tarea3/Proto"
)

type server struct {
	pb.UnimplementedBrokerServer
}

type PlanetaryData struct {
    name string
    server int32
    ammount int32
    version []int32
}

var server_files []PlanetaryData
var broker_address string

// ================================ Aux Func ================================ //
func failOnError(err error, msg string) {
    /*
    Function: If err != nil, print error and close.
    Input:
        err: error type
        msg: msg to show if err != nil
    */
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func FindFile(file string) int {
    /*
    Function: Search if some file exist in server_files
    Input:
        file: string, file name (e.g. Tatooine.txt)
    */
    for i, data := range server_files {
       if data.name == file {
           return i
       }
    }

    return -1
}

func ShowFiles(){
    /*
	Function: Show local memory regiters
	*/
	for _, file := range server_files {
		fmt.Printf("[Data] %v.\n", file)
	}
}

func Menu(){

    fmt.Println("[*] Type GetRebeldsNumber Planet City to get rebelds in that city.")
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
        }else if(data[0] == "Files()"){
            ShowFiles()
        }else{
            fmt.Printf("[->] Comando '%v' desconocido.\n", input)
        }
    }
}



// =========================== Registro planetario ========================== //
func AskRebelds(command string) int32{
    /*
    Function: Ask to Broker about rebelds in some planet-city
    Input:
        command: string, full input to send to broker
    */
    
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
            ammount: 0,
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
    server_files[index].ammount    = r.Ammount

    // Return rebelds
    return r.Ammount
}



// ========================================================================== //
func main(){

    // Get params
	address := flag.String("ba", "localhost:50049", "Address of Broker")
	flag.Parse()

	// Set broker address
	broker_address = *address

    Menu()
}
