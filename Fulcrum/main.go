package main

import (
    "fmt"
	"net"
	"log"
    "os"
    "strconv"
    "math/rand"
    "bufio"
    "io/ioutil"
    "strings"
    "time"

	"context"

	"google.golang.org/grpc"
	pb "github.com/IIIAsmodeusIII/Tarea3/Proto"
)

const (
	port = ":50050"
    debug  = true
)

type server struct {
	pb.UnimplementedFulcrumServer
}

type PlanetaryData struct {
    name string
    version []int32
}

var server_files []PlanetaryData
var server_index int32


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



// =========================== Registro planetario ========================== //
func LineExist(file, city_name string) int{
    /*
    Function: Returns index of a line in certain file
    Input:
        file: file name, string. (e.g. example_file.txt)
        city_name: city_name to search (e.g. Mos_Pelgo)
    Output:
        if city exist, returns line_number. Else, returs -1.
    */

    // Open file
    f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    failOnError(err, "No se puede abrir el archivo: " + file)

    // Create scanner (To read line by line, not the entire file)
    scan      := bufio.NewScanner(f)

    // Iterate line by line
    i := 0
    for scan.Scan() {

        // Split data
        line   := strings.Split(scan.Text(), " ")
        city   := line[1]

        // return line_number
        if(city_name == city){
            f.Close()
            return i
        }

        i = i + 1
    }

    f.Close()
    return -1
}

func UpdateLine(file string, line int, new_data string){
    /*
    Function: Update certain line of a file with new data
    Input:
        file: file name, string. (e.g. example_file.txt)
        line: line number to update, int. (e.g. 3)
        new_data: data to replace line in file, string. (e.g. Tattoine Mos Palgo 100)
    */

    // Open file
    input, err := ioutil.ReadFile(file)
    failOnError(err, "No se pudo abrir el archivo: " + file)

    // Create new content in specific line
    lines := strings.Split(string(input), "\n")
    lines[line] = new_data

    // Write again updated content
    output := strings.Join(lines, "\n")
    err = ioutil.WriteFile(file, []byte(output), 0644)
    failOnError(err, "No se pudo escribir sobre el archivo: " + file)
}

func WriteLine(file, data string){
    /*
    Function: Append a new line in file with data
    Input:
        file: file name, string. (e.g. example_file.txt)
        data: data to append, string. (e.g. Tattoine Mos Palgo 100)
    */

    // Open "Registro planetario"
    f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    failOnError(err, "No se puede abrir el archivo: " + file)

    // Write or append new data
    _, err = f.WriteString(data)
    failOnError(err, "No se puede escribir sobre el archivo:" + file)

    f.Close()
}

func ReadLine(file string, line int) string{
    /*
    Function: Get a specific line from a file
    Input:
        file: file name, string. (e.g. example_file.txt)
        line: line number to read, int. (e.g. 3)
    */

    // Read file
    input, err := ioutil.ReadFile(file)
    failOnError(err, "No se pudo abrir el archivo: " + file)

    // Return specific line
    lines := strings.Split(string(input), "\n")
    return lines[line]
}



// ================================= Server ================================= //
/*
func (s *server) CRUD(ctx context.Context, req *pb.Command) (*pb.Data, error) {

    // Get command data
    data := strings.Split(req.Command, " ")
    comm := data[0]
    file := data[1] + ".txt"
    city := data[2]

	return &pb.Data{
		Code: int64(200),
	}, nil
}
*/
func (s *server) GetRebelds(ctx context.Context, req *pb.GetRebeldsReq) (*pb.GetRebeldsResp, error) {

    fmt.Printf("[GetRebelds Request] Planet: %v. City: %v.\n", req.Planet, req.City)

    // Get command data
    file := req.Planet + ".txt"

    // Check planet exist
    planet := FindFile(file)
    if planet == -1 {

        new_version := []int32{int32(0),int32(0),int32(0),}
        fmt.Printf("[GetRebelds Response] Ammount: %v. Version: %v.\n", 0, new_version)

        return &pb.GetRebeldsResp{
            Ammount: int32(0),
            Server: int32(-1),
            Version: new_version,
    	}, nil
    }

    // Check city exist
    line := LineExist(file, req.City)
    if line == -1 {

        new_version := []int32{int32(0),int32(0),int32(0),}
        fmt.Printf("[GetRebelds Response] Ammount: %v. Version: %v.\n", 0, new_version)

        return &pb.GetRebeldsResp{
            Ammount: int32(0),
            Server: int32(-1),
            Version: new_version,
    	}, nil
    }

    // Get data
    response, _ := strconv.Atoi(strings.Split(ReadLine(file, line), " ")[2])
    fmt.Printf("[GetRebelds Response] Ammount: %v. Version: %v.\n", response, server_files[planet].version)

	// Send response
	return &pb.GetRebeldsResp{
		Ammount: int32(response),
        Server: server_index,
        Version: server_files[planet].version,
	}, nil
}

func StartServer(){
    /*
    Function: Start a server with Fulcrum GRPC service.
    */

	// Set server
	lis, err := net.Listen("tcp", port)
    failOnError(err, "Error al crear listener.")

	s := grpc.NewServer()
	pb.RegisterFulcrumServer(s, &server{})

	log.Printf("Servidor escuchando en %v\n\n", lis.Addr())
	err = s.Serve(lis)
    failOnError(err, "Fallo al montar servidor.")
}



// ========================================================================== //
func main(){
    rand.Seed(time.Now().UnixNano())

    server_index = int32(0)
	StartServer()
}
